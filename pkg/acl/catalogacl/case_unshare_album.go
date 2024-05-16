package catalogacl

import (
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type UnShareAlbumCase struct {
	RevokeScopeRepository aclcore.ScopeWriter
}

func (u *UnShareAlbumCase) StopSharingAlbum(albumId catalog.AlbumId, email usermodel.UserId) error {
	return u.RevokeScopeRepository.DeleteScopes(aclcore.ScopeId{
		Type:          aclcore.AlbumVisitorScope,
		GrantedTo:     email,
		ResourceOwner: albumId.Owner,
		ResourceId:    albumId.FolderName.String(),
	})
}
