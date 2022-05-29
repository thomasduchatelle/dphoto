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
	Id       string    `json:"id"`       // Id is an encoded version of the business id of the media
	Type     string    `json:"type"`     // Type is PHOTO or VIDEO
	Filename string    `json:"filename"` // Filename is user-friendly and have the right extension
	Time     time.Time `json:"time"`     // Time is the datetime at which the media has been taken
	Source   string    `json:"source"`   // Source is the camera that capture the media, taken from the file metadata
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
		id, err := common.EncodeMediaId(media.Signature)
		if err != nil {
			return common.InternalError(err)
		}

		resp[i] = Media{
			Id:       id,
			Type:     string(media.Type),
			Filename: media.Filename,
			Time:     media.Details.DateTime,
			Source:   strings.Trim(fmt.Sprintf("%s %s", media.Details.Make, media.Details.Model), " "),
		}
	}

	return common.Ok(resp)
}

func main() {
	common.Bootstrap()

	lambda.Start(Handler)
}
