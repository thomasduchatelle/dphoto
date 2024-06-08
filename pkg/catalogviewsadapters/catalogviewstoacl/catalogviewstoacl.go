package catalogviewstoacl

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"slices"
)

// TODO is catalogviewstoacl the right package to have these translations catalog -> ACL ? catalogacl is doing the same ...

type ScopeReadRepositoryPort interface {
	ListScopesByOwner(ctx context.Context, owner ownermodel.Owner, types ...aclcore.ScopeType) ([]*aclcore.Scope, error)
	ListScopesByUser(ctx context.Context, id usermodel.UserId, types ...aclcore.ScopeType) ([]*aclcore.Scope, error)
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

	resourceIdsByOwner := make(map[ownermodel.Owner][]string)
	for _, id := range albumIds {
		before, _ := resourceIdsByOwner[id.Owner]
		resourceIdsByOwner[id.Owner] = append(before, id.FolderName.String())
	}

	for owner, resourceIds := range resourceIdsByOwner {
		scopes, err := f.ScopeRepository.ListScopesByOwner(ctx, owner, aclcore.MainOwnerScope, aclcore.AlbumVisitorScope, aclcore.AlbumContributorScope)
		if err != nil {
			return nil, err
		}

		for _, scope := range scopes {
			if scope.Type == aclcore.MainOwnerScope {
				for _, albumId := range albumIds {
					if albumId.Owner == scope.ResourceOwner {
						availability := catalogviews.OwnerAvailability(scope.GrantedTo)

						if list, ok := grid[albumId]; ok {
							grid[albumId] = append(list, availability)
						} else {
							grid[albumId] = []catalogviews.Availability{availability}
						}
					}
				}

			} else if len(scope.ResourceId) > 0 && slices.Contains(resourceIds, scope.ResourceId) && (scope.Type == aclcore.AlbumVisitorScope || scope.Type == aclcore.AlbumContributorScope) {
				albumId := catalog.AlbumId{Owner: scope.ResourceOwner, FolderName: catalog.NewFolderName(scope.ResourceId)}
				availability := catalogviews.VisitorAvailability(scope.GrantedTo)

				if list, ok := grid[albumId]; ok {
					grid[albumId] = append(list, availability)
				} else {
					grid[albumId] = []catalogviews.Availability{availability}
				}
			}
		}
	}

	return grid, nil
}
