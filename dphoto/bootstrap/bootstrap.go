package bootstrap

import (
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/catalogdynamo"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
)

func init() {
	config.Listen(func(cfg config.Config) {
		log.Debugln("connecting catalog adapters (dynamodb)")
		catalog.db = catalogdynamo.Must(catalogdynamo.NewRepository(cfg.GetAWSSession(), cfg.GetString("catalog.dynamodb.table")))
	})
}
