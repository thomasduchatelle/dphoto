package catalogacl

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type FindAlbumPort interface {
	FindAlbum(ctx context.Context, albumId catalog.AlbumId) (*catalog.Album, error)
}

type AlbumSharedObserver interface {
	AlbumShared(ctx context.Context, albumId catalog.AlbumId, userEmail usermodel.UserId) error
}

type ShareAlbumCase struct {
	ScopeWriter   aclcore.ScopeWriter
	FindAlbumPort FindAlbumPort
	Observers     []AlbumSharedObserver
}

func (s *ShareAlbumCase) ShareAlbumWith(albumId catalog.AlbumId, userEmail usermodel.UserId, scope aclcore.ScopeType) error {
	ctx := context.TODO()
	if scope != aclcore.AlbumVisitorScope && scope != aclcore.AlbumContributorScope {
		return errors.Errorf("'%s' scope is not allowed for album shring.", scope)
	}

	_, err := s.FindAlbumPort.FindAlbum(ctx, albumId)
	if err != nil {
		return errors.Wrapf(err, "album %s cannot be shared to %s", albumId, userEmail) // it can be a catalog.AlbumNotFoundError
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
		err = observer.AlbumShared(ctx, albumId, userEmail)
		if err != nil {
			return err
		}
	}

	return nil
}
