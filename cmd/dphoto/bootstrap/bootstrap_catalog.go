package bootstrap

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	_ "github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysers"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogarchivesync"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogdynamo"
)

func init() {
	config.Listen(func(cfg config.Config) {
		log.Debugln("connecting catalog adapters (dynamodb)")
		table := cfg.GetString(config.CatalogDynamodbTable)

		log.Infoln("Updating indexes ...")
		err := appdynamodb.CreateTableIfNecessary(context.TODO(), table, dynamodb.NewFromConfig(cfg.GetAWSV2Config()), false)
		if err != nil {
			panic("Failed while updating indexes: " + err.Error())
		}

		repository := catalogdynamo.Must(catalogdynamo.NewRepository(cfg.GetAWSV2Config(), table))
		catalog.Init(repository, catalogarchivesync.New())
	})
}
