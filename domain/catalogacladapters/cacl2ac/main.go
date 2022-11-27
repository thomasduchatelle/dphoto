package cacl2ac

import (
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/accesscontrol"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogacl"
)

type ScopeRepository interface {
	accesscontrol.ScopesReader
	accesscontrol.ReverseScopesReader
}

func NewAccessControlAdapter(repository ScopeRepository, email string) catalogacl.AccessControlAdapter {
	return &adapter{
		email:           email,
		scopeRepository: repository,
	}
}

func NewAccessControlAdapterFromToken(repository ScopeRepository, claims accesscontrol.Claims) catalogacl.AccessControlAdapter {
	delegate := NewAccessControlAdapter(repository, claims.Subject)
	return &accessTokenAdapter{
		AccessControlAdapter: delegate,
		owner:                claims.Owner,
	}
}

type adapter struct {
	email           string
	scopeRepository ScopeRepository
}

func (a *adapter) Owner() (string, error) {
	scopes, err := a.scopeRepository.ListUserScopes(a.email, accesscontrol.MainOwnerScope)
	if err != nil {
		return "", err
	}

	if len(scopes) == 0 {
		return "", errors.Errorf("%s is not a main user.", a.email)
	}

	return scopes[0].ResourceOwner, nil
}

func (a *adapter) SharedWithUserAlbum() ([]catalog.AlbumId, error) {
	shared, err := a.scopeRepository.ListUserScopes(a.email, accesscontrol.AlbumVisitorScope)

	var albums []catalog.AlbumId
	for _, share := range shared {
		albums = append(albums, catalog.AlbumId{
			Owner:      share.ResourceOwner,
			FolderName: share.ResourceId,
		})
	}

	return albums, err
}

func (a *adapter) SharedByUserGrid(owner string) (map[string][]string, error) {
	scopes, err := a.scopeRepository.ListOwnerScopes(owner, accesscontrol.AlbumVisitorScope)

	grid := make(map[string][]string)
	for _, scope := range scopes {
		list, _ := grid[scope.ResourceId]
		grid[scope.ResourceId] = append(list, scope.GrantedTo)
	}

	return grid, err
}

type accessTokenAdapter struct {
	catalogacl.AccessControlAdapter
	owner string
}

func (a *accessTokenAdapter) Owner() (string, error) {
	return a.owner, nil
}
