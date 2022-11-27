// Package catalogacllogic contains the logic to authorise access to catalog resources based on user requesting it.
package catalogacllogic

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

// NewAccessControlAdapter creates an adapter catalogacl -> accesscontrol which will always request DB layer
func NewAccessControlAdapter(repository ScopeRepository, email string) catalogacl.AccessControlAdapter {
	return &adapter{
		email:           email,
		scopeRepository: repository,
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

func (a *adapter) CanListMediasFromAlbum(owner string, folderName string) error {
	scopes, err := a.scopeRepository.FindScopesById(
		accesscontrol.ScopeId{
			Type:          accesscontrol.MainOwnerScope,
			GrantedTo:     a.email,
			ResourceOwner: owner,
		},
		accesscontrol.ScopeId{
			Type:          accesscontrol.AlbumVisitorScope,
			GrantedTo:     a.email,
			ResourceOwner: owner,
			ResourceId:    folderName,
		},
	)
	if err != nil {
		return err
	}
	if len(scopes) == 0 {
		return errors.Wrapf(accesscontrol.AccessForbiddenError, "listing medias in %s/%s denied.", owner, folderName)
	}

	return nil
}
