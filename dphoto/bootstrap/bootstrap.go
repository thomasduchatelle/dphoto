package bootstrap

import (
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/dynamo"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
)

func init() {
	config.Listen(func(cfg config.Config) {
		log.Debugln("connecting catalog adapters (dynamodb)")
		catalog.Repository = dynamo.Must(dynamo.NewRepository(cfg.GetAWSSession(), cfg.GetString("owner"), cfg.GetString("catalog.dynamodb.table")))
	})
}
