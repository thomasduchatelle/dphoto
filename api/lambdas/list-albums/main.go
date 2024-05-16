package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogaclview"
	"strings"
	"time"
)

type AlbumDTO struct {
	Name          string            `json:"name"`
	Owner         string            `json:"owner"`
	FolderName    string            `json:"folderName"`
	Start         time.Time         `json:"start"`
	End           time.Time         `json:"end"`
	TotalCount    int               `json:"totalCount"`
	SharedWith    map[string]string `json:"sharedWith,omitempty"`
	DirectlyOwned bool              `json:"directlyOwned"`
}

func Handler(request events.APIGatewayV2HTTPRequest) (common.Response, error) {
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

		levelConversion := map[aclcore.ScopeType]string{
			aclcore.AlbumVisitorScope:     "visitor",
			aclcore.AlbumContributorScope: "contributor",
		}

		restAlbums := make([]AlbumDTO, len(albums))
		for i, a := range albums {
			sharedWith := make(map[string]string)
			for email, scope := range a.SharedWith {
				if role, ok := levelConversion[scope]; ok {
					sharedWith[email.Value()] = role
				}
			}
			restAlbums[i] = AlbumDTO{
				End:           a.End,
				FolderName:    strings.TrimPrefix(a.FolderName.String(), "/"),
				Name:          a.Name,
				Owner:         a.Owner.String(),
				Start:         a.Start,
				TotalCount:    a.TotalCount,
				SharedWith:    sharedWith,
				DirectlyOwned: a.DirectlyOwned,
			}
		}
		return common.Ok(restAlbums)
	})
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
