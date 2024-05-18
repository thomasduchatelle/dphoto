package bootstrap

import (
	"context"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

type catalogPort struct {
}

func (c catalogPort) FindAlbum(albumId catalog.AlbumId) (*catalog.Album, error) {
	return catalog.FindAlbum(albumId)
}

func init() {
	config.Listen(func(cfg config.Config) {
		ctx := context.TODO()

		repository := pkgfactory.AclRepository(ctx)
		createUser := &aclcore.CreateUser{
			ScopesReader: repository,
			ScopeWriter:  repository,
		}
		cmd.CreateUserCase = createUser.CreateUser

		sharedAlbum := &catalogacl.ShareAlbumCase{
			ScopeWriter: repository,
			CatalogPort: new(catalogPort),
		}
		cmd.ShareAlbumCase = sharedAlbum.ShareAlbumWith

		unSharedAlbum := &catalogacl.UnShareAlbumCase{
			RevokeScopeRepository: repository,
		}
		cmd.UnShareAlbumCase = unSharedAlbum.StopSharingAlbum
	})
}
