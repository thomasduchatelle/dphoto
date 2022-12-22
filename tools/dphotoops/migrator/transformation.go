package migrator

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type TransformationRun struct {
	Session           *session.Session
	DynamoDBClient    *dynamodb.DynamoDB
	DynamoDBTableName string
	TopicARN          string
	Counter           Counter
}

type Counter map[string]int

func (c Counter) Inc(key string, delta int) {
	c[key] += delta
}

type PreScanTransformation interface {
	PreScan(run *TransformationRun) error
}

type ScanTransformation interface {
	GeneratePatches(run *TransformationRun, item map[string]*dynamodb.AttributeValue) ([]*dynamodb.WriteRequest, error)
}
