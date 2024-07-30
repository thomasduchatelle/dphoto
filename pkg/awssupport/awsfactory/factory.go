package awsfactory

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/thomasduchatelle/dphoto/pkg/singletons"
)

const (
	LocalstackEndpoint = "http://localhost:4566" // LocalstackEndpoint is the default location for localstack
)

type AWSFactory interface {
	GetCfg() aws.Config
	GetDynamoDBClient() *dynamodb.Client
	GetS3Client() *s3.Client
	GetSNSClient() *sns.Client
	GetSQSClient() *sqs.Client
}

func StaticAWSFactory(ctx context.Context, config StaticCredentials) (AWSFactory, error) {
	cfg, err := config.awsConfig(ctx)
	return &ClientsFactory{Cfg: cfg}, err
}

func ContextualAWSFactory(ctx context.Context) (AWSFactory, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	return &ClientsFactory{Cfg: cfg}, err
}

func LocalstackAWSFactory(ctx context.Context, endpoint string) (AWSFactory, error) {
	const region = "us-east-1"
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("localstack", "localstack", "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		})),
	)
	return &LocalstackClientsFactory{ClientsFactory{Cfg: cfg}}, err
}

type ClientsFactory struct {
	Cfg aws.Config
}

func (d *ClientsFactory) GetCfg() aws.Config {
	return d.Cfg
}

// GetDynamoDBClient returns a singleton instance of a DynamoDB client
func (d *ClientsFactory) GetDynamoDBClient() *dynamodb.Client {
	singleton, _ := singletons.Singleton(func() (*dynamodb.Client, error) {
		return dynamodb.NewFromConfig(d.Cfg), nil
	})
	return singleton
}

func (d *ClientsFactory) GetS3Client() *s3.Client {
	return s3.NewFromConfig(d.Cfg)
}

func (d *ClientsFactory) GetSNSClient() *sns.Client {
	return sns.NewFromConfig(d.Cfg)
}

func (d *ClientsFactory) GetSQSClient() *sqs.Client {
	return sqs.NewFromConfig(d.Cfg)
}

type LocalstackClientsFactory struct {
	ClientsFactory
}

func (d *LocalstackClientsFactory) GetS3Client() *s3.Client {
	WithUsePathPrefix := func(options *s3.Options) {
		options.UsePathStyle = true // required for localstack testing on UNIX
	}

	return s3.NewFromConfig(d.Cfg, WithUsePathPrefix)
}
