package dynamotestutils

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"sort"
	"testing"
)

func SetContent(t *testing.T, db *dynamodb.DynamoDB, table string, entries []map[string]*dynamodb.AttributeValue) {
	err := clearContent(db, table)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
	err = insertAll(db, table, entries)
	if err != nil {
		assert.FailNow(t, err.Error())
	}
}

func clearContent(db *dynamodb.DynamoDB, table string) error {
	var deletionRequests []*dynamodb.WriteRequest

	stream := dynamoutils.NewScanStream(db, table)
	for stream.HasNext() {
		entry := stream.Next()
		key := make(map[string]*dynamodb.AttributeValue)
		key["PK"], _ = entry["PK"]
		key["SK"], _ = entry["SK"]
		deletionRequests = append(deletionRequests, &dynamodb.WriteRequest{DeleteRequest: &dynamodb.DeleteRequest{Key: key}})
	}

	return dynamoutils.BufferedWriteItems(db, deletionRequests, table, dynamoutils.DynamoWriteBatchSize)
}

func insertAll(db *dynamodb.DynamoDB, table string, entries []map[string]*dynamodb.AttributeValue) error {
	if len(entries) == 0 {
		return nil
	}

	var requests []*dynamodb.WriteRequest
	for _, entry := range entries {
		requests = append(requests, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{Item: entry},
		})
	}

	_, err := db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			table: requests,
		}})
	return err
}

func AssertAfter(t *testing.T, db *dynamodb.DynamoDB, table string, expected []map[string]*dynamodb.AttributeValue) bool {
	content, err := dynamoutils.AsSlice(dynamoutils.NewScanStream(db, table))
	if !assert.NoError(t, err) {
		return false
	}

	sort.Slice(content, func(i, j int) bool {
		if *content[i]["PK"].S == *content[j]["PK"].S {
			return *content[i]["SK"].S < *content[j]["SK"].S
		}
		return *content[i]["PK"].S < *content[j]["PK"].S
	})

	return assert.Equal(t, expected, content)
}
