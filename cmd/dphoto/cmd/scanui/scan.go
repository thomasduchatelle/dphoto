package scanui

import (
	"context"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/ui"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/internal/screen"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

type ScanProgress struct {
	screen           *screen.AutoRefreshScreen
	scanningLine     *screen.ProgressLine
	analysedLine     *screen.ProgressLine
	onAnalysedCalled bool
}

func ScanWithProgress(owner string, volume backup.SourceVolume, options ...backup.Options) (ui.SuggestionRecordRepositoryPort, error) {
	suggestions, err := doScan(owner, volume, options...)
	return NewSuggestionRepository(owner, suggestions), err
}

func doScan(owner string, volume backup.SourceVolume, options ...backup.Options) ([]*backup.ScannedFolder, error) {
	ctx := context.TODO()

	progress := newScanProgress()
	options = append(options, backup.OptionsWithListener(progress))
	options = append(options, config.BackupOptions()...)

	scanner := pkgfactory.NewMultiFilesScanner(ctx)
	suggestions, err := scanner(ctx, owner, volume, options...)
	progress.stop()

	if err != nil {
		return nil, err
	}

	return suggestions, err
}

func newScanProgress() *ScanProgress {
	table := screen.NewTable(" ", 2, 20, 80, 25)
	scanningLine, scanningSegment := screen.NewProgressLine(table, "Scanning...")
	analysedLine, analysedSegment := screen.NewProgressLine(table, "Analysing...")

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

func (s *ScanProgress) OnScanComplete(total backup.MediaCounter) {
	s.scanningLine.SwapSpinner(1)
	s.scanningLine.SetLabel(aurora.Sprintf("%d files has been found.", aurora.Cyan(total.Count)))
}

func (s *ScanProgress) OnAnalysed(done, total backup.MediaCounter, extra backup.ExtraCounts) {
	s.onAnalysedCalled = true
	if !total.IsZero() {
		s.analysedLine.SetBar(done.Count, total.Count)
		s.analysedLine.SetExplanation(fmt.Sprintf("%d / %d files %s", done.Count, total.Count, extra.String()))

		if done.Count == total.Count {
			s.analysedLine.SwapSpinner(1)
			s.analysedLine.SetLabel("Analyse complete")
		}
	}
}

func (s *ScanProgress) stop() {
	if !s.onAnalysedCalled {
		s.analysedLine.SwapSpinner(1)
		s.analysedLine.SetLabel("Analyse skipped")
	}
	s.screen.Stop()
}
