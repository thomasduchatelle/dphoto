package migrator

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/catalogdynamo"
)

type TransformationUpDateIndex struct{}

func (i TransformationUpDateIndex) PreScan(run *TransformationRun) error {
	log.Infoln("Updating indexes ...")
	_, err := catalogdynamo.NewRepository(run.Session, run.DynamoDBTableName)
	return errors.Wrap(err, "failed to init DynamoDb structure (indexes, ...)")
}
