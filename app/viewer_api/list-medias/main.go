package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/app/viewer_api/common"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogmodel"
	"github.com/thomasduchatelle/dphoto/domain/oauth"
	"strings"
	"time"
)

type Media struct {
	Id     string    `json:"id"`
	Type   string    `json:"type"`
	Path   string    `json:"path"`
	Time   time.Time `json:"time"`   // Time is the datetime at which the media has been taken
	Source string    `json:"source"` // Source is the camera that capture the media, taken from the file metadata
}

func Handler(request events.APIGatewayProxyRequest) (common.Response, error) {
	owner, _ := request.PathParameters["owner"]
	folderName, _ := request.PathParameters["folderName"]

	if resp, deny := common.ValidateRequest(&request, oauth.NewAuthoriseQuery("owner").WithOwner(owner, "READ")); deny {
		return resp, nil
	}

	log.Infof("find medias for album %s/%s", owner, folderName)
	medias, err := catalog.ListMedias(owner, fmt.Sprintf("/%s", folderName), catalogmodel.PageRequest{})
	if err != nil {
		return common.InternalError(err)
	}

	resp := make([]Media, len(medias.Content), len(medias.Content))
	for i, media := range medias.Content {
		restKey := fmt.Sprintf("/api/v1/owners/%s/medias/%s/%d", owner, media.Signature.SignatureSha256, media.Signature.SignatureSize)

		resp[i] = Media{
			Id:     restKey,
			Type:   string(media.Type),
			Path:   fmt.Sprintf("%s/%s", restKey, media.Filename),
			Time:   media.Details.DateTime,
			Source: strings.Trim(fmt.Sprintf("%s %s", media.Details.Make, media.Details.Model), " "),
		}
	}

	return common.Ok(resp)
}

func main() {
	common.Bootstrap()

	lambda.Start(Handler)
}
