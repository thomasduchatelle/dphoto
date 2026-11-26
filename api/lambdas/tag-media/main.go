package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	ctx := context.Background()

	// Extract user from authorizer context (already authenticated and authorized by Lambda Authorizer)
	_, err := common.GetCurrentUserFromContext(&request)
	if err != nil {
		return common.UnauthorizedResponse(err.Error())
	}

	// Note: owner-only permission check is already done by the Lambda Authorizer

	argParser := common.NewArgParser(&request)
	owner := ownermodel.Owner(argParser.ReadPathParameterString("owner"))
	mediaId := catalog.MediaId(argParser.ReadPathParameterString("mediaId"))
	rawTag := argParser.ReadPathParameterString("tag")
	if argParser.HasViolations() {
		return argParser.BadRequest()
	}

	tag, err := catalog.NewTag(rawTag)
	if err != nil {
		return common.BadRequest(map[string]string{"error": err.Error()})
	}

	method := request.RequestContext.HTTP.Method
	switch method {
	case "PUT":
		err = pkgfactory.TagServiceCase(ctx).AddTag(ctx, owner, mediaId, tag)
		if errors.Is(err, catalog.AlbumNotFoundErr) || errors.Is(err, catalog.MediaNotFoundError) {
			return common.NotFound(fmt.Sprintf("media %s/%s hasn't been found", owner, mediaId))
		} else if err != nil {
			return common.InternalError(err)
		}

	case "DELETE":
		err = pkgfactory.TagServiceCase(ctx).RemoveTag(ctx, owner, mediaId, tag)
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
