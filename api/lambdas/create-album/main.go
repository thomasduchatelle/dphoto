package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"time"
)

type CreateAlbumRequestDTO struct {
	Name             string    `json:"name"`
	Start            time.Time `json:"start"`
	End              time.Time `json:"end"`
	ForcedFolderName string    `json:"forcedFolderName,omitempty"`
}

type albumIdDTO struct {
	Owner      string `json:"owner"`
	FolderName string `json:"folderName"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	ctx := context.Background()

	requestDto := &CreateAlbumRequestDTO{}
	err := json.Unmarshal([]byte(request.Body), requestDto)
	if err != nil {
		return common.BadRequest(err.Error())
	}

	return common.RequiresAuthenticated(&request, func(user usermodel.CurrentUser) (common.Response, error) {
		owner, err := pkgfactory.AclCatalogAuthoriser(ctx).CanCreateAlbum(ctx, user)
		if err != nil {
			if errors.Is(err, aclcore.AccessForbiddenError) {
				return common.ForbiddenResponse(err.Error())
			}
			return common.InternalError(err)
		}

		albumId, err := pkgfactory.CreateAlbumCase(ctx).Create(ctx, catalog.CreateAlbumRequest{
			Owner:            *owner,
			Name:             requestDto.Name,
			Start:            requestDto.Start,
			End:              requestDto.End,
			ForcedFolderName: requestDto.ForcedFolderName,
		})
		if err != nil {
			switch {
			case errors.Is(err, catalog.AlbumNameMandatoryErr):
				return common.InvalidRequest("AlbumNameMandatoryErr", err.Error())
			case errors.Is(err, catalog.AlbumStartAndEndDateMandatoryErr):
				return common.InvalidRequest("AlbumStartAndEndDateMandatoryErr", err.Error())
			case errors.Is(err, catalog.AlbumEndDateMustBeAfterStartErr):
				return common.InvalidRequest("AlbumEndDateMustBeAfterStartErr", err.Error())
			case errors.Is(err, catalog.AlbumFolderNameAlreadyTakenErr):
				return common.InvalidRequest("AlbumFolderNameAlreadyTakenErr", err.Error())
			default:
				return common.InternalError(err)
			}
		}

		return common.Created(albumIdDTO{
			Owner:      albumId.Owner.Value(),
			FolderName: albumId.FolderName.String(),
		})
	})
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
