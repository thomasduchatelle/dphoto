package dynamoutils

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Stream is inspired from Java streams to chain transformations in a functional programming style
type Stream interface {
	HasNext() bool                             // HasNext returns true if it has another element and no error has been raised
	Next() map[string]*dynamodb.AttributeValue // Next return current element and move forward the cursor
	Error() error                              // Error returns the error that interrupted the Stream
	Count() int64                              // Count return the number of element found so far
}

// GetStreamAdapter is building the batch query before passing it to dynamodb.
type GetStreamAdapter interface {
	BatchGet([]map[string]*dynamodb.AttributeValue) (*dynamodb.BatchGetItemOutput, error)
}

// DynamoBatchGetItem is implemented by dynamodb.New()
type DynamoBatchGetItem interface {
	BatchGetItem(*dynamodb.BatchGetItemInput) (*dynamodb.BatchGetItemOutput, error)
}

// DynamoQuery is implemented by dynamodb.New()
type DynamoQuery interface {
	Query(d *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
}

// DynamoBatchWriteItem is implemented by dynamodb.New()
type DynamoBatchWriteItem interface {
	BatchWriteItem(*dynamodb.BatchWriteItemInput) (*dynamodb.BatchWriteItemOutput, error)
}
