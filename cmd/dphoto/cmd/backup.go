package cmd

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/backupui"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysiscache"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/filesystemvolume"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/s3volume"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"io"
	"os"
	"path"
	"strings"
)

var (
	newS3Volume    func(volumePath string) (backup.SourceVolume, error)
	cacheDirectory string
)

var (
	backupCmdArg = struct {
		noCache bool
		confirm bool
	}{}
)

var backupCmd = &cobra.Command{
	Use:   "backup [--no-cache] [--ask] <source path>",
	Short: "Backup photos and videos to personal cloud",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		volumePath := args[0]

		progress := backupui.NewProgress()
		volume, err := newSmartVolume(volumePath)
		printer.FatalIfError(err, 1)

		multiFilesBackup := pkgfactory.NewMultiFilesBackup(ctx)
		options := []backup.Options{
			backup.OptionWithListener(progress).WithCachedAnalysis(addCacheAnalysis(!backupCmdArg.noCache)),
		}
		options = append(options, config.BackupOptions()...)
		report, err := multiFilesBackup(ctx, ownermodel.Owner(Owner), volume, options...)
		printer.FatalIfError(err, 2)

		progress.Stop()

		backupui.PrintBackupStats(report, volumePath)
	},
}

// addCacheAnalysis is shared between 'backup' and 'scan'
func addCacheAnalysis(cache bool) backup.AnalyserDecorator {
	if cache {
		decorator, err := analysiscache.NewCacheDecorator(cacheDirectory)
		if err != nil {
			panic(fmt.Sprintf("PANIC - cache couldn't be initiated: %s", err.Error()))
		}

		if closable, isClosable := decorator.(io.Closer); isClosable {
			postRunFunctions = append(postRunFunctions, func() error {
				return closable.Close()
			})
		}

		return decorator
	}

	return nil
}

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().BoolVarP(&backupCmdArg.noCache, "no-cache", "c", false, "set to true to ignore cache (and not building it)")

	config.Listen(func(cfg config.Config) {
		newS3Volume = func(volumePath string) (backup.SourceVolume, error) {
			return s3volume.New(s3.NewFromConfig(cfg.GetAWSV2Config()), volumePath)
		}

		defaultCacheDir := path.Join(cfg.GetStringOrDefault(config.LocalHome, os.ExpandEnv("$HOME/.dphoto")), "cache")
		cacheDirectory = cfg.GetStringOrDefault(config.BackupCacheDirectory, defaultCacheDir)
	})
}

func newSmartVolume(volumePath string) (backup.SourceVolume, error) {
	if strings.HasPrefix(volumePath, "s3://") {
		return newS3Volume(volumePath)
	}

	return filesystemvolume.New(volumePath), nil
}
