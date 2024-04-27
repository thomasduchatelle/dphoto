package awsfactory

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var DefaultConfigFactory = ConfigFactoryFunc(func(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx)
})

func NewContextualConfigFactory() ConfigFactory {
	return DefaultConfigFactory
}
