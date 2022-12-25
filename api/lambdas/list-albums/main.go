package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	common2 "github.com/thomasduchatelle/ephoto/api/lambdas/common"
	"github.com/thomasduchatelle/ephoto/pkg/acl/catalogaclview"
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

func Handler(request events.APIGatewayProxyRequest) (common2.Response, error) {
	parser := common2.NewArgParser(&request)
	onlyDirectlyOwned := parser.ReadQueryParameterBool("onlyDirectlyOwned", false)

	if parser.HasViolations() {
		return parser.BadRequest()
	}

	return common2.RequiresCatalogView(&request, func(catalogView *catalogaclview.View) (common2.Response, error) {
		albums, err := catalogView.ListAlbums(catalogaclview.ListAlbumsFilter{OnlyDirectlyOwned: onlyDirectlyOwned})
		if err != nil {
			return common2.Response{}, err
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
		return common2.Ok(restAlbums)
	})
}

func main() {
	common2.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
