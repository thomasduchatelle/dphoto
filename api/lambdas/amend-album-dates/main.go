package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
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

	// Extract user from authorizer context (already authenticated and authorized by Lambda Authorizer)
	_, err = common.GetCurrentUserFromContext(&request)
	if err != nil {
		return common.UnauthorizedResponse(err.Error())
	}

	// Note: CanAmendAlbumDates permission check is already done by the Lambda Authorizer

	albumId := catalog.AlbumId{Owner: ownermodel.Owner(owner), FolderName: catalog.NewFolderName(folderName)}

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
}

func main() {
	common.BootstrapCatalogAndArchiveDomains()

	lambda.Start(Handler)
}
