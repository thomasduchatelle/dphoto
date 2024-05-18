package catalogviewstoacl

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type ScopeReadRepositoryPort interface {
	ListScopesByOwner(ctx context.Context, owner ownermodel.Owner, types ...aclcore.ScopeType) ([]*aclcore.Scope, error)
	ListScopesByUser(ctx context.Context, id usermodel.UserId, types ...aclcore.ScopeType) ([]*aclcore.Scope, error)
}

type FindAlbumSharingToAdapter struct {
	ScopeRepository ScopeReadRepositoryPort
}

func (f *FindAlbumSharingToAdapter) GetAlbumSharingGrid(ctx context.Context, owner ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error) {
	scopes, err := f.ScopeRepository.ListScopesByOwner(ctx, owner, aclcore.AlbumVisitorScope, aclcore.AlbumContributorScope)

	if err != nil || len(scopes) == 0 {
		return nil, err
	}

	grid := make(map[catalog.AlbumId][]usermodel.UserId)
	for _, scope := range scopes {
		albumId := catalog.AlbumId{Owner: owner, FolderName: catalog.NewFolderName(scope.ResourceId)}
		if list, ok := grid[albumId]; ok {
			grid[albumId] = append(list, scope.GrantedTo)
		} else {
			grid[albumId] = []usermodel.UserId{scope.GrantedTo}
		}
	}

	return grid, nil
}

func (f *FindAlbumSharingToAdapter) ListAlbumIdsSharedWithUser(ctx context.Context, userId usermodel.UserId) ([]catalog.AlbumId, error) {
	shared, err := f.ScopeRepository.ListScopesByUser(ctx, userId, aclcore.AlbumVisitorScope)

	var albums []catalog.AlbumId
	for _, share := range shared {
		albums = append(albums, catalog.AlbumId{Owner: share.ResourceOwner, FolderName: catalog.NewFolderName(share.ResourceId)})
	}

	return albums, err
}
