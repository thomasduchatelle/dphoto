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
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type RenameAlbumRequestDTO struct {
	Name       string `json:"name"`
	FolderName string `json:"folderName,omitempty"`
}

type albumIdDTO struct {
	Owner      string `json:"owner"`
	FolderName string `json:"folderName"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	ctx := context.Background()

	argParser := common.NewArgParser(&request)
	owner := ownermodel.Owner(argParser.ReadPathParameterString("owner"))
	folderName := catalog.NewFolderName("/" + argParser.ReadPathParameterString("folderName"))

	if argParser.HasViolations() {
		return argParser.BadRequest()
	}

	albumId := catalog.AlbumId{Owner: owner, FolderName: folderName}

	requestDto := &RenameAlbumRequestDTO{}
	err := json.Unmarshal([]byte(request.Body), requestDto)
	if err != nil {
		return common.BadRequest(err.Error())
	}

	return common.RequiresAuthenticated(&request, func(user usermodel.CurrentUser) (common.Response, error) {
		err := pkgfactory.AclCatalogAuthoriser(ctx).CanRenameAlbum(ctx, user, albumId)
		if err != nil {
			if errors.Is(err, aclcore.AccessForbiddenError) {
				return common.ForbiddenResponse(err.Error())
			}
			return common.InternalError(err)
		}

		newAlbumId, err := common.Factory.RenameAlbumCase(ctx).RenameAlbum(ctx, albumId, requestDto.Name, requestDto.FolderName)
		if err != nil {
			switch {
			case errors.Is(err, catalog.AlbumNameMandatoryErr):
				return common.UnprocessableEntityResponse("AlbumNameMandatoryErr", err.Error())
			case errors.Is(err, catalog.AlbumFolderNameAlreadyTakenErr):
				return common.UnprocessableEntityResponse("AlbumFolderNameAlreadyTakenErr", err.Error())
			case errors.Is(err, catalog.AlbumNotFoundErr):
				return common.NotFound(map[string]string{
					"error": err.Error(),
				})
			default:
				return common.InternalError(err)
			}
		}

		return common.Ok(albumIdDTO{
			Owner:      newAlbumId.Owner.Value(),
			FolderName: common.ConvertFolderNameForREST(newAlbumId.FolderName),
		})
	})
}

func main() {
	common.BootstrapCatalogAndArchiveDomains()

	lambda.Start(Handler)
}
