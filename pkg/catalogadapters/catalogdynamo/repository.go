package catalogdynamo

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

type RepositoryContract interface {
	catalog.RepositoryAdapter
	catalog.TransferMediasPort
}

type Repository struct {
	client        *dynamodb.Client
	table         string
	localDynamodb bool // localDynamodb is set to true to disable some feature - not available on localstack - like tagging
}

// NewRepository creates the repository and connect to the database
func NewRepository(cfg aws.Config, tableName string) (RepositoryContract, error) {
	rep := &Repository{
		client:        dynamodb.NewFromConfig(cfg),
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
