package dynamoutils

import (
	"context"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

type CreateOrUpdateTableInput struct {
	Client     *dynamodb.DynamoDB
	TableName  string
	Definition *dynamodb.CreateTableInput
}

// CreateOrUpdateTable creates the table if it doesn't exist ; or update it.
// Implementation is not mature and is subject to a lot of limitations (order in which fields and index are deleted)
func CreateOrUpdateTable(ctx context.Context, input *CreateOrUpdateTableInput) error {
	if input.TableName == "" {
		return errors.Errorf("tableName must not be empty")
	}

	mdc := log.WithFields(log.Fields{
		"TableBackup": input.TableName,
	})
	mdc.Debugf("CreateTableIfNecessary > describe table '%s'", input.TableName)

	dynamoClient := input.Client
	table, err := dynamoClient.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: &input.TableName,
	})

	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == dynamodb.ErrCodeResourceNotFoundException {
		table = nil
	} else if err != nil {
		return errors.Wrap(err, "failed to find existing table")
	}

	if table == nil {
		mdc.Infoln("Creating dynamodb table...")
		_, err = dynamoClient.CreateTable(input.Definition)

		return errors.Wrapf(err, "failed to create table %s", input.TableName)

	} else {
		updates := generatedSecondaryUpdatesIndexes(table, input.Definition)

		for i, update := range updates {
			if update.Delete != nil {
				mdc.Infof("[%d/%d] Deleting table index %s", i+1, len(updates), *update.Delete.IndexName)
			} else if update.Create != nil {
				mdc.Infof("[%d/%d] Creating table index %s", i+1, len(updates), *update.Create.IndexName)
			} else if update.Update != nil {
				mdc.Infof("[%d/%d] Updating table index %s", i+1, len(updates), *update.Create.IndexName)
			}

			err := waitDbToBeReady(ctx, dynamoClient, input.TableName)
			if err != nil {
				return err
			}

			_, err = dynamoClient.UpdateTable(&dynamodb.UpdateTableInput{
				AttributeDefinitions:        input.Definition.AttributeDefinitions,
				GlobalSecondaryIndexUpdates: []*dynamodb.GlobalSecondaryIndexUpdate{update},
				TableName:                   &input.TableName,
			})
			if err != nil {
				return errors.Wrapf(err, "failed to update table %s", input.TableName)
			}
		}

		if len(updates) > 0 && len(input.Definition.Tags) > 0 {
			_, err := dynamoClient.TagResource(&dynamodb.TagResourceInput{
				ResourceArn: table.Table.TableArn,
				Tags:        input.Definition.Tags,
			})
			if err != nil {
				return errors.Wrapf(err, "failed to tag the DynamoDB table")
			}
		}
	}

	return nil
}

func waitDbToBeReady(ctx context.Context, client *dynamodb.DynamoDB, tableName string) error {
	for {
		tick := time.Tick(10 * time.Second)
		select {
		case <-ctx.Done():
			return errors.Errorf("wait has been cancelled")

		case <-tick:
			table, err := client.DescribeTable(&dynamodb.DescribeTableInput{
				TableName: &tableName,
			})
			if err != nil {
				return errors.Wrapf(err, "failed describing table waiting all index have been processed")
			}

			indexesActive := true
			for _, index := range table.Table.GlobalSecondaryIndexes {
				indexesActive = indexesActive && *index.IndexStatus == dynamodb.IndexStatusActive
			}
			if *table.Table.TableStatus == dynamodb.TableStatusActive && indexesActive {
				// ready !
				return nil
			}
		}
	}
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
