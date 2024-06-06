package migrator

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
)

type TransformationUpDateIndex struct{}

func (i TransformationUpDateIndex) PreScan(run *TransformationRun) error {
	err := appdynamodb.CreateTableIfNecessary(context.TODO(), run.DynamoDBTableName, run.DynamoDBClient, false)
	return errors.Wrapf(err, "failed to init DynamoDb structure (indexes, ...)")
}
