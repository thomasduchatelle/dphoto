package bootstrap

import (
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/backuparchive"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/backupcatalog"
)

func init() {
	config.Listen(func(cfg config.Config) {
		backup.ConcurrentAnalyser = cfg.GetIntOrDefault(config.BackupConcurrencyAnalyser, 4)
		backup.ConcurrentCataloguer = cfg.GetIntOrDefault(config.BackupConcurrencyCataloguer, 2)
		backup.ConcurrentUploader = cfg.GetIntOrDefault(config.BackupConcurrencyUploader, 2)
		backup.BatchSize = dynamoutils.DynamoReadBatchSize // optimise the cataloguer and scanning

		backup.Init(backupcatalog.New(), backuparchive.New())
	})
}
