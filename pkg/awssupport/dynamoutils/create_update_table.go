package dynamoutils

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

type CreateOrUpdateTableInput struct {
	Client     *dynamodb.Client
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

	client := input.Client
	table, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: &input.TableName,
	})

	var resourceNotFoundException *types.ResourceNotFoundException
	if errors.As(err, &resourceNotFoundException) {
		mdc.Infoln("Creating dynamodb table...")
		_, err = client.CreateTable(ctx, input.Definition)

		return errors.Wrapf(err, "failed to create table %s", input.TableName)

	} else if err != nil {
		return errors.Wrap(err, "failed to read existing table structure")
	}

	updates := generatedSecondaryUpdatesIndexes(table, input.Definition)

	if len(updates) == 0 {
		mdc.Infof("No change required on dynamodb table - update complete.")
	}

	for i, update := range updates {
		if update.Delete != nil {
			mdc.Infof("[%d/%d] Deleting table index %s", i+1, len(updates), *update.Delete.IndexName)
		} else if update.Create != nil {
			mdc.Infof("[%d/%d] Creating table index %s", i+1, len(updates), *update.Create.IndexName)
		} else if update.Update != nil {
			mdc.Infof("[%d/%d] Updating table index %s", i+1, len(updates), *update.Create.IndexName)
		}

		err := waitDbToBeReady(ctx, client, input.TableName)
		if err != nil {
			return err
		}

		_, err = client.UpdateTable(ctx, &dynamodb.UpdateTableInput{
			AttributeDefinitions:        input.Definition.AttributeDefinitions,
			GlobalSecondaryIndexUpdates: []types.GlobalSecondaryIndexUpdate{update},
			TableName:                   &input.TableName,
		})
		if err != nil {
			return errors.Wrapf(err, "failed to update table %s", input.TableName)
		}
	}

	if len(updates) > 0 && len(input.Definition.Tags) > 0 {
		_, err := client.TagResource(ctx, &dynamodb.TagResourceInput{
			ResourceArn: table.Table.TableArn,
			Tags:        input.Definition.Tags,
		})
		if err != nil {
			return errors.Wrapf(err, "failed to tag the DynamoDB table")
		}
	}

	return nil
}

func waitDbToBeReady(ctx context.Context, client *dynamodb.Client, tableName string) error {
	for {
		tick := time.Tick(10 * time.Second)
		select {
		case <-ctx.Done():
			return errors.Errorf("wait has been cancelled")

		case <-tick:
			table, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
				TableName: &tableName,
			})
			if err != nil {
				return errors.Wrapf(err, "failed describing table waiting all index have been processed")
			}

			indexesActive := true
			for _, index := range table.Table.GlobalSecondaryIndexes {
				indexesActive = indexesActive && index.IndexStatus == types.IndexStatusActive
			}
			if table.Table.TableStatus == types.TableStatusActive && indexesActive {
				// ready !
				return nil
			}
		}
	}
}

func generatedSecondaryUpdatesIndexes(existing *dynamodb.DescribeTableOutput, expected *dynamodb.CreateTableInput) []types.GlobalSecondaryIndexUpdate {

	expectedIndexes := make(map[string]types.GlobalSecondaryIndex)
	for _, indexDefinition := range expected.GlobalSecondaryIndexes {
		expectedIndexes[*indexDefinition.IndexName] = indexDefinition
	}

	existingIndexes := make(map[string]interface{})
	var updates []types.GlobalSecondaryIndexUpdate

	for _, existingIndex := range existing.Table.GlobalSecondaryIndexes {
		if _, mustBeKept := expectedIndexes[*existingIndex.IndexName]; !mustBeKept {
			updates = append(updates, types.GlobalSecondaryIndexUpdate{
				Delete: &types.DeleteGlobalSecondaryIndexAction{
					IndexName: existingIndex.IndexName,
				},
			})
		}

		existingIndexes[*existingIndex.IndexName] = nil
	}

	for expectedIndexName, expectedIndex := range expectedIndexes {
		if _, mustNotBeCreated := existingIndexes[expectedIndexName]; !mustNotBeCreated {
			updates = append(updates, types.GlobalSecondaryIndexUpdate{
				Create: &types.CreateGlobalSecondaryIndexAction{
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
