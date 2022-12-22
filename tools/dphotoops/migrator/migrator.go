package migrator

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/dynamoutils"
)

func Migrate(tableName string, arn string, repopulateCache bool, scripts []interface{}) (int, error) {
	awsSession := session.Must(session.NewSession())
	run := &TransformationRun{
		Session:           awsSession,
		DynamoDBClient:    dynamodb.New(awsSession),
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

	var patches []*dynamodb.WriteRequest

	err := run.DynamoDBClient.ScanPages(&dynamodb.ScanInput{
		TableName: &tableName,
	}, func(output *dynamodb.ScanOutput, _ bool) bool {
		for _, item := range output.Items {
			run.Counter.Inc("RECORDS", 1)
			for _, tr := range scanTransformations {
				newPatches, err := tr.GeneratePatches(run, item)
				if err != nil {
					log.WithError(err).Errorf("Failed to apply transformation %+v: %s", tr, err.Error())
					return false
				}

				patches = append(patches, newPatches...)
			}
		}

		return true
	})
	if err != nil {
		return 0, err
	}

	log.Infof("Types count: %+v\n", run.Counter)

	if len(patches) > 0 {
		log.Infof("Running %d updates ... ", len(patches))
		err = dynamoutils.BufferedWriteItems(run.DynamoDBClient, patches, tableName, dynamoutils.DynamoWriteBatchSize)
		if err != nil {
			return 0, err
		}
	} else {
		log.Infof("Nothing to migrate, everything is already up to date.")
	}

	total, _ := run.Counter["RECORDS"]
	return total, nil

}
