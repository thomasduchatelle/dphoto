package cmd

import (
	"github.com/spf13/cobra"
	backupui2 "github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/backupui"
	config2 "github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/filesystemvolume"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/s3volume"
	"strings"
)

var newS3Volume func(volumePath string) (backup.SourceVolume, error)

var backupCmd = &cobra.Command{
	Use:   "backup <source path>",
	Short: "Backup photos and videos to personal cloud",
	Long:  `Backup photos and videos to personal cloud`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		volumePath := args[0]

		progress := backupui2.NewProgress()
		volume, err := newSmartVolume(volumePath)
		printer.FatalIfError(err, 1)

		report, err := backup.Backup(Owner, volume, backup.OptionWithListener(progress))
		printer.FatalIfError(err, 2)

		progress.Stop()

		backupui2.PrintBackupStats(report, volumePath)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	config2.Listen(func(cfg config2.Config) {
		newS3Volume = func(volumePath string) (backup.SourceVolume, error) {
			return s3volume.New(cfg.GetAWSSession(), volumePath)
		}
	})
}

func newSmartVolume(volumePath string) (backup.SourceVolume, error) {
	if strings.HasPrefix(volumePath, "s3://") {
		return newS3Volume(volumePath)
	}

	return filesystemvolume.New(volumePath), nil
}
