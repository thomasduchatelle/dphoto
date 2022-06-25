package bootstrap

import (
	log "github.com/sirupsen/logrus"
	_ "github.com/thomasduchatelle/dphoto/domain/backupadapters/analysers"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/catalogarchive"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/catalogdynamo"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
)

func init() {
	config.Listen(func(cfg config.Config) {
		log.Debugln("connecting catalog adapters (dynamodb)")
		catalog.Init(catalogdynamo.Must(catalogdynamo.NewRepository(cfg.GetAWSSession(), cfg.GetString("catalog.dynamodb.table"))), catalogarchive.New())
	})
}
