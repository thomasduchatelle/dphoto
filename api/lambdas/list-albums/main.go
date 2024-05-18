package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
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

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (common.Response, error) {
	parser := common.NewArgParser(&request)
	onlyDirectlyOwned := parser.ReadQueryParameterBool("onlyDirectlyOwned", false)

	if parser.HasViolations() {
		return parser.BadRequest()
	}

	return common.RequiresAuthenticated(&request, func(user usermodel.CurrentUser) (common.Response, error) {
		filter := catalogviews.ListAlbumsFilter{OnlyDirectlyOwned: onlyDirectlyOwned}
		albums, err := pkgfactory.AlbumView(ctx).ListAlbums(ctx, user, filter)
		if err != nil {
			return common.HandleError(err)
		}

		restAlbums := make([]AlbumDTO, len(albums))
		for i, a := range albums {
			sharedWith := make(map[string]string)
			for _, visitor := range a.Visitors {
				sharedWith[visitor.Value()] = "visitor"
			}
			restAlbums[i] = AlbumDTO{
				End:           a.End,
				FolderName:    strings.TrimPrefix(a.FolderName.String(), "/"),
				Name:          a.Name,
				Owner:         a.Owner.String(),
				Start:         a.Start,
				TotalCount:    a.MediaCount,
				SharedWith:    sharedWith,
				DirectlyOwned: a.OwnedByCurrentUser,
			}
		}
		return common.Ok(restAlbums)

	})
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
