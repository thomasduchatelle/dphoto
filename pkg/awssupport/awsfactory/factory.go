package awsfactory

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/singletons"
)

type ConfigFactoryFunc func(ctx context.Context) (aws.Config, error)

type ConfigFactory interface {
	NewConfig(ctx context.Context) (aws.Config, error)
}

func (f ConfigFactoryFunc) NewConfig(ctx context.Context) (aws.Config, error) {
	return f(ctx)
}

type AWSFactory struct {
	Cfg aws.Config
}

// NewAWSFactory creates a new factory for AWS clients.
func NewAWSFactory(ctx context.Context, configFactory ConfigFactory) (*AWSFactory, error) {
	cfg, err := configFactory.NewConfig(ctx)
	return &AWSFactory{
		Cfg: cfg,
	}, err
}

// GetDynamoDBClient returns a singleton instance of a DynamoDB client
func (d *AWSFactory) GetDynamoDBClient() *dynamodb.Client {
	singleton, _ := singletons.Singleton(func() (*dynamodb.Client, error) {
		return dynamodb.NewFromConfig(d.Cfg), nil
	})
	return singleton
}
