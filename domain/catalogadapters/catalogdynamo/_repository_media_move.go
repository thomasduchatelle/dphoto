package catalogdynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/dynamoutils"
	"time"
)

func (r *rep) TransferMedias(filter *catalog.UpdateMediaFilter, newFolderName string) ([]string, error) {
	queries, err := r.findMediasQueries(filter, aws.String("PK, SK"))
	if err != nil {
		return nil, err
	}

	it := dynamoutils.NewQueryStream(r.db, queries)
	if !it.HasNext() {
		return nil, nil
	}

	var updates []*dynamodb.UpdateItemInput

	for it.HasNext() {
		mediaKey := it.Next()

		newAlbumKey := AlbumIndexedKey(filter.Owner, newFolderName)
		updateValues, err := dynamodbattribute.MarshalMap(map[string]string{
			":albumPK": newAlbumKey.AlbumIndexPK,
		})
		if err != nil {
			return nil, err
		}

		updates = append(updates, &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: updateValues,
			Key:                       mediaKey,
			TableName:                 &r.table,
			UpdateExpression:          aws.String("SET AlbumIndexPK=:albumPK"),
		})
	}
	if err = it.Error(); err != nil {
		return "", 0, err
	}

	err = dynamoutils.BufferedWriteItems(r.db, moveOrders, r.table, dynamoWriteBatchSize)
	if err != nil {
		return "", 0, err
	}

	err = r.updateItems(updates)
	if err != nil {
		return "", 0, err
	}

	err = r.markMoveTransactionReady(transactionId)

	log.WithFields(log.Fields{
		"Filter":      filter,
		"Destination": newFolderName,
	}).Infoln(it.Count(), "media virtually moved to new album")

	return transactionId, int(it.Count()), err
}

func (r *rep) startMoveTransaction(owner string) (string, map[string]*dynamodb.AttributeValue, error) {
	transactionId := fmt.Sprintf("MOVE_ORDER#%s#%s#%s", owner, time.Now().Format(time.RFC3339), uuid.New())

	transaction, err := marshalMoveTransaction(transactionId, transactionPreparing)
	if err != nil {
		return transactionId, nil, err
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		Item:      transaction,
		TableName: &r.table,
	})
	return transactionId, transaction, err
}

func (r *rep) markMoveTransactionReady(transactionId string) error {
	transaction, err := marshalMoveTransaction(transactionId, transactionReady)
	if err != nil {
		return err
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		Item:      transaction,
		TableName: &r.table,
	})
	return err
}

func (r *rep) FindReadyMoveTransactions(owner string) ([]*catalog.MoveTransaction, error) {
	ownedTransactionsPrefix := fmt.Sprintf("MOVE_ORDER#%s", owner) // see startMoveTransaction
	crawler := dynamoutils.NewQueryStream(r.db, []*dynamodb.QueryInput{
		{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":status":     mustAttribute(transactionReady),
				":ownedStart": mustAttribute(ownedTransactionsPrefix + "#"),
				":ownedEnd":   mustAttribute(ownedTransactionsPrefix + "$"),
			},
			IndexName:              aws.String("MoveTransaction"),
			KeyConditionExpression: aws.String("MoveTransactionStatus = :status AND PK BETWEEN :ownedStart AND :ownedEnd"),
			TableName:              &r.table,
		},
	})

	var transactionPartitionKeys []*catalog.MoveTransaction
	for crawler.HasNext() {
		var transaction MediaMoveTransactionData
		err := dynamodbattribute.UnmarshalMap(crawler.Next(), &transaction)
		if err != nil {
			return nil, err
		}

		count, err := r.countNumberOfMediaToBeMoved(transaction.PK)
		if err != nil {
			return nil, err
		}
		transactionPartitionKeys = append(transactionPartitionKeys, &catalog.MoveTransaction{
			TransactionId: transaction.PK,
			Count:         count,
		})
	}

	return transactionPartitionKeys, crawler.Error()
}

