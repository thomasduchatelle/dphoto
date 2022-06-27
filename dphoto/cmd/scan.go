package cmd

import (
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/domain/backup"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/dphoto/cmd/adapters/backupproxy"
	"github.com/thomasduchatelle/dphoto/dphoto/cmd/ui"
	"github.com/thomasduchatelle/dphoto/dphoto/printer"
	"os"
	"path"
	"time"
)

var (
	scanArgs = struct {
		nonInteractive bool
		skipRejects    bool
		rejectFile     string
	}{}
)

var scan = &cobra.Command{
	Use:   "scan <folder to scan>",
	Short: "Discover directory structure to suggest new albums to create",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		volume := args[0]

		smartVolume, err := newSmartVolume(volume)
		recordRepository, rejects, err := backupproxy.ScanWithCache(Owner, smartVolume, backup.OptionSkipRejects(scanArgs.skipRejects))
		printer.FatalIfError(err, 2)

		if len(rejects) > 0 && scanArgs.rejectFile != "" {
			err = writeRejectsInFile(rejects, scanArgs.rejectFile)
			printer.FatalIfError(err, 3)
		}

		if recordRepository.Count() == 0 {
			fmt.Println(aurora.Yellow(fmt.Sprintf("No new media found on volume %s.", aurora.Cyan(volume))))
		} else if scanArgs.nonInteractive {
			err = ui.NewSimpleSession(backupproxy.NewAlbumRepository(Owner), recordRepository).Render()
			printer.FatalIfError(err, 1)
		} else {
			err = ui.NewInteractiveSession(&uiCatalogAdapter{backupproxy.NewBackupHandler(Owner, newSmartVolume)}, backupproxy.NewAlbumRepository(Owner), recordRepository, Owner).Start()
			printer.FatalIfError(err, 1)
		}

		os.Exit(0)
	},
}

func writeRejectsInFile(rejects []backup.FoundMedia, file string) error {
	err := os.MkdirAll(path.Dir(file), 0744)
	if err != nil {
		return err
	}

	rejectFile, err := os.Create(file)
	if err != nil {
		return err
	}

	defer rejectFile.Close()

	for _, media := range rejects {
		_, err = rejectFile.WriteString(fmt.Sprintf("%s\n", media))
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(scan)

	scan.Flags().BoolVarP(&scanArgs.nonInteractive, "non-interactive", "I", false, "Disable interactive output and only display the scan results.")
	scan.Flags().BoolVarP(&scanArgs.skipRejects, "skip-errors", "s", false, "Unreadable files, or files without date, will be reported as 'rejects' and printed in rejected file.")
	scan.Flags().StringVar(&scanArgs.rejectFile, "rejects", "", "Unreadable files, or files without date, will be listed in the given file. Requires to use --skip-errors.")
}

type uiCatalogAdapter struct {
	ui.BackupSuggestionPort
}

func (o uiCatalogAdapter) Create(request ui.RecordCreation) error {
	return catalog.Create(catalog.CreateAlbum{
		Owner:            request.Owner,
		Name:             request.Name,
		Start:            request.Start,
		End:              request.End,
		ForcedFolderName: request.FolderName,
	})
}

func (o *uiCatalogAdapter) RenameAlbum(folderName, newName string, renameFolder bool) error {
	return catalog.RenameAlbum(Owner, folderName, newName, renameFolder)
}

func (o *uiCatalogAdapter) UpdateAlbum(folderName string, start, end time.Time) error {
	return catalog.UpdateAlbum(Owner, folderName, start, end)
}

func (o *uiCatalogAdapter) DeleteAlbum(folderName string) error {
	return catalog.DeleteAlbum(Owner, folderName, false)
}
