package bootstrap

import (
	"github.com/thomasduchatelle/dphoto/domain/backup"
	"github.com/thomasduchatelle/dphoto/domain/backupadapters/backuparchive"
	"github.com/thomasduchatelle/dphoto/domain/backupadapters/backupcatalog"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
)

func init() {
	config.Listen(func(cfg config.Config) {
		backup.ConcurrentAnalyser = cfg.GetIntOrDefault("backup.concurrency.analyser", 4)
		backup.ConcurrentCataloguer = cfg.GetIntOrDefault("backup.concurrency.cataloguer", 2)
		backup.ConcurrentUploader = cfg.GetIntOrDefault("backup.concurrency.uploader", 2)
		backup.BatchSize = backupcatalog.RecommendedBatchSize

		backup.Init(backupcatalog.New(), backuparchive.New())
	})
}
