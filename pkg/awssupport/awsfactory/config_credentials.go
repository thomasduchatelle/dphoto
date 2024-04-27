package awsfactory

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type StaticConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Token           string
	Endpoint        string
}

func (j *StaticConfig) NewConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(
		ctx,
		config.WithRegion(j.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(j.AccessKeyID, j.SecretAccessKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if j.Endpoint != "" {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           j.Endpoint,
					SigningRegion: j.Region,
				}, nil
			}

			// returning EndpointNotFoundError will allow the service to fallback to its default resolution
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})),
	)
}
