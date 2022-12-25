package acldynamodb

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogdynamo"
)

type GrantRepository interface {
	aclcore.ScopesReader
	aclcore.ReverseScopesReader
	aclcore.ScopeWriter
}

func New(sess *session.Session, tableName string, createTable bool) (GrantRepository, error) {
	if createTable {
		catalogRepository, err := catalogdynamo.NewRepository(sess, tableName)
		if err != nil {
			return nil, err
		}
		if tableCreator, ok := catalogRepository.(*catalogdynamo.Repository); ok {
			err = tableCreator.CreateTableIfNecessary()
			if err != nil {
				return nil, err
			}
		}
	}

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
