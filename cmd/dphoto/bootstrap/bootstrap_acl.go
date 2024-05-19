package bootstrap

import (
	"context"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

func init() {
	config.Listen(func(cfg config.Config) {
		ctx := context.TODO()

		repository := pkgfactory.AclRepository(ctx)
		createUser := &aclcore.CreateUser{
			ScopesReader: repository,
			ScopeWriter:  repository,
		}
		cmd.CreateUserCase = createUser.CreateUser

		sharedAlbum := pkgfactory.AclCatalogShare(ctx)
		cmd.ShareAlbumCase = sharedAlbum.ShareAlbumWith

		unSharedAlbum := pkgfactory.AclCatalogUnShare(ctx)
		cmd.UnShareAlbumCase = unSharedAlbum.StopSharingAlbum
	})
}
