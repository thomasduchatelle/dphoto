package catalogacl

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type UnShareAlbumCase struct {
	RevokeScopeRepository aclcore.ScopeWriter
	Observers             []AlbumUnSharedObserver
}

type AlbumUnSharedObserver interface {
	AlbumUnShared(ctx context.Context, albumId catalog.AlbumId, userEmail usermodel.UserId) error
}

func (u *UnShareAlbumCase) StopSharingAlbum(albumId catalog.AlbumId, email usermodel.UserId) error {
	err := u.RevokeScopeRepository.DeleteScopes(aclcore.ScopeId{
		Type:          aclcore.AlbumVisitorScope,
		GrantedTo:     email,
		ResourceOwner: albumId.Owner,
		ResourceId:    albumId.FolderName.String(),
	})
	if err != nil {
		return err
	}

	for _, observer := range u.Observers {
		err = observer.AlbumUnShared(context.TODO(), albumId, email)
		if err != nil {
			return err
		}
	}
	return nil
}
