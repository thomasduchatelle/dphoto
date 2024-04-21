package bootstrap

import (
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclscopedynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

type catalogPort struct {
}

func (c catalogPort) FindAlbum(owner, folderName string) (*catalog.Album, error) {
	return catalog.FindAlbum(owner, folderName)
}

func init() {
	config.Listen(func(cfg config.Config) {
		repository := aclscopedynamodb.Must(aclscopedynamodb.New(cfg.GetAWSV2Config(), cfg.GetString(config.CatalogDynamodbTable)))
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
