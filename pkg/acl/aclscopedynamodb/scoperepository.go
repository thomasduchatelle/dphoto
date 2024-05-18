package aclscopedynamodb

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
)

type GrantRepository interface {
	aclcore.ScopeReadRepository
	aclcore.ScopeWriter
	aclcore.IdentityQueriesScopeRepository

	aclcore.ScopesReader        // TODO remove: duplicates aclcore.ScopeReadRepository
	aclcore.ReverseScopesReader // TODO remove: duplicates aclcore.ScopeReadRepository
}

func New(client *dynamodb.Client, tableName string) (*Repository, error) {
	return &Repository{
		client: client,
		table:  tableName,
	}, nil
}

func Must(repository GrantRepository, err error) GrantRepository {
	if err != nil {
		panic(err)
	}
	return repository
}

type Repository struct {
	client *dynamodb.Client
	table  string
}
