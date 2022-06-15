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
			result, err := db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]*dynamodb.WriteRequest{
					table: buffer,
				},
			})

			if err != nil {
				return err
			}

			buffer = buffer[:0]
			if unprocessed, ok := result.UnprocessedItems[table]; ok && len(unprocessed) > 0 {
				buffer = append(buffer, unprocessed...)
			}

			return nil

		}, retry, func(err error, duration time.Duration) {
			log.WithFields(log.Fields{
				"Table":    table,
				"Duration": duration,
			}).WithError(err).Warnf("Retrying inserting media (buffer len %d)", len(buffer))
		})
		if err != nil {
			return errors.Wrapf(err, "failed to insert batch %+v", buffer)
		}
	}

	return nil
}
