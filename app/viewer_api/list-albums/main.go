package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thomasduchatelle/dphoto/app/viewer_api/common"
	"github.com/thomasduchatelle/dphoto/domain/acl/catalogaclview"
	"strings"
	"time"
)

type AlbumDTO struct {
	Name       string    `json:"name"`
	Owner      string    `json:"owner"`
	FolderName string    `json:"folderName"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	TotalCount int       `json:"totalCount"`
	SharedTo   []string  `json:"sharedTo,omitempty"`
}

func Handler(request events.APIGatewayProxyRequest) (common.Response, error) {
	parser := common.NewArgParser(&request)
	onlyDirectlyOwned := parser.ReadQueryParameterBool("onlyDirectlyOwned", false)

	if parser.HasViolations() {
		return parser.BadRequest()
	}

	return common.RequiresCatalogView(&request, func(catalogView *catalogaclview.View) (common.Response, error) {
		albums, err := catalogView.ListAlbums(catalogaclview.ListAlbumsFilter{OnlyDirectlyOwned: onlyDirectlyOwned})
		if err != nil {
			return common.Response{}, err
		}

		restAlbums := make([]AlbumDTO, len(albums))
		for i, a := range albums {
			restAlbums[i] = AlbumDTO{
				End:        a.End,
				FolderName: strings.TrimPrefix(a.FolderName, "/"),
				Name:       a.Name,
				Owner:      a.Owner,
				Start:      a.Start,
				TotalCount: a.TotalCount,
				SharedTo:   a.SharedTo,
			}
		}
		return common.Ok(restAlbums)
	})
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
