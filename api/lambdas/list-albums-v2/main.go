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

// TODO NOTES - can the details be IN the view to not "reload" ?
// TODO NOTES - can the endpoint be 'visibleAlbums' or 'catalog' or '/v2/albums' ; is there is anything else to add to the BFF endpoint ?

type OwnerDetailsDTO struct {
	Name  string           `json:"name"`
	Users []UserDetailsDTO `json:"users"`
}

type SharedWithDTO struct {
	Role string         `json:"role"`
	User UserDetailsDTO `json:"user"`
}

type UserDetailsDTO struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

type AlbumDTO struct {
	Name          string          `json:"name"`
	Owner         OwnerDetailsDTO `json:"owner"`
	FolderName    string          `json:"folderName"` // FolderName is stripped from the leading '/'
	Start         time.Time       `json:"start"`
	End           time.Time       `json:"end"`
	TotalCount    int             `json:"totalCount"`
	SharedWith    []SharedWithDTO `json:"sharedWith,omitempty"`
	DirectlyOwned bool            `json:"directlyOwned"` // DirectlyOwned is TRUE when album is owned by the user requesting the API
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
		for i, album := range albums {
			sharedWith := make(map[string]string)
			for _, visitor := range album.Visitors {
				sharedWith[visitor.Value()] = "visitor"
			}
			restAlbums[i] = AlbumDTO{
				End:           album.End,
				FolderName:    strings.TrimPrefix(album.FolderName.String(), "/"),
				Name:          album.Name,
				Owner:         album.Owner.String(),
				Start:         album.Start,
				TotalCount:    album.MediaCount,
				SharedWith:    sharedWith,
				DirectlyOwned: album.OwnedByCurrentUser,
			}
		}
		return common.Ok(restAlbums)

	})
}

func main() {
	common.BootstrapCatalogDomain()

	lambda.Start(Handler)
}
