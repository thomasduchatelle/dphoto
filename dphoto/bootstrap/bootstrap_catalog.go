package bootstrap

import (
	log "github.com/sirupsen/logrus"
	_ "github.com/thomasduchatelle/dphoto/domain/backupadapters/analysers"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/catalogarchivesync"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/catalogdynamo"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
)

func init() {
	config.Listen(func(cfg config.Config) {
		log.Debugln("connecting catalog adapters (dynamodb)")
		repository := catalogdynamo.Must(catalogdynamo.NewRepository(cfg.GetAWSSession(), cfg.GetString(config.CatalogDynamodbTable)))
		if repo, ok := repository.(*catalogdynamo.Repository); ok {
			log.Infoln("Updating indexes ...")
			err := repo.CreateTableIfNecessary()
			if err != nil {
				panic("Failed while updating indexes: " + err.Error())
			}

		} else {
			log.Warn("catalogdynamo.NewRepository hasn't return the right type to update indexes")
		}

		catalog.Init(repository, catalogarchivesync.New())
	})
}
