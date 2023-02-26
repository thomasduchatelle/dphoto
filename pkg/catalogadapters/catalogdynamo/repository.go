package catalogdynamo

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

const (
	DynamoWriteBatchSize = 25
	DynamoReadBatchSize  = 100
)

type Repository struct {
	db            *dynamodb.DynamoDB
	table         string
	localDynamodb bool // localDynamodb is set to true to disable some feature - not available on localstack - like tagging
}

// NewRepository creates the repository and connect to the database
func NewRepository(awsSession *session.Session, tableName string) (catalog.RepositoryAdapter, error) {
	rep := &Repository{
		db:            dynamodb.New(awsSession),
		table:         tableName,
		localDynamodb: false,
	}

	return rep, nil
}

// Must panics if there is an error
func Must(repository catalog.RepositoryAdapter, err error) catalog.RepositoryAdapter {
	if err != nil {
		panic(err)
	}

	return repository
}
