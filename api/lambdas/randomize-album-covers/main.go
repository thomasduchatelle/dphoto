package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

type CoverDTO struct {
	MediaId  string `json:"mediaId"`
	Filename string `json:"filename"`
	Origin   string `json:"origin"`
}

type RandomizeCoversResponseDTO struct {
	Covers []CoverDTO `json:"covers"`
}

type coverRandomizer interface {
	Randomize(ctx context.Context, albumId catalog.AlbumId) ([]catalog.Cover, error)
}

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	ctx := context.Background()
	return handle(ctx, request, pkgfactory.RandomizeAlbumCoversCase(ctx))
}

func handle(ctx context.Context, request events.APIGatewayV2HTTPRequest, randomizer coverRandomizer) (common.Response, error) {
	argParser := common.NewArgParser(&request)
	owner := ownermodel.Owner(argParser.ReadPathParameterString("owner"))
	folderName := catalog.NewFolderName(argParser.ReadPathParameterString("folderName"))

	if argParser.HasViolations() {
		return argParser.BadRequest()
	}

	if _, err := common.GetCurrentUserFromContext(&request); err != nil {
		return common.UnauthorizedResponse(err.Error())
	}

	albumId := catalog.AlbumId{Owner: owner, FolderName: folderName}

	covers, err := randomizer.Randomize(ctx, albumId)
	if err != nil {
		switch {
		case errors.Is(err, catalog.AlbumNotFoundErr):
			return common.NotFound(map[string]string{"error": err.Error()})
		case errors.Is(err, aclcore.AccessForbiddenError):
			return common.ForbiddenResponse(err.Error())
		default:
			return common.InternalError(err)
		}
	}

	return common.Ok(RandomizeCoversResponseDTO{Covers: convertCoversForREST(covers)})
}

func convertCoversForREST(covers []catalog.Cover) []CoverDTO {
	dto := make([]CoverDTO, len(covers))
	for i, c := range covers {
		dto[i] = CoverDTO{
			MediaId:  c.MediaId.Value(),
			Filename: c.Filename,
			Origin:   string(c.Origin),
		}
	}
	return dto
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
