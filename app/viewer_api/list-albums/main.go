package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/app/viewer_api/common"
	"github.com/thomasduchatelle/dphoto/domain/accesscontrol"
	"github.com/thomasduchatelle/dphoto/domain/catalogacl"
	"time"
)

type AlbumDTO struct {
	Name       string    `json:"name"`
	Owner      string    `json:"owner"`
	FolderName string    `json:"folderName"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	TotalCount int       `json:"totalCount"`
}

func Handler(request events.APIGatewayProxyRequest) (common.Response, error) {
	parser := common.NewArgParser(&request)
	onlyDirectlyOwned := parser.ReadQueryParameterBool("onlyDirectlyOwned", false)

	if parser.HasViolations() {
		return parser.BadRequest()
	}

	return common.RequiresAuthorisation(&request, func(catalogView catalogacl.View) (common.Response, error) {
		albums, err := catalogView.ListAlbums(catalogacl.ListAlbumsFilter{OnlyDirectlyOwned: onlyDirectlyOwned})
		if err != nil {
			if errors.Is(err, accesscontrol.AccessForbiddenError) {
				return common.Response{}, err
			}
			return common.InternalError(errors.Wrapf(err, "failed to fetch albums"))
		}

		restAlbums := make([]AlbumDTO, len(albums))
		for i, a := range albums {
			restAlbums[i] = AlbumDTO{
				End:        a.End,
				FolderName: a.FolderName,
				Name:       a.Name,
				Owner:      a.Owner,
				Start:      a.Start,
				TotalCount: a.TotalCount,
			}
		}
		return common.Ok(restAlbums)
	})
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
