package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
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
	ctx := context.Background()
	owner := request.PathParameters["owner"]
	folderName := request.PathParameters["folderName"]

	albumId := catalog.NewAlbumIdFromStrings(owner, folderName)

	return common.RequiresAuthenticated(&request, func(user usermodel.CurrentUser) (common.Response, error) {
		log.Infof("list medias for album %s/%s", owner, folderName)
		err := pkgfactory.AclCatalogAuthoriser(ctx).IsAuthorisedToListMedias(ctx, user, albumId)
		if errors.Is(err, catalogacl.ErrAccessDenied) {
			return common.ForbiddenResponse(err.Error())
		}
		if err != nil {
			return common.InternalError(err)
		}

		medias, err := pkgfactory.CatalogMediaQueries(ctx).ListMedias(ctx, albumId)

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
	})
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
