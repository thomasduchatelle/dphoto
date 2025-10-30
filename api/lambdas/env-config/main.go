package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/meta"
)

// TODO AGENT delete this file, it's mapping i the API gatewaym and the stub data used when running locally

func Handler() (common.Response, error) {
	return common.Ok(struct {
		GoogleClientId string `json:"googleClientId"`
		Version        string `json:"version"`
	}{
		GoogleClientId: os.Getenv("GOOGLE_LOGIN_CLIENT_ID"),
		Version:        meta.Version(),
	})
}

func main() {
	lambda.Start(Handler)
}
