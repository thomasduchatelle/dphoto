package catalogviewstoacl

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type ResourceIds map[ownermodel.Owner][]string

func (i ResourceIds) Append(owner ownermodel.Owner, resourceId string) {
	if _, ok := i[owner]; !ok {
		i[owner] = make([]string, 0)
	}
	i[owner] = append(i[owner], resourceId)
}

func NewResourceIds() ResourceIds {
	return make(ResourceIds)
}

type ScopeReadRepositoryPort interface {
	ListScopesByOwner(ctx context.Context, owner ownermodel.Owner, types ...aclcore.ScopeType) ([]*aclcore.Scope, error)
	ListScopesByUser(ctx context.Context, id usermodel.UserId, types ...aclcore.ScopeType) ([]*aclcore.Scope, error)
	ListScopesByResource(ctx context.Context, resourceIds ResourceIds, types ...aclcore.ScopeType) ([]*aclcore.Scope, error)
}

type CatalogToACLAdapter struct {
	ScopeRepository ScopeReadRepositoryPort
}

func (f *CatalogToACLAdapter) GetAlbumSharingGrid(ctx context.Context, owner ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error) {
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

func (f *CatalogToACLAdapter) ListAlbumIdsSharedWithUser(ctx context.Context, userId usermodel.UserId) ([]catalog.AlbumId, error) {
	shared, err := f.ScopeRepository.ListScopesByUser(ctx, userId, aclcore.AlbumVisitorScope)

	var albums []catalog.AlbumId
	for _, share := range shared {
		albums = append(albums, catalog.AlbumId{Owner: share.ResourceOwner, FolderName: catalog.NewFolderName(share.ResourceId)})
	}

	return albums, err
}

func (f *CatalogToACLAdapter) ListUsersWhoCanAccessAlbum(ctx context.Context, albumIds ...catalog.AlbumId) (map[catalog.AlbumId][]catalogviews.Availability, error) {
	grid := make(map[catalog.AlbumId][]catalogviews.Availability)

	ids := NewResourceIds()
	for _, id := range albumIds {
		ids.Append(id.Owner, id.FolderName.String())
	}

	visitorScopes, err := f.ScopeRepository.ListScopesByResource(ctx, ids, aclcore.AlbumVisitorScope, aclcore.AlbumContributorScope)
	if err != nil {
		return nil, err
	}

	for _, visitorScope := range visitorScopes {
		albumId := catalog.AlbumId{Owner: visitorScope.ResourceOwner, FolderName: catalog.NewFolderName(visitorScope.ResourceId)}
		availability := catalogviews.VisitorAvailability(visitorScope.GrantedTo)

		if list, ok := grid[albumId]; ok {
			grid[albumId] = append(list, availability)
		} else {
			grid[albumId] = []catalogviews.Availability{availability}
		}

	}

	for owner := range ids {
		ownerScopes, err := f.ScopeRepository.ListScopesByOwner(ctx, owner, aclcore.MainOwnerScope)
		if err != nil {
			return nil, err
		}

		for _, ownerScope := range ownerScopes {
			for _, albumId := range albumIds {
				if albumId.Owner == ownerScope.ResourceOwner {
					availability := catalogviews.OwnerAvailability(ownerScope.GrantedTo)

					if list, ok := grid[albumId]; ok {
						grid[albumId] = append(list, availability)
					} else {
						grid[albumId] = []catalogviews.Availability{availability}
					}

					break
				}
			}
		}

	}

	return grid, nil
}
