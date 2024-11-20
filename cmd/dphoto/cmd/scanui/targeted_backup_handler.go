package scanui

import (
	"context"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/backupui"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/ui"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"strings"
)

type TargetedBackupHandler struct {
	Owner             string
	SubVolumeResolver func(absolutePath string) (backup.SourceVolume, error)
	ScanOptions       backup.Options // ScanOptions are the options used for the original scan
}

func NewBackupHandler(owner string, resolver func(absolutePath string) (backup.SourceVolume, error), options backup.Options) ui.BackupSuggestionPort {
	return &TargetedBackupHandler{
		Owner:             owner,
		SubVolumeResolver: resolver,
		ScanOptions:       options,
	}
}

func (t *TargetedBackupHandler) BackupSuggestion(record *ui.SuggestionRecord, existing *ui.ExistingRecord, renderer ui.InteractiveRendererPort) error {
	ctx := context.TODO()
	if folder, ok := record.Original.(*backup.ScannedFolder); ok {
		subVolume, err := t.SubVolumeResolver(folder.AbsolutePath)
		if err != nil {
			return err
		}

		restriction := ""
		if existing != nil {
			restriction = fmt.Sprintf(", restrictited to album %s [%s -> %s]", existing.Name, existing.Start.Format("2006-01-02"), existing.End.Format("2006-01-02"))
		}
		renderer.Print(fmt.Sprintf("Backing up %s%s [Y/n] ?", aurora.Cyan(subVolume), restriction))
		val, err := renderer.ReadAnswer()
		if err != nil || strings.ToLower(val) != "y\n" && strings.ToLower(val) != "\n" {
			return err
		}

		renderer.TakeOverScreen()

		listener := backupui.NewProgress()
		options := []backup.Options{
			backup.OptionsWithListener(listener),
			t.ScanOptions,
		}
		options = append(options, config.BackupOptions()...)

		if existing != nil {
			options = append(options, backup.OptionsOnlyAlbums(existing.FolderName))
		}

		multiFilesBackup := pkgfactory.NewMultiFilesBackup(ctx)
		report, err := multiFilesBackup(ctx, ownermodel.Owner(t.Owner), subVolume, options...)
		listener.Stop()

		if err != nil {
			return err
		}

		backupui.PrintBackupStats(report, subVolume.String())
		renderer.Print("Hit enter to go back.")
		_, err = renderer.ReadAnswer()
		return err
	}

	return errors.Errorf("Original not supported: %+v", record.Original)
}
