package bootstrap

import (
	"context"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

func init() {
	config.Listen(func(cfg config.Config) {
		ctx := context.TODO()

		createUser := pkgfactory.CreateUserCase(ctx)
		cmd.CreateUserCase = createUser.CreateUser

		sharedAlbum := pkgfactory.AclCatalogShare(ctx)
		cmd.ShareAlbumCase = sharedAlbum.ShareAlbumWith

		unSharedAlbum := pkgfactory.AclCatalogUnShare(ctx)
		cmd.UnShareAlbumCase = unSharedAlbum.StopSharingAlbum
	})
}
