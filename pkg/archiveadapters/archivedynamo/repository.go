// Package archivedynamo extends catalogdynamo to add media locations to the main table
package archivedynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	dynamoutils2 "github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
)

func New(sess *session.Session, tableName string) (archive.ARepositoryAdapter, error) {
	return &repository{
		db:    dynamodb.New(sess),
		table: tableName,
	}, nil
}

func Must(a archive.ARepositoryAdapter, err error) archive.ARepositoryAdapter {
	if err != nil {
		panic(err)
	}

	return a
}

type repository struct {
	db    *dynamodb.DynamoDB
	table string
}

func (r *repository) FindById(owner, id string) (string, error) {
	item, err := r.db.GetItem(&dynamodb.GetItemInput{
		Key:       marshalMediaLocationPK(owner, id),
		TableName: &r.table,
	})
	if err != nil {
		return "", errors.Wrapf(err, "FindById %s, %s failed", owner, id)
	}

	if len(item.Item) == 0 {
		return "", archive.NotFoundError
	}

	_, location, err := unmarshalMediaLocation(item.Item)
	return location, err
}

func (r *repository) AddLocation(owner, id, key string) error {
	location, err := marshalMediaLocation(owner, id, key)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal location")
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		Item:      location,
		TableName: &r.table,
	})
	return errors.Wrapf(err, "failed to upsert media location %s - %s - %s", owner, id, key)
}

func (r *repository) FindByIds(owner string, ids []string) (map[string]string, error) {
	keys := make([]map[string]*dynamodb.AttributeValue, len(ids), len(ids))
	for i, id := range ids {
		keys[i] = marshalMediaLocationPK(owner, id)
	}

	locations := make(map[string]string)

	stream := dynamoutils2.NewGetStream(dynamoutils2.NewGetBatchItem(r.db, r.table, ""), keys, dynamoutils2.DynamoReadBatchSize)
	for stream.HasNext() {
		id, key, err := unmarshalMediaLocation(stream.Next())
		if err != nil {
			return nil, err
		}

		locations[id] = key
	}

	return locations, stream.Error()
}

func (r *repository) UpdateLocations(owner string, locations map[string]string) error {
	requests := make([]*dynamodb.WriteRequest, len(locations), len(locations))
	i := 0
	for id, key := range locations {
		location, err := marshalMediaLocation(owner, id, key)
		requests[i] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: location,
			},
		}
		if err != nil {
			return err
		}

		i++
	}

	return dynamoutils2.BufferedWriteItems(r.db, requests, r.table, dynamoutils2.DynamoWriteBatchSize)
}

func (r *repository) FindIdsFromKeyPrefix(keyPrefix string) (map[string]string, error) {
	pairs := make(map[string]string)

	err := r.db.QueryPages(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":prefix": {S: &keyPrefix},
		},
		KeyConditionExpression: aws.String("LocationKeyPrefix = :prefix"),
		IndexName:              aws.String("ReverseLocationIndex"),
		TableName:              &r.table,
	}, func(output *dynamodb.QueryOutput, last bool) bool {
		for _, item := range output.Items {
			mediaId, storeKey, err := unmarshalMediaLocation(item)
			if err != nil {
				log.WithError(err).Errorf("failed unmarshaling item %+v for prefix %s - skipping", item, keyPrefix)
			} else {
				pairs[mediaId] = storeKey
			}
		}

		return true
	})

	if len(pairs) == 0 {
		pairs = nil
	}
	return pairs, err
}
