package adapters

import (
	"duchatelle.io/dphoto/dphoto/catalog"
	"duchatelle.io/dphoto/dphoto/catalog/adapters/dynamo"
	config2 "duchatelle.io/dphoto/dphoto/internal/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	config2.Listen(func(cfg config2.Config) {
		log.Debugln("connecting catalog adapters (dynamodb)")
		catalog.Repository = dynamo.Must(dynamo.NewRepository(cfg.GetAWSSession(), cfg.GetString("catalog.dynamodb.table")))
	})
}
