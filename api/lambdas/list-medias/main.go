package main

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

type Media struct {
	Id       string    `json:"id"`       // Id is an encoded version of the business id of the media
	Type     string    `json:"type"`     // Type is PHOTO or VIDEO
	Filename string    `json:"filename"` // Filename is user-friendly and have the right extension
	Time     time.Time `json:"time"`     // Time is the datetime at which the media has been taken
	Source   string    `json:"source"`   // Source is the camera that capture the media, taken from the file metadata
}

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	ctx := context.Background()
	owner := request.PathParameters["owner"]
	folderName := request.PathParameters["folderName"]

	albumId := catalog.NewAlbumIdFromStrings(owner, folderName)

	// Extract user from authorizer context (already authenticated and authorized by Lambda Authorizer)
	_, err := common.GetCurrentUserFromContext(&request)
	if err != nil {
		return common.UnauthorizedResponse(err.Error())
	}

	// Note: IsAuthorisedToListMedias permission check is already done by the Lambda Authorizer

	log.Infof("list medias for album %s/%s", owner, folderName)

	medias, err := pkgfactory.CatalogMediaQueries(ctx).ListMedias(ctx, albumId)
	if err != nil {
		return common.InternalError(err)
	}

	resp := make([]Media, len(medias), len(medias))
	for i, media := range medias {
		resp[i] = Media{
			Id:       string(media.Id),
			Type:     string(media.Type),
			Filename: media.Filename,
			Time:     media.Details.DateTime,
			Source:   strings.Join([]string{media.Details.Make, media.Details.Model}, " "),
		}
	}

	return common.Ok(resp)
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
