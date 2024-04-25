package appdynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"time"
)

const (
	tableVersion = "2.1" // tableVersion should be bumped manually when schema is updated
)

// CreateTableIfNecessary creates the table if it doesn't exist ; or update it.
func CreateTableIfNecessary(ctx context.Context, table string, client *dynamodb.Client, localDynamodb bool) error {
	mdc := log.WithFields(log.Fields{
		"TableBackup":  table,
		"TableVersion": tableVersion,
	})
	mdc.Debugf("CreateTableIfNecessary > describe table '%s'", table)

	var secondaryIndexProvisionedThroughput *types.ProvisionedThroughput
	if localDynamodb {
		secondaryIndexProvisionedThroughput = &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		}
	}

	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("PK"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("SK"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("AlbumIndexPK"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("AlbumIndexSK"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("LocationId"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("LocationKeyPrefix"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("ResourceOwner"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("AbsoluteExpiryTime"), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("PK"), KeyType: types.KeyTypeHash},
			{AttributeName: aws.String("SK"), KeyType: types.KeyTypeRange},
		},
		TableName: &table,
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("AlbumIndex"),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("AlbumIndexPK"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("AlbumIndexSK"), KeyType: types.KeyTypeRange},
				},
				Projection:            &types.Projection{ProjectionType: types.ProjectionTypeAll},
				ProvisionedThroughput: secondaryIndexProvisionedThroughput,
			},
			{
				IndexName: aws.String("ReverseLocationIndex"), // from 'archivedynamo' extension
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("LocationKeyPrefix"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("LocationId"), KeyType: types.KeyTypeRange},
				},
				Projection:            &types.Projection{ProjectionType: types.ProjectionTypeAll},
				ProvisionedThroughput: secondaryIndexProvisionedThroughput,
			},
			{
				IndexName: aws.String("ReverseGrantIndex"), // from 'acl' extension
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("ResourceOwner"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("SK"), KeyType: types.KeyTypeRange},
				},
				Projection:            &types.Projection{ProjectionType: types.ProjectionTypeAll},
				ProvisionedThroughput: secondaryIndexProvisionedThroughput,
			},
			{
				IndexName: aws.String("RefreshTokenExpiration"), // from 'acl' extension
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("SK"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("AbsoluteExpiryTime"), KeyType: types.KeyTypeRange},
				},
				Projection:            &types.Projection{ProjectionType: types.ProjectionTypeInclude, NonKeyAttributes: []string{"PK"}},
				ProvisionedThroughput: secondaryIndexProvisionedThroughput,
			},
		},
		ProvisionedThroughput: secondaryIndexProvisionedThroughput,
	}

	if !localDynamodb {
		// Localstack dynamodb doesn't support tags
		createTableInput.Tags = []types.Tag{
			{Key: aws.String("Version"), Value: aws.String(tableVersion)},
			{Key: aws.String("LastUpdated"), Value: aws.String(time.Now().Format("2006-01-02 15:04:05"))},
		}
	}

	return dynamoutils.CreateOrUpdateTable(ctx, &dynamoutils.CreateOrUpdateTableInput{
		Client:     client,
		TableName:  table,
		Definition: createTableInput,
	})
}
