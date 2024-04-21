package aclscopedynamodb

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
)

type GrantRepository interface {
	aclcore.ScopesReader
	aclcore.ReverseScopesReader
	aclcore.ScopeWriter
	aclcore.IdentityQueriesScopeRepository
}

func New(cfg aws.Config, tableName string) (GrantRepository, error) {
	return &repository{
		client: dynamodb.NewFromConfig(cfg),
		table:  tableName,
	}, nil
}

func Must(repository GrantRepository, err error) GrantRepository {
	if err != nil {
		panic(err)
	}
	return repository
}

type repository struct {
	client *dynamodb.Client
	table  string
}
