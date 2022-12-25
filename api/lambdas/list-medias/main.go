package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	common2 "github.com/thomasduchatelle/ephoto/api/lambdas/common"
	"github.com/thomasduchatelle/ephoto/pkg/acl/catalogaclview"
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

func Handler(request events.APIGatewayProxyRequest) (common2.Response, error) {
	owner := request.PathParameters["owner"]
	folderName := request.PathParameters["folderName"]

	return common2.RequiresCatalogView(&request, func(catalogView *catalogaclview.View) (common2.Response, error) {
		log.Infof("list medias for album %s/%s", owner, folderName)
		medias, err := catalogView.ListMediasFromAlbum(owner, fmt.Sprintf("/%s", folderName))
		if err != nil {
			return common2.Response{}, err
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

		return common2.Ok(resp)
	})
}

func main() {
	common2.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
