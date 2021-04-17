package dynamo

import (
	"duchatelle.io/dphoto/dphoto/catalog"
	"duchatelle.io/dphoto/dphoto/config"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	versionTagName       = "Version"
	tableVersion         = "1.1" // to be bumped manually when schema is updated
	dynamoWriteBatchSize = 25
	dynamoReadBatchSize  = 100
	defaultPage          = 50
)

func init() {
	catalog.Repository = &rep{
		db:                      dynamodb.New(config.AwsConfig.GetSession()),
		table:                   "AlbumMedia",
		RootOwner:               "#ROOT",
		findMovedMediaBatchSize: int64(dynamoWriteBatchSize),
	}
}

type rep struct {
	db                      *dynamodb.DynamoDB
	table                   string
	RootOwner               string // PK ID, '#ROOT' if single tenant.
	localDynamodb           bool   // some feature like tagging are not available
	findMovedMediaBatchSize int64
}

func (r *rep) InsertMedias(medias []catalog.CreateMediaRequest) error {
	requests := make([]*dynamodb.WriteRequest, len(medias)*2)
	for index, media := range medias {
		mediaEntry, locationEntry, err := r.marshalMedia(&media)
		if err != nil {
			return errors.Wrapf(err, "Failed mapping media %s", fmt.Sprint(media))
		}

		requests[index*2] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: locationEntry,
			},
		}
		requests[index*2+1] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: mediaEntry,
			},
		}
	}

	return r.bufferedWriteItems(requests)
}

func (r *rep) FindMedias(folderName string, request catalog.PageRequest) (*catalog.MediaPage, error) {
	albumKey := r.albumIndexedKey(r.RootOwner, folderName)
	albumIndexPK, err := dynamodbattribute.Marshal(albumKey.AlbumIndexPK)
	if err != nil {
		return nil, err
	}

	limit := aws.Int64(defaultPage)
	if request.Size > 0 {
		limit = aws.Int64(request.Size)
	}

	startKey, err := r.unmarshalPageToken(request.NextPage)
	if err != nil {
		return nil, err
	}

	page, err := r.db.Query(&dynamodb.QueryInput{
		ExclusiveStartKey: startKey,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":albumKey":    albumIndexPK,
			":excludeMeta": r.mustAttribute("$"),
		},
		IndexName:              aws.String("AlbumIndex"),
		KeyConditionExpression: aws.String("AlbumIndexPK = :albumKey AND AlbumIndexSK >= :excludeMeta"),
		Limit:                  limit,
		TableName:              &r.table,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load medias for album %s with page request %+v", folderName, request)
	}

	medias := make([]*catalog.MediaMeta, 0, len(page.Items))
	for _, attributes := range page.Items {
		media, err := r.unmarshalMediaMetaData(attributes)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal media %v", attributes)
		}
		medias = append(medias, media)
	}

	nextPage, err := r.marshalPageToken(page.LastEvaluatedKey)
	return &catalog.MediaPage{
		NextPage: nextPage,
		Content:  medias,
	}, err
}

func (r *rep) FindExistingSignatures(signatures []*catalog.MediaSignature) ([]*catalog.MediaSignature, error) {
	keys := make([]map[string]*dynamodb.AttributeValue, len(signatures))
	for index, signature := range signatures {
		key, err := dynamodbattribute.MarshalMap(r.mediaPrimaryKey(signature))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal media keys from signature %+v", signature)
		}

		keys[index] = key
	}

	stream := r.bufferedBatchGetItem(keys, aws.String("SignatureSize, SignatureHash"))

	found := make([]*catalog.MediaSignature, 0, len(signatures))
	for stream.HasNext() {
		attributes := stream.Next()

		var signature MediaData
		err := dynamodbattribute.UnmarshalMap(attributes, &signature)
		if err != nil {
			return nil, err
		}

		found = append(found, &catalog.MediaSignature{
			SignatureSha256: signature.SignatureHash,
			SignatureSize:   signature.SignatureSize,
		})
	}

	return found, stream.Error()
}

func (r *rep) FindMediaLocationsSignatures(signatures []*catalog.MediaSignature) ([]*catalog.MediaSignatureAndLocation, error) {
	keys := make([]map[string]*dynamodb.AttributeValue, len(signatures))
	for index, signature := range signatures {
		key, err := dynamodbattribute.MarshalMap(r.mediaLocationPrimaryKey(signature))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal media keys from signature %+v", signature)
		}

		keys[index] = key
	}

	stream := r.bufferedBatchGetItem(keys, aws.String("FolderName, Filename, SignatureHash, SignatureSize"))

	locations := make([]*catalog.MediaSignatureAndLocation, 0, len(signatures))
	for stream.HasNext() {
		attributes := stream.Next()

		var location MediaLocationData
		err := dynamodbattribute.UnmarshalMap(attributes, &location)
		if err != nil {
			return nil, err
		}

		locations = append(locations, &catalog.MediaSignatureAndLocation{
			Location: catalog.MediaLocation{
				FolderName: location.FolderName,
				Filename:   location.Filename,
			},
			Signature: catalog.MediaSignature{
				SignatureSha256: location.SignatureHash,
				SignatureSize:   location.SignatureSize,
			},
		})
	}

	return locations, stream.Error()
}

