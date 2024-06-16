package pkgfactory

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclscopedynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/pkg/singletons"
)

func AclRepository(ctx context.Context) *aclscopedynamodb.Repository {
	return singletons.MustSingleton(func() (*aclscopedynamodb.Repository, error) {
		return aclscopedynamodb.New(AWSFactory(ctx).GetDynamoDBClient(), AWSNames.DynamoDBName())
	})
}

func AclQueries(ctx context.Context) *aclcore.ScopeQueries {
	return singletons.MustSingleton(func() (*aclcore.ScopeQueries, error) {
		return &aclcore.ScopeQueries{
			ScopeReadRepository: AclRepository(ctx),
		}, nil
	})
}

func AclCatalogShare(ctx context.Context) *catalogacl.ShareAlbumCase {
	return &catalogacl.ShareAlbumCase{
		ScopeWriter:   AclRepository(ctx),
		FindAlbumPort: AlbumQueries(ctx),
		Observers: []catalogacl.AlbumSharedObserver{
			CommandHandlerAlbumSize(ctx),
		},
	}
}

func AclCatalogUnShare(ctx context.Context) *catalogacl.UnShareAlbumCase {
	return &catalogacl.UnShareAlbumCase{
		RevokeScopeRepository: AclRepository(ctx),
		Observers: []catalogacl.AlbumUnSharedObserver{
			CommandHandlerAlbumSize(ctx),
		},
	}
}

func AclCatalogAuthoriser(ctx context.Context) *catalogacl.CatalogAuthorizer {
	return &catalogacl.CatalogAuthorizer{
		HasPermissionPort: AclQueries(ctx),
	}
}
