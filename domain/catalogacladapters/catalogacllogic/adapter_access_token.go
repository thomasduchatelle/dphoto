package catalogacllogic

import (
	"github.com/thomasduchatelle/dphoto/domain/accesscontrol"
	"github.com/thomasduchatelle/dphoto/domain/catalogacl"
)

// NewAccessControlAdapterFromToken creates an adapter catalogacl -> accesscontrol with pre-authorised scopes from the Access Token
func NewAccessControlAdapterFromToken(repository ScopeRepository, claims accesscontrol.Claims) catalogacl.AccessControlAdapter {
	delegate := NewAccessControlAdapter(repository, claims.Subject)
	return &accessTokenAdapter{
		AccessControlAdapter: delegate,
		owner:                claims.Owner,
	}
}

type accessTokenAdapter struct {
	catalogacl.AccessControlAdapter
	owner string
}

func (a *accessTokenAdapter) Owner() (string, error) {
	return a.owner, nil
}

func (a *accessTokenAdapter) CanListMediasFromAlbum(owner string, folderName string) error {
	if a.owner == owner {
		return nil
	}

	return a.AccessControlAdapter.CanListMediasFromAlbum(owner, folderName)
}

func (a *accessTokenAdapter) CanReadMedia(owner string, mediaId string) error {
	if a.owner == owner {
		return nil
	}

	return a.AccessControlAdapter.CanReadMedia(owner, mediaId)
}
