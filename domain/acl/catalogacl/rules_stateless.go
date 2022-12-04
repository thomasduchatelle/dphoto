// Package catalogacl contains the logic to authorise access to catalog resources based on user requesting it.
package catalogacl

import (
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
)

type ScopeRepository interface {
	aclcore.ScopesReader
	aclcore.ReverseScopesReader
}

type MediaAlbumResolver interface {
	FindAlbumOfMedia(owner, mediaId string) (string, error)
}

type CatalogRules interface {
	Owner() (string, error)
	SharedWithUserAlbum() ([]catalog.AlbumId, error)
	SharedByUserGrid(owner string) (map[string][]string, error)
	CanListMediasFromAlbum(owner string, folderName string) error
	CanReadMedia(owner string, id string) error
}

// NewCatalogRules creates an adapter catalogacl -> aclcore which will always request DB layer
func NewCatalogRules(repository ScopeRepository, mediaAlbumResolver MediaAlbumResolver, email string) CatalogRules {
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
	email              string
	mediaAlbumResolver MediaAlbumResolver
	scopeRepository    ScopeRepository
}

func (r *rules) SharedWithUserAlbum() ([]catalog.AlbumId, error) {
	shared, err := r.scopeRepository.ListUserScopes(r.email, aclcore.AlbumVisitorScope)

	var albums []catalog.AlbumId
	for _, share := range shared {
		albums = append(albums, catalog.AlbumId{
			Owner:      share.ResourceOwner,
			FolderName: share.ResourceId,
		})
	}

	return albums, err
}

func (r *rules) SharedByUserGrid(owner string) (map[string][]string, error) {
	scopes, err := r.scopeRepository.ListOwnerScopes(owner, aclcore.AlbumVisitorScope)

	if err != nil || len(scopes) == 0 {
		return nil, err
	}

	grid := make(map[string][]string)
	for _, scope := range scopes {
		list, _ := grid[scope.ResourceId]
		grid[scope.ResourceId] = append(list, scope.GrantedTo)
	}

	return grid, nil
}

func (r *rules) CanListMediasFromAlbum(owner string, folderName string) error {
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
			ResourceId:    folderName,
		},
	)
	if err != nil {
		return err
	}
	if len(scopes) == 0 {
		return errors.Wrapf(aclcore.AccessForbiddenError, "listing medias in %s/%s denied.", owner, folderName)
	}

	return nil
}

func (r *rules) CanReadMedia(owner string, mediaId string) error {
	folderName, err := r.mediaAlbumResolver.FindAlbumOfMedia(owner, mediaId)
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
			ResourceId:    folderName,
		},
		aclcore.ScopeId{
			Type:          aclcore.MediaVisitorScope,
			GrantedTo:     r.email,
			ResourceOwner: owner,
			ResourceId:    mediaId,
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
