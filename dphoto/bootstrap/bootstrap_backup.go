package bootstrap

import (
	"github.com/thomasduchatelle/dphoto/domain/backup"
	"github.com/thomasduchatelle/dphoto/domain/backupadapters/backuparchive"
	"github.com/thomasduchatelle/dphoto/domain/backupadapters/backupcatalog"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/catalogdynamo"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
)

func init() {
	config.Listen(func(cfg config.Config) {
		backup.ConcurrentAnalyser = cfg.GetIntOrDefault(config.BackupConcurrencyAnalyser, 4)
		backup.ConcurrentCataloguer = cfg.GetIntOrDefault(config.BackupConcurrencyCataloguer, 2)
		backup.ConcurrentUploader = cfg.GetIntOrDefault(config.BackupConcurrencyUploader, 2)
		backup.BatchSize = catalogdynamo.DynamoReadBatchSize // optimise the cataloguer and scanning

		backup.Init(backupcatalog.New(), backuparchive.New())
	})
}
