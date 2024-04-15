package dynamoutilsv2

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Stream is inspired from Java streams to chain transformations in a functional programming style
type Stream interface {
	HasNext() bool                         // HasNext returns true if it has another element and no error has been raised
	Next() map[string]types.AttributeValue // Next return current element and move forward the cursor
	Error() error                          // Error returns the error that interrupted the Stream
	Count() int64                          // Count return the number of element found so far
}

// GetStreamAdapter is building the batch query before passing it to dynamodb.
type GetStreamAdapter interface {
	BatchGet(ctx context.Context, key []map[string]types.AttributeValue) (*dynamodb.BatchGetItemOutput, error)
}

// DynamoBatchGetItem is implemented by dynamodb.New()
type DynamoBatchGetItem interface {
	BatchGetItem(ctx context.Context, params *dynamodb.BatchGetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchGetItemOutput, error)
}

// DynamoQuery is implemented by dynamodb.New()
type DynamoQuery interface {
	Query(ctx context.Context, d *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

// DynamoBatchWriteItem is implemented by dynamodb.New()
type DynamoBatchWriteItem interface {
	BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)
}
