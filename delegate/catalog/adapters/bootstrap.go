package adapters

import (
	"duchatelle.io/dphoto/dphoto/catalog"
	"duchatelle.io/dphoto/dphoto/catalog/adapters/dynamo"
	"duchatelle.io/dphoto/dphoto/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	config.Listen(func(cfg config.Config) {
		log.Debugln("connecting catalog adapters (dynamodb)")
		catalog.Repository = dynamo.Must(dynamo.NewRepository(cfg.GetAWSSession(), cfg.GetString("owner"), cfg.GetString("catalog.dynamodb.table")))
	})
}
