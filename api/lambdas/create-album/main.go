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

	// Extract user from authorizer context (already authenticated and authorized by Lambda Authorizer)
	user, err := common.GetCurrentUserFromContext(&request)
	if err != nil {
		return common.UnauthorizedResponse(err.Error())
	}

	// Note: CanCreateAlbum permission check is already done by the Lambda Authorizer
	// We just need to get the owner from the authorizer result
	if user.Owner == nil {
		return common.ForbiddenResponse("user does not have an owner")
	}

	albumId, err := common.Factory.CreateAlbumCase(ctx).Create(ctx, catalog.CreateAlbumRequest{
		Owner:            *user.Owner,
		Name:             requestDto.Name,
		Start:            requestDto.Start,
		End:              requestDto.End,
		ForcedFolderName: requestDto.ForcedFolderName,
	})
	if err != nil {
		switch {
		case errors.Is(err, catalog.AlbumNameMandatoryErr):
			return common.UnprocessableEntityResponse("AlbumNameMandatoryErr", err.Error())
		case errors.Is(err, catalog.AlbumStartAndEndDateMandatoryErr):
			return common.UnprocessableEntityResponse("AlbumStartAndEndDateMandatoryErr", err.Error())
		case errors.Is(err, catalog.AlbumEndDateMustBeAfterStartErr):
			return common.UnprocessableEntityResponse("AlbumEndDateMustBeAfterStartErr", err.Error())
		case errors.Is(err, catalog.AlbumFolderNameAlreadyTakenErr):
			return common.UnprocessableEntityResponse("AlbumFolderNameAlreadyTakenErr", err.Error())
		default:
			return common.InternalError(err)
		}
	}

	return common.Created(albumIdDTO{
		Owner:      albumId.Owner.Value(),
		FolderName: common.ConvertFolderNameForREST(albumId.FolderName),
	})
}

func main() {
	common.BootstrapCatalogAndArchiveDomains()

	lambda.Start(Handler)
}
