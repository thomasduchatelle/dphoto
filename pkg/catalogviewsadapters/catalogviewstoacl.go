package catalogviewsadapters

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type ListScopesByOwnerPort interface {
	ListScopesByOwner(owner ownermodel.Owner, types ...aclcore.ScopeType) ([]*aclcore.Scope, error)
}

type FindAlbumSharingToAdapter struct {
	ListScopesByOwnerPort ListScopesByOwnerPort
}

func (f *FindAlbumSharingToAdapter) GetAlbumSharingGrid(ctx context.Context, owner ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error) {
	scopes, err := f.ListScopesByOwnerPort.ListScopesByOwner(owner, aclcore.AlbumVisitorScope, aclcore.AlbumContributorScope)

	if err != nil || len(scopes) == 0 {
		return nil, err
	}

	grid := make(map[catalog.AlbumId][]usermodel.UserId)
	for _, scope := range scopes {
		albumId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName(scope.ResourceId)}
		if list, ok := grid[albumId]; ok {
			grid[albumId] = append(list, usermodel.UserId(scope.GrantedTo))
		} else {
			grid[albumId] = []usermodel.UserId{usermodel.UserId(scope.GrantedTo)}
		}
	}

	return grid, nil
}
