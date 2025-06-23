package main

import (
	"context"
	"fmt"
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

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	ctx := context.Background()

	owner := request.PathParameters["owner"]
	folderName := request.PathParameters["folderName"]

	albumId := catalog.AlbumId{Owner: ownermodel.Owner(owner), FolderName: catalog.NewFolderName(folderName)}

	if owner == "" || folderName == "" {
		return common.BadRequest("Missing required path parameters: owner or folderName")
	}

	return common.RequiresAuthenticated(&request, func(user usermodel.CurrentUser) (common.Response, error) {

		err := pkgfactory.AclCatalogAuthoriser(ctx).CanDeleteAlbum(ctx, user, albumId)
		if err != nil {
			if errors.Is(err, aclcore.AccessForbiddenError) {
				return common.ForbiddenResponse(err.Error())
			}
			return common.InternalError(err)
		}

		err = common.Factory.CreateAlbumDeleteCase(ctx).DeleteAlbum(ctx, albumId)
		if err != nil {
			switch {
			case errors.Is(err, catalog.OrphanedMediasErr):
				return common.UnprocessableEntityResponse("OrphanedMedias", err.Error())
			case errors.Is(err, aclcore.AccessForbiddenError):
				return common.UnauthorizedResponse(fmt.Sprintf("You are not authorized to delete the album %s", albumId))
			default:
				return common.UnprocessableEntityResponse("InternalError", err.Error())
			}
		}

		return common.NoContent()
	})
}

func main() {
	common.BootstrapCatalogAndArchiveDomains()

	// This lambda can handle both POST (create) and DELETE (delete) album requests.
	lambda.Start(Handler)
}
