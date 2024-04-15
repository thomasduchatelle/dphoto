// Package archivedynamo extends catalogdynamo to add media locations to the main table
package archivedynamo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	dynamoutils "github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutilsv2"
)

func New(cfg aws.Config, tableName string) (archive.ARepositoryAdapter, error) {
	return &repository{
		db:    dynamodb.NewFromConfig(cfg),
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
	db    *dynamodb.Client
	table string
}

func (r *repository) FindById(owner, id string) (string, error) {
	item, err := r.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
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

	_, err = r.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      location,
		TableName: &r.table,
	})
	return errors.Wrapf(err, "failed to upsert media location %s - %s - %s", owner, id, key)
}

func (r *repository) FindByIds(owner string, ids []string) (map[string]string, error) {
	keys := make([]map[string]types.AttributeValue, len(ids), len(ids))
	for i, id := range ids {
		keys[i] = marshalMediaLocationPK(owner, id)
	}

	locations := make(map[string]string)

	stream := dynamoutils.NewGetStream(context.TODO(), dynamoutils.NewGetBatchItem(r.db, r.table, ""), keys, dynamoutils.DynamoReadBatchSize)
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
	requests := make([]types.WriteRequest, len(locations), len(locations))
	i := 0
	for id, key := range locations {
		location, err := marshalMediaLocation(owner, id, key)
		requests[i] = types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: location,
			},
		}
		if err != nil {
			return err
		}

		i++
	}

	return dynamoutils.BufferedWriteItems(context.TODO(), r.db, requests, r.table, dynamoutils.DynamoWriteBatchSize)
}

func (r *repository) FindIdsFromKeyPrefix(keyPrefix string) (map[string]string, error) {
	pairs := make(map[string]string)

	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.Key("LocationKeyPrefix").Equal(expression.Value(keyPrefix))).
		Build()
	if err != nil {
		return nil, err
	}

	paginator := dynamodb.NewQueryPaginator(r.db, &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		IndexName:                 aws.String("ReverseLocationIndex"),
		TableName:                 &r.table,
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		for _, item := range page.Items {
			mediaId, storeKey, err := unmarshalMediaLocation(item)
			if err != nil {
				log.WithError(err).Errorf("failed unmarshaling item %+v for prefix %s - skipping", item, keyPrefix)
			} else {
				pairs[mediaId] = storeKey
			}
		}
	}

	if len(pairs) == 0 {
		pairs = nil
	}
	return pairs, nil
}
