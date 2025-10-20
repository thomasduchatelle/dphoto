package pkgfactory

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcognitoadapter"
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

func CognitoRepository(ctx context.Context) aclcore.CognitoRepository {
	return singletons.MustSingletonKey("CognitoRepository", func() (aclcore.CognitoRepository, error) {
		userPoolId := AWSNames.CognitoUserPoolId()
		if userPoolId == "" {
			log.Warn("COGNITO_USER_POOL_ID is not configured, Cognito user creation will be skipped")
			return nil, nil
		}

		return aclcognitoadapter.New(AWSFactory(ctx).GetCfg(), userPoolId)
	})
}

func CreateUserCase(ctx context.Context) *aclcore.CreateUser {
	return singletons.MustSingletonKey("CreateUserCase", func() (*aclcore.CreateUser, error) {
		repository := AclRepository(ctx)
		return &aclcore.CreateUser{
			ScopesReader:      repository,
			ScopeWriter:       repository,
			CognitoRepository: CognitoRepository(ctx),
		}, nil
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
		ScopeWriter:       AclRepository(ctx),
		FindAlbumPort:     AlbumQueries(ctx),
		CognitoRepository: CognitoRepository(ctx),
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
		HasPermissionPort:  AclQueries(ctx),
		CatalogQueriesPort: CatalogMediaQueries(ctx),
	}
}
