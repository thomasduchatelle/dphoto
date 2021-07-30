package cmd

import (
	"duchatelle.io/dphoto/dphoto/backup"
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"duchatelle.io/dphoto/dphoto/cmd/backupui"
	"duchatelle.io/dphoto/dphoto/cmd/printer"
	"github.com/spf13/cobra"
	"path/filepath"
)

var (
	backupArgs = struct {
		remote bool
	}{}
)

var backupCmd = &cobra.Command{
	Use:   "backup [--remote] <source path>",
	Short: "Backup photos and videos to personal cloud",
	Long:  `Backup photos and videos to personal cloud`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		volumePath, err := filepath.Abs(args[0])
		printer.FatalWithMessageIfError(err, 2, "provided argument must be a valid file path")

		progress := backupui.NewProgress()

		report, err := backup.StartBackupRunner(Owner, backupmodel.VolumeToBackup{
			UniqueId: volumePath,
			Type:     backupmodel.VolumeTypeFileSystem,
			Path:     volumePath,
			Local:    !backupArgs.remote,
		}, backup.Options{Listener: progress})
		printer.FatalIfError(err, 1)

		progress.Stop()

		backupui.PrintBackupStats(report, volumePath)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().BoolVarP(&backupArgs.remote, "remote", "r", false, "mark the source as remote ; a local buffer will be used to read files only once")
}
