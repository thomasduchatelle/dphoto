package migrator

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
)

type TransformationUpDateIndex struct{}

func (i TransformationUpDateIndex) PreScan(run *TransformationRun) error {
	log.Infoln("Updating indexes ...")
	err := appdynamodb.CreateTableIfNecessary(context.TODO(), run.DynamoDBTableName, run.DynamoDBClient, false)
	return errors.Wrapf(err, "failed to init DynamoDb structure (indexes, ...)")
}
