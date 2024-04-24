package migrator

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TransformationRun struct {
	DynamoDBClient    *dynamodb.Client
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
	GeneratePatches(run *TransformationRun, item map[string]types.AttributeValue) ([]types.WriteRequest, error)
}
