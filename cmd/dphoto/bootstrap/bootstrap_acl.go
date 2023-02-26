package bootstrap

import (
	cmd2 "github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd"
	config2 "github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
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
	config2.Listen(func(cfg config2.Config) {
		repository := aclscopedynamodb.Must(aclscopedynamodb.New(cfg.GetAWSSession(), cfg.GetString(config2.CatalogDynamodbTable)))
		createUser := &aclcore.CreateUser{
			ScopesReader: repository,
			ScopeWriter:  repository,
		}
		cmd2.CreateUserCase = createUser.CreateUser

		sharedAlbum := &catalogacl.ShareAlbumCase{
			ScopeWriter: repository,
			CatalogPort: new(catalogPort),
		}
		cmd2.ShareAlbumCase = sharedAlbum.ShareAlbumWith

		unSharedAlbum := &catalogacl.UnShareAlbumCase{
			RevokeScopeRepository: repository,
		}
		cmd2.UnShareAlbumCase = unSharedAlbum.StopSharingAlbum
	})
}
