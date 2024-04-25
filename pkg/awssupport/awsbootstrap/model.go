package awsbootstrap

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
)

type ConfigFactory interface {
	NewConfig(ctx context.Context) (aws.Config, error)
}
