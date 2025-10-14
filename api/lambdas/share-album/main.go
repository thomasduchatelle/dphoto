package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	ctx := context.Background()

	owner := request.PathParameters["owner"]
	folderName := request.PathParameters["folderName"]
	email := request.PathParameters["email"]

	// Extract user from authorizer context (already authenticated and authorized by Lambda Authorizer)
	_, err := common.GetCurrentUserFromContext(&request)
	if err != nil {
		return common.UnauthorizedResponse(err.Error())
	}

	// Note: CanShareAlbum permission check is already done by the Lambda Authorizer

	albumId := catalog.NewAlbumIdFromStrings(owner, folderName)
	userId := usermodel.NewUserId(email)

	method := request.RequestContext.HTTP.Method
	switch method {
	case "PUT":
		err := pkgfactory.AclCatalogShare(ctx).ShareAlbumWith(ctx, albumId, userId)
		if errors.Is(err, catalog.AlbumNotFoundErr) {
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
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
