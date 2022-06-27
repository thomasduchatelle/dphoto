package catalogdynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
)

const (
	versionTagName       = "Version"
	tableVersion         = "2.0" // tableVersion must be bumped manually when schema is updated
	dynamoWriteBatchSize = 25
	dynamoReadBatchSize  = 100
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

	table, err := r.db.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: &r.table,
	})

	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
		table = nil
	} else if err != nil {
		return errors.Wrap(err, "failed to find existing table")
	}

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
		Tags: []*dynamodb.Tag{
			{Key: aws.String(versionTagName), Value: aws.String(tableVersion)},
		},
	}

	if table == nil {
		mdc.Infoln("Creating dynamodb table...")
		_, err = r.db.CreateTable(createTableInput)

		return errors.Wrapf(err, "failed to create table %s", r.table)

	} else if r.localDynamodb {
		mdc.Debugln("Local dynamodb update is not supported due to lack of AWS Tag support.")

	} else {
		version, err := r.readTableVersion(table)
		if err != nil {
			return err
		}

		if version != tableVersion {
			mdc.WithFields(log.Fields{
				"Table":           r.table,
				"Version":         tableVersion,
				"PreviousVersion": version,
			}).Errorln("Updating DynamoDB table...")

			updates := generatedSecondaryUpdatesIndexes(table, createTableInput)

			_, err = r.db.UpdateTable(&dynamodb.UpdateTableInput{
				AttributeDefinitions:        createTableInput.AttributeDefinitions,
				GlobalSecondaryIndexUpdates: updates,
				TableName:                   &r.table,
			})
			if err != nil {
				return errors.Wrapf(err, "failed to update table %s", r.table)
			}

			_, err := r.db.TagResource(&dynamodb.TagResourceInput{
				ResourceArn: table.Table.TableArn,
				Tags:        createTableInput.Tags,
			})
			if err != nil {
				return errors.Wrapf(err, "failed to tag the DynamoDB table")
			}
		}
	}

	return nil
}

func generatedSecondaryUpdatesIndexes(existing *dynamodb.DescribeTableOutput, expected *dynamodb.CreateTableInput) []*dynamodb.GlobalSecondaryIndexUpdate {

	expectedIndexes := make(map[string]*dynamodb.GlobalSecondaryIndex)
	for _, indexDefinition := range expected.GlobalSecondaryIndexes {
		expectedIndexes[*indexDefinition.IndexName] = indexDefinition
	}

	existingIndexes := make(map[string]interface{})
	var updates []*dynamodb.GlobalSecondaryIndexUpdate

	for _, existingIndex := range existing.Table.GlobalSecondaryIndexes {
		if _, mustBeKept := expectedIndexes[*existingIndex.IndexName]; !mustBeKept {
			updates = append(updates, &dynamodb.GlobalSecondaryIndexUpdate{
				Delete: &dynamodb.DeleteGlobalSecondaryIndexAction{
					IndexName: existingIndex.IndexName,
				},
			})
		}

		existingIndexes[*existingIndex.IndexName] = nil
	}

	for expectedIndexName, expectedIndex := range expectedIndexes {
		if _, mustNotBeCreated := existingIndexes[expectedIndexName]; !mustNotBeCreated {
			updates = append(updates, &dynamodb.GlobalSecondaryIndexUpdate{
				Create: &dynamodb.CreateGlobalSecondaryIndexAction{
					IndexName:             expectedIndex.IndexName,
					KeySchema:             expectedIndex.KeySchema,
					Projection:            expectedIndex.Projection,
					ProvisionedThroughput: expectedIndex.ProvisionedThroughput,
				},
			})
		}
	}

	return updates
}

func (r *rep) readTableVersion(table *dynamodb.DescribeTableOutput) (version string, err error) {
	resource, err := r.db.ListTagsOfResource(&dynamodb.ListTagsOfResourceInput{
		ResourceArn: table.Table.TableArn,
	})
	if err != nil {
		return "", errors.Wrapf(err, "couldn't read %s tag", versionTagName)
	}

	for _, t := range resource.Tags {
		if *t.Key == versionTagName {
			version = *t.Value
		}
	}

	return
}
