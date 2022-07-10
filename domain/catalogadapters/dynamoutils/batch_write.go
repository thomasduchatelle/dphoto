package dynamoutils

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

// BufferedWriteItems writes items in batch with backoff and replay unprocessed items
func BufferedWriteItems(db DynamoBatchWriteItem, requests []*dynamodb.WriteRequest, table string, dynamoWriteBatchSize int) error {
	batchSize := dynamoWriteBatchSize
	if batchSize > DynamoWriteBatchSize {
		batchSize = DynamoWriteBatchSize
	}

	buffer := make([]*dynamodb.WriteRequest, 0, batchSize)

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

		result, err := db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
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
