package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/app/viewer_api/common"
	"github.com/thomasduchatelle/dphoto/domain/oauth"
	"time"
)

type Media struct {
	Id     string    `json:"id"`
	Type   string    `json:"type"`
	Path   string    `json:"path"`
	Time   time.Time `json:"time"`
	Source string    `json:"source"`
}

func Handler(request events.APIGatewayProxyRequest) (common.Response, error) {
	owner, _ := request.PathParameters["owner"]
	folderName, _ := request.PathParameters["folderName"]

	if resp, deny := common.ValidateRequest(request, oauth.NewAuthoriseQuery("owner").WithOwner(owner, "READ")); deny {
		return resp, nil
	}

	log.Infof("find medias for album %s/%s", owner, folderName)
	return common.Ok([]Media{})
}

func main() {
	common.Bootstrap()

	lambda.Start(Handler)
}
