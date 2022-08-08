package bootstrap

import (
	log "github.com/sirupsen/logrus"
	_ "github.com/thomasduchatelle/dphoto/domain/backupadapters/analysers"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/carchivesync"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/repositorydynamo"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
)

func init() {
	config.Listen(func(cfg config.Config) {
		log.Debugln("connecting catalog adapters (dynamodb)")
		catalog.Init(repositorydynamo.Must(repositorydynamo.NewRepository(cfg.GetAWSSession(), cfg.GetString(config.CatalogDynamodbTable))), carchivesync.New())
	})
}
