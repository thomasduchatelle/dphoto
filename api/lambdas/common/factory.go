package common

import (
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

type catalogPort struct {
}

func (c catalogPort) FindAlbum(owner, folderName string) (*catalog.Album, error) {
	return catalog.FindAlbum(owner, folderName)
}

func GetShareAlbumCase() *catalogacl.ShareAlbumCase {
	return &catalogacl.ShareAlbumCase{
		ScopeWriter: grantRepository,
		CatalogPort: new(catalogPort),
	}
}

func GetUnShareAlbumCase() *catalogacl.UnShareAlbumCase {
	return &catalogacl.UnShareAlbumCase{
		RevokeScopeRepository: grantRepository,
	}
}

func GetIdentityQueries() *aclcore.IdentityQueries {
	return &aclcore.IdentityQueries{
		IdentityRepository: getIdentityDetailsStore(),
		ScopeRepository:    grantRepository,
	}
}
