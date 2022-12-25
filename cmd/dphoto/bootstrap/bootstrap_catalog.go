package bootstrap

import (
	log "github.com/sirupsen/logrus"
	config2 "github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	_ "github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysers"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogarchivesync"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogdynamo"
)

func init() {
	config2.Listen(func(cfg config2.Config) {
		log.Debugln("connecting catalog adapters (dynamodb)")
		repository := catalogdynamo.Must(catalogdynamo.NewRepository(cfg.GetAWSSession(), cfg.GetString(config2.CatalogDynamodbTable)))
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
