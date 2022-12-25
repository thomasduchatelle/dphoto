package migrator

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogdynamo"
)

type TransformationUpDateIndex struct{}

func (i TransformationUpDateIndex) PreScan(run *TransformationRun) error {
	repository := catalogdynamo.Must(catalogdynamo.NewRepository(run.Session, run.DynamoDBTableName))
	if repo, ok := repository.(*catalogdynamo.Repository); ok {
		log.Infoln("Updating indexes ...")
		return errors.Wrap(repo.CreateTableIfNecessary(), "failed to init DynamoDb structure (indexes, ...)")
	} else {
		log.Warn("catalogdynamo.NewRepository hasn't return the right type")
	}

	return nil
}
