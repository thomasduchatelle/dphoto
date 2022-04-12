package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thomasduchatelle/dphoto/app/viewer_api/common"
	version "github.com/thomasduchatelle/dphoto/domain/meta"
)

func Handler() (common.Response, error) {
	return common.Ok(map[string]string{
		"version": version.Version(),
	})
}

func main() {
	lambda.Start(Handler)
}
