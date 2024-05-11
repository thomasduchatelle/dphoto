package catalogdynamo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

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
