package cmd

import (
	"duchatelle.io/dphoto/dphoto/backup"
	"duchatelle.io/dphoto/dphoto/backup/model"
	"duchatelle.io/dphoto/dphoto/cmd/printer"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"path/filepath"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup photos and videos to personal cloud",
	Long:  `Backup photos and videos to personal cloud`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		volumePath, err := filepath.Abs(args[0])
		printer.FatalWithMessageIfError(err, 2, "provided argument must be a valid file path")

		err = backup.StartBackupRunner(model.VolumeToBackup{
			UniqueId: volumePath,
			Type:     model.VolumeTypeFileSystem,
			Path:     volumePath,
			Local:    true,
		})
		printer.FatalIfError(err, 1)

		printer.Success("Backup of %s complete", aurora.Cyan(volumePath))
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
