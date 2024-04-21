package localstack

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	Region   = "us-east-1"
	Endpoint = "http://localhost:4566"
)

var (
	WithUsePathPrefix = func(options *s3.Options) {
		options.UsePathStyle = true // required for localstack testing on UNIX
	}
)

// Config creates a static configuration using localstack running on localhost.
func Config(ctx context.Context) aws.Config {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("localstack", "localstack", "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           Endpoint,
				PartitionID:   "aws",
				SigningRegion: region,
			}, nil
		})),
	)
	if err != nil {
		panic(err)
	}

	return cfg
}

func S3(cfg aws.Config) *s3.Client {
	return s3.NewFromConfig(cfg, WithUsePathPrefix)
}
