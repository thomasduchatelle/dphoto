package catalogdynamo

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/dynamoutils"
	"time"
)

const (
	tableVersion         = "2.0" // tableVersion should be bumped manually when schema is updated
	DynamoWriteBatchSize = 25
	DynamoReadBatchSize  = 100
)

type rep struct {
	db            *dynamodb.DynamoDB
	table         string
	localDynamodb bool // localDynamodb is set to true to disable some feature - not available on localstack - like tagging
}

// NewRepository creates the repository and connect to the database
func NewRepository(awsSession *session.Session, tableName string) (catalog.RepositoryAdapter, error) {
	rep := &rep{
		db:            dynamodb.New(awsSession),
		table:         tableName,
		localDynamodb: false,
	}

	err := rep.CreateTableIfNecessary()
	return rep, err
}

// Must panics if there is an error
func Must(repository catalog.RepositoryAdapter, err error) catalog.RepositoryAdapter {
	if err != nil {
		panic(err)
	}

	return repository
}

// CreateTableIfNecessary creates the table if it doesn't exists ; or update it.
func (r *rep) CreateTableIfNecessary() error {
	mdc := log.WithFields(log.Fields{
		"TableBackup":  r.table,
		"TableVersion": tableVersion,
	})
	mdc.Debugf("CreateTableIfNecessary > describe table '%s'", r.table)

	s := aws.String(dynamodb.ScalarAttributeTypeS)
	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("PK"), AttributeType: s},
			{AttributeName: aws.String("SK"), AttributeType: s},
			{AttributeName: aws.String("AlbumIndexPK"), AttributeType: s},
			{AttributeName: aws.String("AlbumIndexSK"), AttributeType: s},
		},
		BillingMode: aws.String(dynamodb.BillingModePayPerRequest),
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("AlbumIndex"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{AttributeName: aws.String("AlbumIndexPK"), KeyType: aws.String(dynamodb.KeyTypeHash)},
					{AttributeName: aws.String("AlbumIndexSK"), KeyType: aws.String(dynamodb.KeyTypeRange)},
				},
				Projection: &dynamodb.Projection{ProjectionType: aws.String(dynamodb.ProjectionTypeAll)},
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("PK"), KeyType: aws.String(dynamodb.KeyTypeHash)},
			{AttributeName: aws.String("SK"), KeyType: aws.String(dynamodb.KeyTypeRange)},
		},
		TableName: &r.table,
	}

	if !r.localDynamodb {
		// Localstack dynamodb doesn't support tags
		createTableInput.Tags = []*dynamodb.Tag{
			{Key: aws.String("Version"), Value: aws.String(tableVersion)},
			{Key: aws.String("LastUpdated"), Value: aws.String(time.Now().Format("2006-01-02 15:04:05"))},
		}
	}

	return dynamoutils.CreateOrUpdateTable(context.Background(), &dynamoutils.CreateOrUpdateTableInput{
		Client:     r.db,
		TableName:  r.table,
		Definition: createTableInput,
	})
}
