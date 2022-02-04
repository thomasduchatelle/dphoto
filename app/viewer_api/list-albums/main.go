package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/app/viewer_api/common"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/oauth"
)

func Handler(request events.APIGatewayProxyRequest) (common.Response, error) {
	owner, _ := request.PathParameters["owner"]
	if resp, deny := common.ValidateRequest(request, oauth.NewAuthoriseQuery("owner").WithOwner(owner, "READ")); deny {
		return resp, nil
	}

	common.BootstrapCatalogDomain(owner)

	albums, err := catalog.FindAllAlbumsWithStats()
	if err != nil {
		return common.InternalError(errors.Wrapf(err, "failed to fetch albums"))
	}

	return common.Ok(albums)
}

func main() {
	common.Bootstrap()

	lambda.Start(Handler)
}
