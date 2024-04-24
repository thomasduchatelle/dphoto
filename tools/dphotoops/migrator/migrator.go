package migrator

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
)

func Migrate(tableName string, arn string, repopulateCache bool, scripts []interface{}) (int, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return 0, err
	}

	run := &TransformationRun{
		DynamoDBClient:    dynamodb.NewFromConfig(cfg),
		DynamoDBTableName: tableName,
		TopicARN:          arn,
		Counter:           make(map[string]int),
	}

	var scanTransformations []ScanTransformation
	for _, transformation := range scripts {
		if preScan, ok := transformation.(PreScanTransformation); ok {
			err := preScan.PreScan(run)
			if err != nil {
				return 0, err
			}
		}
		if duringScan, ok := transformation.(ScanTransformation); ok {
			scanTransformations = append(scanTransformations, duringScan)
		}
	}

	if len(scanTransformations) == 0 {
		return 0, nil
	}

	log.Infof("Scanning through the all records")

	var patches []types.WriteRequest

	paginator := dynamodb.NewScanPaginator(run.DynamoDBClient, &dynamodb.ScanInput{
		TableName: &tableName,
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return 0, err
		}

		for _, item := range page.Items {
			run.Counter.Inc("RECORDS", 1)
			for _, tr := range scanTransformations {
				newPatches, err := tr.GeneratePatches(run, item)
				if err != nil {
					log.WithError(err).Errorf("Failed to apply transformation %+v: %s", tr, err.Error())
					return 0, err
				}

				patches = append(patches, newPatches...)
			}
		}

	}

	log.Infof("Types count: %+v\n", run.Counter)

	if len(patches) > 0 {
		log.Infof("Running %d updates ... ", len(patches))
		err = dynamoutils.BufferedWriteItems(context.TODO(), run.DynamoDBClient, patches, tableName, dynamoutils.DynamoWriteBatchSize)
		if err != nil {
			return 0, err
		}
	} else {
		log.Infof("Nothing to migrate, everything is already up to date.")
	}

	total, _ := run.Counter["RECORDS"]
	return total, nil

}
