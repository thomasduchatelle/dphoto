package catalogacl

import (
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type ShareAlbumCatalogPort interface {
	FindAlbum(albumId catalog.AlbumId) (*catalog.Album, error)
}

type ShareAlbumCase struct {
	ScopeWriter aclcore.ScopeWriter
	CatalogPort ShareAlbumCatalogPort
}

func (s *ShareAlbumCase) ShareAlbumWith(albumId catalog.AlbumId, userEmail usermodel.UserId, scope aclcore.ScopeType) error {
	if scope != aclcore.AlbumVisitorScope && scope != aclcore.AlbumContributorScope {
		return errors.Errorf("'%s' scope is not allowed for album shring.", scope)
	}

	_, err := s.CatalogPort.FindAlbum(albumId)
	if err != nil {
		return err // it can be a catalog.AlbumNotFoundError
	}

	return s.ScopeWriter.SaveIfNewScope(aclcore.Scope{
		Type:          scope,
		GrantedAt:     aclcore.TimeFunc(),
		GrantedTo:     userEmail,
		ResourceOwner: albumId.Owner,
		ResourceId:    albumId.FolderName.String(),
	})
}
