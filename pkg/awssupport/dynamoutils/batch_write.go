package dynamoutils

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

func BufferedPutItems(ctx context.Context, db DynamoBatchWriteItem, items []interface{}, table string, dynamoWriteBatchSize int) error {
	values, err := attributevalue.MarshalList(items)
	if err != nil {
		return err
	}

	requests := make([]types.WriteRequest, len(items), len(items))
	for i, value := range values {
		if ss, ok := value.(*types.AttributeValueMemberM); ok {
			requests[i] = types.WriteRequest{
				PutRequest: &types.PutRequest{Item: ss.Value},
			}
		} else {
			return errors.Errorf("items must be structures or maps with a valid PK ; got %+v", items[i])
		}
	}

	return BufferedWriteItems(ctx, db, requests, table, dynamoWriteBatchSize)
}

// BufferedWriteItems writes items in batch with backoff and replay unprocessed items
func BufferedWriteItems(ctx context.Context, db DynamoBatchWriteItem, requests []types.WriteRequest, table string, dynamoWriteBatchSize int) error {
	batchSize := dynamoWriteBatchSize
	if batchSize > DynamoWriteBatchSize {
		batchSize = DynamoWriteBatchSize
	}

	buffer := make([]types.WriteRequest, 0, batchSize)

	retry := backoff.NewExponentialBackOff()

	for len(buffer) > 0 || len(requests) > 0 {
		end := cap(buffer) - len(buffer)
		if end > len(requests) {
			end = len(requests)
		}

		if end > 0 {
			buffer = append(buffer, requests[:end]...)
			requests = requests[end:]
		}

		result, err := db.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				table: buffer,
			},
		})
		if err != nil {
			return errors.Wrapf(err, "failed to insert batch %+v", buffer)
		}

		buffer = buffer[:0]
		if unprocessed, ok := result.UnprocessedItems[table]; ok && len(unprocessed) > 0 {
			buffer = append(buffer, unprocessed...)
			log.WithField("Table", table).Warnf("%d unprocessed item(s) while inserting in dynamodb table '%s'. Will retry with exponential backoff.", len(unprocessed), table)
			time.Sleep(retry.NextBackOff())
		}
	}

	return nil
}