// FindFilesToMove returns a page of media to move (25 like a write batch of dynamodb) and the next page token
func (r *rep) FindFilesToMove(transactionId, pageToken string) ([]*catalog.MovedMedia, string, error) {
	startKey, err := r.unmarshalPageToken(pageToken)
	if err != nil {
		return nil, "", err
	}

	exprValues, err := dynamodbattribute.MarshalMap(map[string]string{
		":transactionId": transactionId,
	})
	if err != nil {
		return nil, "", err
	}

	results, err := r.db.Query(&dynamodb.QueryInput{
		ExclusiveStartKey:         startKey,
		ExpressionAttributeValues: exprValues,
		IndexName:                 aws.String("MoveOrder"),
		KeyConditionExpression:    aws.String("MoveTransaction = :transactionId"),
		Limit:                     aws.Int64(r.findMovedMediaBatchSize),
		TableName:                 &r.table,
	})
	if err != nil {
		return nil, "", err
	}

	nextPageToken, err := r.marshalPageToken(results.LastEvaluatedKey)
	if err != nil {
		return nil, "", err
	}

	destinationFolders := make(map[string]string)
	currentLocationKeys := make([]map[string]*dynamodb.AttributeValue, len(results.Items))

	for index, attributes := range results.Items {
		order, err := unmarshalMoveOrder(attributes)
		if err != nil {
			return nil, "", err
		}
		locationKey, err := mediaLocationKeyFromMediaKey(attributes)
		if err != nil {
			return nil, "", err
		}

		destinationFolders[order.PK] = order.DestinationFolder
		currentLocationKeys[index] = locationKey
	}

	stream := dynamoutils.NewGetStream(dynamoutils.NewGetBatchItem(r.db, r.table, *nil), currentLocationKeys, dynamoReadBatchSize)

	movedMedias := make([]*catalog.MovedMedia, 0, len(results.Items))
	for stream.HasNext() {
		attributes := stream.Next()

		var location MediaLocationData
		err := dynamodbattribute.UnmarshalMap(attributes, &location)
		if err != nil {
			return nil, "", err
		}

		movedMedias = append(movedMedias, &catalog.MovedMedia{
			Signature: catalog.MediaSignature{
				SignatureSha256: location.SignatureHash,
				SignatureSize:   location.SignatureSize,
			},
			SourceFolderName: location.FolderName,
			SourceFilename:   location.Filename,
			TargetFolderName: destinationFolders[location.PK],
			TargetFilename:   location.Filename,
		})
	}

	return movedMedias, nextPageToken, stream.Error()
}

func (r *rep) UpdateMediasLocation(owner string, transactionId string, moves []*catalog.MovedMedia) error {
	locations := make([]*dynamodb.WriteRequest, len(moves)*2)
	for i, move := range moves {
		locationItem, err := marshalMediaLocationFromMoveOrder(owner, move)
		if err != nil {
			return err
		}

		moveKey, err := dynamodbattribute.MarshalMap(mediaMoveOrderPrimaryKey(owner, &move.Signature, transactionId))
		if err != nil {
			return err
		}

		locations[i*2] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{Item: locationItem},
		}

		locations[i*2+1] = &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{Key: moveKey},
		}
	}

	return dynamoutils.BufferedWriteItems(r.db, locations, r.table, dynamoWriteBatchSize)
}

func (r *rep) DeleteEmptyMoveTransaction(transactionId string) error {
	count, err := r.countNumberOfMediaToBeMoved(transactionId)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.Errorf("Move transaction must be empty to be deleted. %s contains %d media to move", transactionId, count)
	}

	key, err := dynamodbattribute.MarshalMap(moveTransactionPrimaryKey(transactionId))
	if err != nil {
		return err
	}

	_, err = r.db.DeleteItem(&dynamodb.DeleteItemInput{
		Key:       key,
		TableName: &r.table,
	})
	return err
}

func (r *rep) countNumberOfMediaToBeMoved(transactionId string) (int, error) {
	result, err := r.db.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {S: &transactionId},
		},
		IndexName:              aws.String("MoveOrder"),
		KeyConditionExpression: aws.String("MoveTransaction = :pk"),
		Select:                 aws.String(dynamodb.SelectCount),
		TableName:              &r.table,
	})
	if err != nil {
		return 0, err
	}
	return int(*result.Count), nil
}
