package backupproxy

import (
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/pkg/errors"
	backupui2 "github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/backupui"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/ui"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"strings"
)

type TargetedBackupHandler struct {
	Owner             string
	SubVolumeResolver func(absolutePath string) (backup.SourceVolume, error)
}

func NewBackupHandler(owner string, resolver func(absolutePath string) (backup.SourceVolume, error)) ui.BackupSuggestionPort {
	return &TargetedBackupHandler{
		Owner:             owner,
		SubVolumeResolver: resolver,
	}
}

type BackupAlbumFilter struct {
	FolderName string
}

func (b *BackupAlbumFilter) AcceptAnalysedMedia(media *backup.AnalysedMedia, folderName string) bool {
	return folderName == b.FolderName
}

func (t *TargetedBackupHandler) BackupSuggestion(record *ui.SuggestionRecord, existing *ui.ExistingRecord, renderer ui.InteractiveRendererPort) error {
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

		listener := backupui2.NewProgress()
		options := []backup.Options{
			backup.OptionWithListener(listener),
		}
		if existing != nil {
			options = append(options, backup.OptionOnlyAlbums(existing.FolderName))
		}

		report, err := backup.Backup(t.Owner, subVolume, options...)
		listener.Stop()

		if err != nil {
			return err
		}

		backupui2.PrintBackupStats(report, subVolume.String())
		renderer.Print("Hit enter to go back.")
		_, err = renderer.ReadAnswer()
		return err
	}

	return errors.Errorf("Original not supported: %+v", record.Original)
}
