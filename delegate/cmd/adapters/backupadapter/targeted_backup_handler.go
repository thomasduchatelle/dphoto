package backupadapter

import (
	"duchatelle.io/dphoto/dphoto/backup"
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"duchatelle.io/dphoto/dphoto/cmd/backupui"
	"duchatelle.io/dphoto/dphoto/cmd/ui"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/pkg/errors"
	"strings"
)

type TargetedBackupHandler struct {
	Owner string
}

func NewBackupHandler(owner string) ui.BackupSuggestionPort {
	return &TargetedBackupHandler{
		Owner: owner,
	}
}

type BackupAlbumFilter struct {
	FolderName string
}

func (b *BackupAlbumFilter) AcceptAnalysedMedia(media *backupmodel.AnalysedMedia, folderName string) bool {
	return folderName == b.FolderName
}

func (t *TargetedBackupHandler) BackupSuggestion(record *ui.SuggestionRecord, existing *ui.ExistingRecord, renderer ui.InteractiveRendererPort) error {
	if folder, ok := record.Original.(*backupmodel.ScannedFolder); ok {
		restriction := ""
		if existing != nil {
			restriction = fmt.Sprintf(", restrictited to album %s [%s -> %s]", existing.Name, existing.Start.Format("2006-01-02"), existing.End.Format("2006-01-02"))
		}
		renderer.Print(fmt.Sprintf("Backing up %s%s [Y/n] ?", aurora.Cyan(folder.BackupVolume.Path), restriction))
		val, err := renderer.ReadAnswer()
		if err != nil || strings.ToLower(val) != "y\n" && strings.ToLower(val) != "\n" {
			return err
		}

		renderer.TakeOverScreen()

		listener := backupui.NewProgress()
		options := backup.Options{Listener: listener}
		if existing != nil {
			options.PostAnalyseFilter = &BackupAlbumFilter{FolderName: existing.FolderName}
		}

		report, err := backup.StartBackupRunner(t.Owner, *folder.BackupVolume, options)
		listener.Stop()

		if err != nil {
			return err
		}

		backupui.PrintBackupStats(report, folder.BackupVolume.Path)
		renderer.Print("Hit enter to go back.")
		_, err = renderer.ReadAnswer()
		return err
	}

	return errors.Errorf("Original not supported: %+v", record.Original)
}
