package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thomasduchatelle/dphoto/app/viewer_api/common"
	"os"
)

func Handler() (common.Response, error) {
	return common.Ok(struct {
		GoogleClientId string `json:"googleClientId"`
	}{
		GoogleClientId: os.Getenv("GOOGLE_LOGIN_CLIENT_ID"),
	})
}

func main() {
	lambda.Start(Handler)
}
