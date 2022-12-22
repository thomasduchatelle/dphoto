package bootstrap

import (
	"github.com/thomasduchatelle/dphoto/domain/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/domain/acl/acldynamodb"
	"github.com/thomasduchatelle/dphoto/domain/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/dphoto/cmd"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
)

type catalogPort struct {
}

func (c catalogPort) FindAlbum(owner, folderName string) (*catalog.Album, error) {
	return catalog.FindAlbum(owner, folderName)
}

func init() {
	config.Listen(func(cfg config.Config) {
		repository := acldynamodb.Must(acldynamodb.New(cfg.GetAWSSession(), cfg.GetString(config.CatalogDynamodbTable), false))
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
