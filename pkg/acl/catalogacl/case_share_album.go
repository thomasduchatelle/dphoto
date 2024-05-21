package catalogacl

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type ShareAlbumCatalogPort interface {
	FindAlbum(albumId catalog.AlbumId) (*catalog.Album, error)
}

type AlbumSharedObserver interface {
	AlbumShared(ctx context.Context, albumId catalog.AlbumId, userEmail usermodel.UserId) error
}

type ShareAlbumCase struct {
	ScopeWriter aclcore.ScopeWriter
	CatalogPort ShareAlbumCatalogPort
	Observers   []AlbumSharedObserver
}

func (s *ShareAlbumCase) ShareAlbumWith(albumId catalog.AlbumId, userEmail usermodel.UserId, scope aclcore.ScopeType) error {
	if scope != aclcore.AlbumVisitorScope && scope != aclcore.AlbumContributorScope {
		return errors.Errorf("'%s' scope is not allowed for album shring.", scope)
	}

	_, err := s.CatalogPort.FindAlbum(albumId)
	if err != nil {
		return err // it can be a catalog.AlbumNotFoundError
	}

	err = s.ScopeWriter.SaveIfNewScope(aclcore.Scope{
		Type:          scope,
		GrantedAt:     aclcore.TimeFunc(),
		GrantedTo:     userEmail,
		ResourceOwner: albumId.Owner,
		ResourceId:    albumId.FolderName.String(),
	})
	if err != nil {
		return err
	}

	for _, observer := range s.Observers {
		err = observer.AlbumShared(context.TODO(), albumId, userEmail)
		if err != nil {
			return err
		}
	}

	return nil
}
