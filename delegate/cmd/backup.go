package cmd

import (
	"github.com/thomasduchatelle/dphoto/delegate/backup"
	"github.com/thomasduchatelle/dphoto/delegate/backup/backupmodel"
	"github.com/thomasduchatelle/dphoto/delegate/cmd/backupui"
	"github.com/thomasduchatelle/dphoto/delegate/cmd/printer"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
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
		volumePath := args[0]

		progress := backupui.NewProgress()

		volume := newSmartVolume(volumePath)
		if backupArgs.remote {
			volume.Local = false
		}
		report, err := backup.StartBackupRunner(Owner, volume, backup.Options{Listener: progress})
		printer.FatalIfError(err, 1)

		progress.Stop()

		backupui.PrintBackupStats(report, volumePath)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().BoolVarP(&backupArgs.remote, "remote", "r", false, "mark the source as remote ; a local buffer will be used to read files only once")
}

func newSmartVolume(volumePath string) backupmodel.VolumeToBackup {
	volumeType := VolumeTypeFromPath(volumePath)
	volumeId := volumePath
	if volumeType == backupmodel.VolumeTypeFileSystem {
		volumeId, _ = filepath.Abs(volumePath)
	}
	return backupmodel.VolumeToBackup{
		UniqueId: volumeId,
		Type:     volumeType,
		Path:     volumePath,
		Local:    true,
	}
}

func VolumeTypeFromPath(arg string) backupmodel.VolumeType {
	if strings.HasPrefix(arg, "s3://") {
		return backupmodel.VolumeTypeS3
	}

	return backupmodel.VolumeTypeFileSystem
}
