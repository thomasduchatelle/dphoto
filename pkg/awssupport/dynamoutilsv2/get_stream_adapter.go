package dynamoutilsv2

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type simpleBatchGetExecutor struct {
	delegate             DynamoBatchGetItem
	tableName            string
	projectionExpression string
}

func NewGetBatchItem(delegate DynamoBatchGetItem, tableName string, projectionExpression string) GetStreamAdapter {
	return &simpleBatchGetExecutor{
		delegate:             delegate,
		tableName:            tableName,
		projectionExpression: projectionExpression,
	}
}

func (s *simpleBatchGetExecutor) BatchGet(ctx context.Context, keys []map[string]types.AttributeValue) (*dynamodb.BatchGetItemOutput, error) {
	var expression *string
	if s.projectionExpression != "" {
		expression = &s.projectionExpression
	}

	return s.delegate.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			s.tableName: {
				Keys:                 keys,
				ProjectionExpression: expression,
			},
		},
	})
}
