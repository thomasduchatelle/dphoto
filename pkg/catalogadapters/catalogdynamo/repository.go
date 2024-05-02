package catalogdynamo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

type RepositoryContract interface {
	catalog.RepositoryAdapter
	catalog.TransferMediasPort
}

type Repository struct {
	client *dynamodb.Client
	table  string
}

// NewRepository creates the repository and connect to the database
func NewRepository(client *dynamodb.Client, tableName string) *Repository {
	return &Repository{
		client: client,
		table:  tableName,
	}
}

// Must panics if there is an error
func Must(repository catalog.RepositoryAdapter, err error) catalog.RepositoryAdapter {
	if err != nil {
		panic(err)
	}

	return repository
}
