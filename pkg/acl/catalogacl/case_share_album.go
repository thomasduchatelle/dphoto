package catalogacl

import (
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

type ShareAlbumCatalogPort interface {
	FindAlbum(owner, folderName string) (*catalog.Album, error)
}

type ShareAlbumCase struct {
	ScopeWriter aclcore.ScopeWriter
	CatalogPort ShareAlbumCatalogPort
}

func (s *ShareAlbumCase) ShareAlbumWith(owner, folderName, userEmail string, scope aclcore.ScopeType) error {
	if scope != aclcore.AlbumVisitorScope && scope != aclcore.AlbumContributorScope {
		return errors.Errorf("'%s' scope is not allowed for album shring.", scope)
	}

	_, err := s.CatalogPort.FindAlbum(owner, folderName)
	if err != nil {
		return err // it can be a catalog.AlbumNotFoundError
	}

	return s.ScopeWriter.SaveIfNewScope(aclcore.Scope{
		Type:          scope,
		GrantedAt:     aclcore.TimeFunc(),
		GrantedTo:     userEmail,
		ResourceOwner: owner,
		ResourceId:    folderName,
	})
}
