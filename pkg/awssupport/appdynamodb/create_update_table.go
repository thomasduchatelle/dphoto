package appdynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"time"
)

const (
	tableVersion = "2.1" // tableVersion should be bumped manually when schema is updated
)

// CreateTableIfNecessary creates the table if it doesn't exist ; or update it.
func CreateTableIfNecessary(table string, db *dynamodb.DynamoDB, localDynamodb bool) error {
	mdc := log.WithFields(log.Fields{
		"TableBackup":  table,
		"TableVersion": tableVersion,
	})
	mdc.Debugf("CreateTableIfNecessary > describe table '%s'", table)

	s := aws.String(dynamodb.ScalarAttributeTypeS)
	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("PK"), AttributeType: s},
			{AttributeName: aws.String("SK"), AttributeType: s},
			{AttributeName: aws.String("AlbumIndexPK"), AttributeType: s},
			{AttributeName: aws.String("AlbumIndexSK"), AttributeType: s},
			{AttributeName: aws.String("LocationId"), AttributeType: s},
			{AttributeName: aws.String("LocationKeyPrefix"), AttributeType: s},
			{AttributeName: aws.String("ResourceOwner"), AttributeType: s},
			{AttributeName: aws.String("AbsoluteExpiryTime"), AttributeType: s},
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
			{
				IndexName: aws.String("ReverseLocationIndex"), // from 'archivedynamo' extension
				KeySchema: []*dynamodb.KeySchemaElement{
					{AttributeName: aws.String("LocationKeyPrefix"), KeyType: aws.String(dynamodb.KeyTypeHash)},
					{AttributeName: aws.String("LocationId"), KeyType: aws.String(dynamodb.KeyTypeRange)},
				},
				Projection: &dynamodb.Projection{ProjectionType: aws.String(dynamodb.ProjectionTypeAll)},
			},
			{
				IndexName: aws.String("ReverseGrantIndex"), // from 'acl' extension
				KeySchema: []*dynamodb.KeySchemaElement{
					{AttributeName: aws.String("ResourceOwner"), KeyType: aws.String(dynamodb.KeyTypeHash)},
					{AttributeName: aws.String("SK"), KeyType: aws.String(dynamodb.KeyTypeRange)},
				},
				Projection: &dynamodb.Projection{ProjectionType: aws.String(dynamodb.ProjectionTypeAll)},
			},
			{
				IndexName: aws.String("RefreshTokenExpiration"), // from 'acl' extension
				KeySchema: []*dynamodb.KeySchemaElement{
					{AttributeName: aws.String("SK"), KeyType: aws.String(dynamodb.KeyTypeHash)},
					{AttributeName: aws.String("AbsoluteExpiryTime"), KeyType: aws.String(dynamodb.KeyTypeRange)},
				},
				Projection: &dynamodb.Projection{ProjectionType: aws.String(dynamodb.ProjectionTypeInclude), NonKeyAttributes: []*string{aws.String("PK")}},
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("PK"), KeyType: aws.String(dynamodb.KeyTypeHash)},
			{AttributeName: aws.String("SK"), KeyType: aws.String(dynamodb.KeyTypeRange)},
		},
		TableName: &table,
	}

	if !localDynamodb {
		// Localstack dynamodb doesn't support tags
		createTableInput.Tags = []*dynamodb.Tag{
			{Key: aws.String("Version"), Value: aws.String(tableVersion)},
			{Key: aws.String("LastUpdated"), Value: aws.String(time.Now().Format("2006-01-02 15:04:05"))},
		}
	}

	return dynamoutils.CreateOrUpdateTable(context.Background(), &dynamoutils.CreateOrUpdateTableInput{
		Client:     db,
		TableName:  table,
		Definition: createTableInput,
	})
}
