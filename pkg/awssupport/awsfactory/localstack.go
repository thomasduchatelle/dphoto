package awsfactory

import "context"

func NewLocalstackConfig(ctx context.Context) *StaticConfig {
	return &StaticConfig{
		Region:          "us-east-1",
		AccessKeyID:     "localstack",
		SecretAccessKey: "localstack",
		Endpoint:        "http://localhost:4566",
	}
}
