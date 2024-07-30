package config

import (
	"github.com/thomasduchatelle/dphoto/pkg/backup"
)

func BackupOptions() []backup.Options {
	return []backup.Options{
		backup.WithConcurrentAnalyser(config.GetIntOrDefault(BackupConcurrencyAnalyser, 4)),
		backup.WithConcurrentCataloguer(config.GetIntOrDefault(BackupConcurrencyCataloguer, 2)),
		backup.WithConcurrentUploader(config.GetIntOrDefault(BackupConcurrencyUploader, 2)),
	}
}
