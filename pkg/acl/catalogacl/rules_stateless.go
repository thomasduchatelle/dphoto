// Package catalogacl contains the logic to authorise access to catalog resources based on user requesting it.
package catalogacl

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type ScopeRepository interface {
	aclcore.ScopesReader
	aclcore.ReverseScopesReader
}

type MediaAlbumResolver interface {
	FindAlbumOfMedia(owner ownermodel.Owner, mediaId catalog.MediaId) (catalog.AlbumId, error)
}

type CatalogRules interface {
	Owner() (*ownermodel.Owner, error)
	SharedWithUserAlbum() ([]catalog.AlbumId, error)
	SharedByUserGrid(owner ownermodel.Owner) (map[string]map[usermodel.UserId]aclcore.ScopeType, error)
	CanListMediasFromAlbum(id catalog.AlbumId) error
	CanReadMedia(owner ownermodel.Owner, id catalog.MediaId) error

	CanManageAlbum(id catalog.AlbumId) error
}

// NewCatalogRules creates an adapter catalogacl -> aclcore which will always request DB layer
func NewCatalogRules(repository ScopeRepository, mediaAlbumResolver MediaAlbumResolver, email usermodel.UserId) CatalogRules {
	return &rules{
		CoreRules: aclcore.CoreRules{
			ScopeReader: repository,
			Email:       email,
		},
		email:              email,
		mediaAlbumResolver: mediaAlbumResolver,
		scopeRepository:    repository,
	}
}

type rules struct {
	aclcore.CoreRules
	email              usermodel.UserId
	mediaAlbumResolver MediaAlbumResolver
	scopeRepository    ScopeRepository
}

func (r *rules) SharedWithUserAlbum() ([]catalog.AlbumId, error) {
	shared, err := r.scopeRepository.ListUserScopes(r.email, aclcore.AlbumVisitorScope)

	var albums []catalog.AlbumId
	for _, share := range shared {
		albums = append(albums, catalog.AlbumId{Owner: share.ResourceOwner, FolderName: catalog.NewFolderName(share.ResourceId)})
	}

	return albums, err
}

func (r *rules) SharedByUserGrid(owner ownermodel.Owner) (map[string]map[usermodel.UserId]aclcore.ScopeType, error) {
	scopes, err := r.scopeRepository.ListOwnerScopes(owner, aclcore.AlbumVisitorScope, aclcore.AlbumContributorScope)

	if err != nil || len(scopes) == 0 {
		return nil, err
	}

	grid := make(map[string]map[usermodel.UserId]aclcore.ScopeType)
	for _, scope := range scopes {
		list, ok := grid[scope.ResourceId]
		if !ok || len(list) == 0 {
			list = make(map[usermodel.UserId]aclcore.ScopeType)
			grid[scope.ResourceId] = list
		}

		list[scope.GrantedTo] = scope.Type
	}

	return grid, nil
}

func (r *rules) CanListMediasFromAlbum(albumId catalog.AlbumId) error {
	scopes, err := r.scopeRepository.FindScopesById(
		aclcore.ScopeId{
			Type:          aclcore.MainOwnerScope,
			GrantedTo:     r.email,
			ResourceOwner: albumId.Owner,
		},
		aclcore.ScopeId{
			Type:          aclcore.AlbumVisitorScope,
			GrantedTo:     r.email,
			ResourceOwner: albumId.Owner,
			ResourceId:    albumId.FolderName.String(),
		},
	)
	if err != nil {
		return err
	}
	if len(scopes) == 0 {
		return errors.Wrapf(aclcore.AccessForbiddenError, "listing medias in %s denied.", albumId)
	}

	return nil
}

func (r *rules) CanReadMedia(owner ownermodel.Owner, mediaId catalog.MediaId) error {
	albumId, err := r.mediaAlbumResolver.FindAlbumOfMedia(owner, mediaId)
	if err != nil {
		return errors.Wrapf(aclcore.AccessForbiddenError, err.Error())
	}

	scopes, err := r.scopeRepository.FindScopesById(
		aclcore.ScopeId{
			Type:          aclcore.MainOwnerScope,
			GrantedTo:     r.email,
			ResourceOwner: owner,
		},
		aclcore.ScopeId{
			Type:          aclcore.AlbumVisitorScope,
			GrantedTo:     r.email,
			ResourceOwner: owner,
			ResourceId:    albumId.FolderName.String(),
		},
		aclcore.ScopeId{
			Type:          aclcore.MediaVisitorScope,
			GrantedTo:     r.email,
			ResourceOwner: owner,
			ResourceId:    mediaId.Value(),
		},
	)
	if err != nil {
		return err
	}
	if len(scopes) == 0 {
		return errors.Wrapf(aclcore.AccessForbiddenError, "reading media %s/%s has been denied.", owner, mediaId)
	}

	return nil
}

func (r *rules) CanManageAlbum(albumId catalog.AlbumId) error {
	scopes, err := r.scopeRepository.FindScopesById(
		aclcore.ScopeId{Type: aclcore.MainOwnerScope, GrantedTo: r.email, ResourceOwner: albumId.Owner},
	)
	if err != nil {
		return err
	}

	if len(scopes) == 0 {
		return errors.Wrapf(aclcore.AccessForbiddenError, fmt.Sprintf("%s is not allowed to managed album %s", r.email, albumId))
	}
	return nil
}
