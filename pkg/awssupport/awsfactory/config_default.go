package awsfactory

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func NewContextualConfigFactory(ctx context.Context) ConfigFactory {
	return ConfigFactoryFunc(func(ctx context.Context) (aws.Config, error) {
		return config.LoadDefaultConfig(ctx)
	})
}
