package dynamoutils

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

func (s *simpleBatchGetExecutor) BatchGet(keys []map[string]*dynamodb.AttributeValue) (*dynamodb.BatchGetItemOutput, error) {
	var expression *string
	if s.projectionExpression != "" {
		expression = &s.projectionExpression
	}

	return s.delegate.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			s.tableName: {
				Keys:                 keys,
				ProjectionExpression: expression,
			},
		},
	})
}
