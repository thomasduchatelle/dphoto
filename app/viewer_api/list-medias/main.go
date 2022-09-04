package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/app/viewer_api/common"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
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

	if resp, deny := common.ValidateRequest(&request, common.Or(common.CanReadAsOwner(owner), common.CanReadAlbum(owner, folderName))); deny {
		return resp, nil
	}

	log.Infof("find medias for album %s/%s", owner, folderName)
	medias, err := catalog.ListMedias(owner, fmt.Sprintf("/%s", folderName), catalog.PageRequest{})
	if err != nil {
		return common.InternalError(err)
	}

	resp := make([]Media, len(medias.Content), len(medias.Content))
	for i, media := range medias.Content {
		resp[i] = Media{
			Id:       media.Id,
			Type:     string(media.Type),
			Filename: media.Filename,
			Time:     media.Details.DateTime,
			Source:   strings.Trim(fmt.Sprintf("%s %s", media.Details.Make, media.Details.Model), " "),
		}
	}

	return common.Ok(resp)
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
