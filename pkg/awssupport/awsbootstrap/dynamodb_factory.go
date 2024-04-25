package awsbootstrap

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBFactory struct {
	ConfigFactory
}

func (d *DynamoDBFactory) GetDynamoDBClient(ctx context.Context) (*dynamodb.Client, error) {
	cfg, err := d.NewConfig(ctx)
	if err != nil {
		return nil, err
	}

	return dynamodb.NewFromConfig(cfg), nil
}
