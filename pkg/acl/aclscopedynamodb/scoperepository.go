package aclscopedynamodb

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
)

type GrantRepository interface {
	aclcore.ScopesReader
	aclcore.ReverseScopesReader
	aclcore.ScopeWriter
}

func New(sess *session.Session, tableName string) (GrantRepository, error) {
	return &repository{
		db:    dynamodb.New(sess),
		table: tableName,
	}, nil
}

func Must(repository GrantRepository, err error) GrantRepository {
	if err != nil {
		panic(err)
	}
	return repository
}

type repository struct {
	db    *dynamodb.DynamoDB
	table string
}
