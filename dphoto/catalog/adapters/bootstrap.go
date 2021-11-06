package adapters

import (
	"github.com/thomasduchatelle/dphoto/dphoto/catalog"
	"github.com/thomasduchatelle/dphoto/dphoto/catalog/adapters/dynamo"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	config.Listen(func(cfg config.Config) {
		log.Debugln("connecting catalog adapters (dynamodb)")
		catalog.Repository = dynamo.Must(dynamo.NewRepository(cfg.GetAWSSession(), cfg.GetString("owner"), cfg.GetString("catalog.dynamodb.table")))
	})
}
