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
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	ctx := context.Background()

	owner := request.PathParameters["owner"]
	folderName := request.PathParameters["folderName"]
	email := request.PathParameters["email"]

	return common.RequiresAuthenticated(&request, func(user usermodel.CurrentUser) (common.Response, error) {
		albumId := catalog.NewAlbumIdFromStrings(owner, folderName)
		userId := usermodel.NewUserId(email)

		if err := pkgfactory.AclCatalogAuthoriser(ctx).CanShareAlbum(ctx, user, albumId); err != nil {
			if errors.Is(err, aclcore.AccessForbiddenError) {
				return common.ForbiddenResponse(err.Error())
			}
			return common.InternalError(err)
		}

		method := request.RequestContext.HTTP.Method
		switch method {
		case "PUT":
			err := pkgfactory.AclCatalogShare(ctx).ShareAlbumWith(ctx, albumId, userId)
			if errors.Is(err, catalog.AlbumNotFoundError) {
				return common.NotFound(fmt.Sprintf("%s/%s hasn't been found", owner, folderName))
			} else if err != nil {
				return common.InternalError(err)
			}

		case "DELETE":
			err := pkgfactory.AclCatalogUnShare(ctx).StopSharingAlbum(albumId, userId)
			if err != nil {
				return common.InternalError(err)
			}

		default:
			return common.BadRequest(fmt.Sprintf("%s method is not supported", method))
		}

		return common.NoContent()
	})
}

func translateScope(role string) (aclcore.ScopeType, error) {
	switch role {
	case "visitor":
		return aclcore.AlbumVisitorScope, nil
	case "contributor":
		return aclcore.AlbumContributorScope, nil
	default:
		return "", errors.Errorf("'%s' role is not supported. Expected 'visitor' or 'contributor'", role)
	}
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
