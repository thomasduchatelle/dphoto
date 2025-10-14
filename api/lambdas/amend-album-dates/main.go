package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type AmendAlbumDatesRequestDTO struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	ctx := context.Background()

	owner := request.PathParameters["owner"]
	folderName := request.PathParameters["folderName"]

	requestDto := &AmendAlbumDatesRequestDTO{}
	err := json.Unmarshal([]byte(request.Body), requestDto)
	if err != nil {
		return common.BadRequest(err.Error())
	}

	return common.RequiresAuthenticated(&request, func(user usermodel.CurrentUser) (common.Response, error) {
		albumId := catalog.AlbumId{Owner: ownermodel.Owner(owner), FolderName: catalog.NewFolderName(folderName)}

		if err := pkgfactory.AclCatalogAuthoriser(ctx).CanAmendAlbumDates(ctx, user, albumId); err != nil {
			if errors.Is(err, aclcore.AccessForbiddenError) {
				return common.ForbiddenResponse(err.Error())
			}
			return common.InternalError(err)
		}

		err = common.Factory.AmendAlbumDatesCase(ctx).AmendAlbumDates(ctx, albumId, requestDto.Start, requestDto.End)
		if err != nil {
			switch {
			case errors.Is(err, catalog.AlbumNotFoundErr):
				return common.NotFound(err.Error())
			case errors.Is(err, catalog.AlbumStartAndEndDateMandatoryErr):
				return common.UnprocessableEntityResponse("AlbumStartAndEndDateMandatoryErr", err.Error())
			case errors.Is(err, catalog.AlbumEndDateMustBeAfterStartErr):
				return common.UnprocessableEntityResponse("AlbumEndDateMustBeAfterStartErr", err.Error())
			case errors.Is(err, catalog.OrphanedMediasErr):
				return common.UnprocessableEntityResponse("OrphanedMediasErr", err.Error())
			default:
				return common.InternalError(err)
			}
		}

		return common.NoContent()
	})
}

func main() {
	common.BootstrapCatalogAndArchiveDomains()

	lambda.Start(Handler)
}
