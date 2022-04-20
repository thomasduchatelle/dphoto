package backupadapter

import (
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/dphoto/backup"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"
	"github.com/thomasduchatelle/dphoto/dphoto/cmd/screen"
	"github.com/thomasduchatelle/dphoto/dphoto/cmd/ui"
)

type ScanProgress struct {
	screen       *screen.AutoRefreshScreen
	scanningLine *screen.ProgressLine
	analysedLine *screen.ProgressLine
}

func ScanWithCache(owner string, volume backupmodel.VolumeToBackup, options backup.ScanOptions) (ui.SuggestionRecordRepositoryPort, []backupmodel.FoundMedia, error) {
	previousResult, rejectCount, err := restore(volume.UniqueId)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to restore previous scan result for volume %s", volume.String())
	}

	if len(previousResult) > 0 {
		useState, ok := ui.NewSimpleForm().ReadBool("Previous result has been found for this volume, do you want to restore it?", "Y/n")
		if !ok || useState {
			return NewSuggestionRepository(owner, previousResult, rejectCount), nil, err
		}
	}

	suggestions, rejects, err := doScan(volume, options)
	return NewSuggestionRepository(owner, suggestions, len(rejects)), rejects, err
}

func doScan(volume backupmodel.VolumeToBackup, options backup.ScanOptions) ([]*backupmodel.ScannedFolder, []backupmodel.FoundMedia, error) {
	progress := newScanProgress()
	options.Listeners = append(options.Listeners, progress)

	suggestions, rejects, err := backup.ScanVolume(volume, options)
	progress.screen.Stop()

	if err != nil {
		return nil, nil, err
	}

	err = Store(volume.UniqueId, suggestions, len(rejects))
	return suggestions, rejects, err
}

func newScanProgress() *ScanProgress {
	table := screen.NewTable(" ", 2, 20, 80, 25)
	scanningLine, scanningSegment := screen.NewProgressLine(table, "Scanning...")
	analysedLine, analysedSegment := screen.NewProgressLine(table, "Analysed...")

	progressScreen := screen.NewAutoRefreshScreen(
		screen.RenderingOptions{Width: 180},
		scanningSegment,
		analysedSegment,
	)

	return &ScanProgress{
		screen:       progressScreen,
		scanningLine: scanningLine,
		analysedLine: analysedLine,
	}
}

func (s *ScanProgress) OnScanComplete(total uint) {
	s.scanningLine.SwapSpinner(1)
	s.scanningLine.SetLabel(aurora.Sprintf("%d files has been found.", aurora.Cyan(total)))
}

func (s *ScanProgress) OnAnalyseProgress(count, total uint) {
	s.analysedLine.SetBar(count, total)
	s.analysedLine.SetExplanation(fmt.Sprintf("%d / %d", count, total))

	if count < total {
		s.analysedLine.SetLabel("Reading...")
	} else {
		s.analysedLine.SwapSpinner(1)
		s.analysedLine.SetLabel("Reading completed.")
	}
}