func (r *rep) bufferedBatchGetItem(keys []map[string]*dynamodb.AttributeValue, projectionExpression *string) Stream {
	return NewGetStream(r, keys, projectionExpression, dynamoReadBatchSize)
}

func (r *rep) bufferedQueriesCrawler(queries []*dynamodb.QueryInput) Stream {
	return NewQueryStream(r, queries)
}

func (r *rep) bufferedWriteItems(requests []*dynamodb.WriteRequest) error {
	retry := backoff.NewExponentialBackOff()
	buffer := make([]*dynamodb.WriteRequest, 0, dynamoWriteBatchSize)

	for len(buffer) > 0 || len(requests) > 0 {
		end := cap(buffer) - len(buffer)
		if end > len(requests) {
			end = len(requests)
		}

		if end > 0 {
			buffer = append(buffer, requests[:end]...)
			requests = requests[end:]
		}

		err := backoff.RetryNotify(func() error {
			result, err := r.db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]*dynamodb.WriteRequest{
					r.table: buffer,
				},
			})

			if err != nil {
				return err
			}

			buffer = buffer[:0]
			if unprocessed, ok := result.UnprocessedItems[r.table]; ok && len(unprocessed) > 0 {
				buffer = append(buffer, unprocessed...)
			}

			return nil

		}, retry, func(err error, duration time.Duration) {
			log.WithFields(log.Fields{
				"Duration": duration,
				"Buffer":   buffer,
			}).WithError(err).Warn("Retrying inserting media")
		})
		if err != nil {
			return errors.Wrapf(err, "failed to insert batch %+v", buffer)
		}
	}

	return nil
}

func (r *rep) bufferedUpdateItems(updates []*dynamodb.UpdateItemInput) error {
	for _, update := range updates {
		_, err := r.db.UpdateItem(update)
		if err != nil {
			return err
		}
	}

	return nil
}

// unmarshalPageToken take a base64 page id and decode it into a dynamodb key. Return nil, nil if nextPageToken is blank
func (r *rep) unmarshalPageToken(nextPageToken string) (map[string]*dynamodb.AttributeValue, error) {
	var startKey map[string]*dynamodb.AttributeValue
	if !isBlank(nextPageToken) {
		startKeyJson, err := base64.StdEncoding.DecodeString(nextPageToken)
		if err != nil {
			return startKey, errors.Wrapf(err, "Invalid nextPageToken %s", nextPageToken)
		}

		err = json.Unmarshal(startKeyJson, &startKey)
		if err != nil {
			return startKey, errors.Wrapf(err, "Invalid nextPageToken %s", startKeyJson)
		}
	}

	return startKey, nil
}

// marshalPageToken take a dynamodb key and encode it (JSON+BASE64) to be used by service. Return empty string when key is empty
func (r *rep) marshalPageToken(key map[string]*dynamodb.AttributeValue) (string, error) {
	if len(key) == 0 {
		return "", nil
	}

	keyJson, err := json.Marshal(key)
	return base64.StdEncoding.EncodeToString(keyJson), err
}

// panic if value can't be converted into an attribute
func (r *rep) mustAttribute(value interface{}) *dynamodb.AttributeValue {
	attribute, err := dynamodbattribute.Marshal(value)
	if err != nil {
		panic(err)
	}
	return attribute
}

func (r *rep) FindMediaLocations(signature catalog.MediaSignature) ([]*catalog.MediaLocation, error) {
	mediaPK := r.mediaPrimaryKey(&signature)
	queryValues, err := dynamodbattribute.MarshalMap(map[string]string{
		":mediaPK":    mediaPK.PK,
		":noMetadata": "$",
	})
	if err != nil {
		return nil, err
	}

	result, err := r.db.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: queryValues,
		KeyConditionExpression:    aws.String("PK = :mediaPK AND SK > :noMetadata"),
		TableName:                 &r.table,
	})
	if err != nil {
		return nil, err
	}

	location, orders, err := r.unmarshalMediaItems(result.Items)
	if location == nil {
		return nil, errors.Errorf("Location not found for media %+v", signature)
	}

	locations := make([]*catalog.MediaLocation, len(orders)+1)
	locations[0] = &catalog.MediaLocation{
		FolderName: location.FolderName,
		Filename:   location.Filename,
	}
	for index, order := range orders {
		locations[index+1] = &catalog.MediaLocation{
			FolderName: order.DestinationFolder,
			Filename:   location.Filename,
		}
	}

	return locations, nil
}
