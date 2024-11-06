package config

import (
	"github.com/thomasduchatelle/dphoto/pkg/backup"
)

func BackupOptions() []backup.Options {
	return []backup.Options{
		backup.OptionsConcurrentAnalyserRoutines(config.GetIntOrDefault(BackupConcurrencyAnalyser, 4)),
		backup.OptionsConcurrentCataloguerRoutines(config.GetIntOrDefault(BackupConcurrencyCataloguer, 2)),
		backup.OptionsConcurrentUploaderRoutines(config.GetIntOrDefault(BackupConcurrencyUploader, 2)),
	}
}
