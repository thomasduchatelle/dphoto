package bootstrap

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	_ "github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysers"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

func init() {
	config.Listen(func(cfg config.Config) {
		ctx := context.TODO()

		log.Debugln("connecting catalog adapters (dynamodb)")
		table := cfg.GetString(config.CatalogDynamodbTable)

		if cfg.GetBool(config.Localstack) {
			err := appdynamodb.CreateTableIfNecessary(ctx, table, dynamodb.NewFromConfig(cfg.GetAWSV2Config()), true)
			if err != nil {
				panic("Failed while updating indexes: " + err.Error())
			}
		}

		catalog.Init(pkgfactory.CatalogRepository(ctx))
	})
}
