package catalogacl

import "github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"

type UnShareAlbumCase struct {
	RevokeScopeRepository aclcore.ScopeWriter
}

func (u *UnShareAlbumCase) StopSharingAlbum(owner, folderName, email string) error {
	return u.RevokeScopeRepository.DeleteScopes(aclcore.ScopeId{
		Type:          aclcore.AlbumVisitorScope,
		GrantedTo:     email,
		ResourceOwner: owner,
		ResourceId:    folderName,
	})
}
