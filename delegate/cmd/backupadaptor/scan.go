package backupadaptor

import (
	"duchatelle.io/dphoto/dphoto/backup"
	"duchatelle.io/dphoto/dphoto/backup/model"
	"duchatelle.io/dphoto/dphoto/cmd/screen"
	"duchatelle.io/dphoto/dphoto/cmd/ui"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/pkg/errors"
)

type ScanProgress struct {
	screen       *screen.AutoRefreshScreen
	scanningLine *screen.ProgressLine
	analysedLine *screen.ProgressLine
}

func ScanWithCache(volume string) (ui.RecordRepositoryPort, int, error) {
	previousResult, err := restore(volume)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "failed to restore previous scan result for volume %s", volume)
	}

	if len(previousResult) > 0 {
		useState, ok := ui.NewSimpleForm().ReadBool("Previous result has been found for this volume, do you want to restore it?", "Y/n")
		if !ok || useState {
			return NewSuggestionRepository(previousResult), len(previousResult), err
		}
	}

	suggestions, err := doScan(volume)
	return NewSuggestionRepository(suggestions), len(suggestions), err
}

func doScan(volume string) ([]*backup.FoundAlbum, error) {
	progress := newScanProgress()
	suggestions, err := backup.DiscoverAlbumFromSource(model.VolumeToBackup{
		UniqueId: volume,
		Type:     model.VolumeTypeFileSystem,
		Path:     volume,
		Local:    true,
	}, progress)
	progress.screen.Stop()

	if err != nil {
		return nil, err
	}

	err = Store(volume, suggestions)
	return suggestions, err
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
