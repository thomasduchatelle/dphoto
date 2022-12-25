package backupui

import (
	"fmt"
	screen3 "github.com/thomasduchatelle/dphoto/internal/screen"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
)

type BackupProgress struct {
	screen           *screen3.AutoRefreshScreen
	scanLine         *screen3.ProgressLine
	analyseLine      *screen3.ProgressLine
	uploadLine       *screen3.ProgressLine
	onAnalysedCalled bool
	onUploadCalled   bool
}

func NewProgress() *BackupProgress {
	table := screen3.NewTable(" ", 2, 20, 80, 25)

	segments := make([]screen3.Segment, 3)
	p := &BackupProgress{}
	p.scanLine, segments[0] = screen3.NewProgressLine(table, "Scanning...")
	p.analyseLine, segments[1] = screen3.NewProgressLine(table, "Analysing...")
	p.uploadLine, segments[2] = screen3.NewProgressLine(table, "Uploading ...")

	p.screen = screen3.NewAutoRefreshScreen(
		screen3.RenderingOptions{Width: 180},
		segments...,
	)

	return p
}

func (p *BackupProgress) OnScanComplete(total backup.MediaCounter) {
	if total.Count == 0 {
		p.scanLine.SwapSpinner(1)
		p.scanLine.SetLabel(fmt.Sprintf("Scan complete: no new files found"))
	} else {
		p.scanLine.SwapSpinner(1)
		p.scanLine.SetLabel(fmt.Sprintf("Scan complete: %d files found", total.Count))
	}
}

func (p *BackupProgress) OnAnalysed(done, total backup.MediaCounter) {
	p.onAnalysedCalled = true
	if !total.IsZero() {
		p.analyseLine.SetBar(done.Count, total.Count)
		p.analyseLine.SetExplanation(fmt.Sprintf("%d / %d files", done.Count, total.Count))

		if done.Count == total.Count {
			p.analyseLine.SwapSpinner(1)
			p.analyseLine.SetLabel("Analyse complete")
		}
	}
}

func (p *BackupProgress) OnUploaded(done, total backup.MediaCounter) {
	p.onUploadCalled = true
	if !total.IsZero() {
		p.uploadLine.SetBar(done.Size, total.Size)
		p.uploadLine.SetExplanation(fmt.Sprintf("%s / %s", byteCountIEC(done.Size), byteCountIEC(total.Size)))

		if done.Count == total.Count {
			p.uploadLine.SwapSpinner(1)
			p.uploadLine.SetLabel("Upload complete")
		}
	} else {
		p.uploadLine.SetExplanation(fmt.Sprintf("%s", byteCountIEC(done.Size)))
	}
}

func (p *BackupProgress) Stop() {
	if !p.onAnalysedCalled {
		p.analyseLine.SwapSpinner(1)
		p.analyseLine.SetLabel("Analyse skipped")
	}
	if !p.onUploadCalled {
		p.uploadLine.SwapSpinner(1)
		p.uploadLine.SetLabel("Upload skipped")
	}
	p.screen.Stop()
}
