package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogaclview"
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

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	owner := request.PathParameters["owner"]
	folderName := request.PathParameters["folderName"]

	return common.RequiresCatalogView(&request, func(catalogView *catalogaclview.View) (common.Response, error) {
		log.Infof("list medias for album %s/%s", owner, folderName)
		medias, err := catalogView.ListMediasFromAlbum(owner, fmt.Sprintf("/%s", folderName))
		if err != nil {
			return common.Response{}, err
		}

		resp := make([]Media, len(medias.Content), len(medias.Content))
		for i, media := range medias.Content {
			resp[i] = Media{
				Id:       media.Id,
				Type:     string(media.Type),
				Filename: media.Filename,
				Time:     media.Details.DateTime,
				Source:   strings.Join([]string{media.Details.Make, media.Details.Model}, " "),
			}
		}

		return common.Ok(resp)
	})
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
