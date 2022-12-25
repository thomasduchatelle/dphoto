package bootstrap

import (
	config2 "github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/backuparchive"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/backupcatalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogdynamo"
)

func init() {
	config2.Listen(func(cfg config2.Config) {
		backup.ConcurrentAnalyser = cfg.GetIntOrDefault(config2.BackupConcurrencyAnalyser, 4)
		backup.ConcurrentCataloguer = cfg.GetIntOrDefault(config2.BackupConcurrencyCataloguer, 2)
		backup.ConcurrentUploader = cfg.GetIntOrDefault(config2.BackupConcurrencyUploader, 2)
		backup.BatchSize = catalogdynamo.DynamoReadBatchSize // optimise the cataloguer and scanning

		backup.Init(backupcatalog.New(), backuparchive.New())
	})
}
