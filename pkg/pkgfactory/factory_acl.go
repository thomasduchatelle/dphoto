package pkgfactory

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclscopedynamodb"
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
