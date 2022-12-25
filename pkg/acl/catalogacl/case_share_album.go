package catalogacl

import (
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

func (s *ShareAlbumCase) ShareAlbumWith(owner, folderName, userEmail string) error {
	_, err := s.CatalogPort.FindAlbum(owner, folderName)
	if err != nil {
		return err // it can be a catalog.NotFoundError
	}

	return s.ScopeWriter.SaveIfNewScope(aclcore.Scope{
		Type:          aclcore.AlbumVisitorScope,
		GrantedAt:     aclcore.TimeFunc(),
		GrantedTo:     userEmail,
		ResourceOwner: owner,
		ResourceId:    folderName,
	})
}
