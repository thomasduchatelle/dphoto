package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	version "github.com/thomasduchatelle/dphoto/pkg/meta"
)

func Handler() (common.Response, error) {
	return common.Ok(map[string]string{
		"version": version.Version(),
	})
}

func main() {
	lambda.Start(Handler)
}
