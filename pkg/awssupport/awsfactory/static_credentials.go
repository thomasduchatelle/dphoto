package awsfactory

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type StaticCredentials struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

func (j *StaticCredentials) awsConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(
		ctx,
		config.WithRegion(j.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(j.AccessKeyID, j.SecretAccessKey, "")),
	)
}
