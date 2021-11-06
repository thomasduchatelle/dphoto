package backupui

import (
	"github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"
	"github.com/thomasduchatelle/dphoto/dphoto/cmd/screen"
	"fmt"
)

type BackupProgress struct {
	screen       *screen.AutoRefreshScreen
	scanLine     *screen.ProgressLine
	downloadLine *screen.ProgressLine
	analyseLine  *screen.ProgressLine
	uploadLine   *screen.ProgressLine
}

func NewProgress() *BackupProgress {
	table := screen.NewTable(" ", 2, 20, 80, 25)

	segments := make([]screen.Segment, 4)
	p := &BackupProgress{}
	p.scanLine, segments[0] = screen.NewProgressLine(table, "Scanning...")
	p.downloadLine, segments[1] = screen.NewProgressLine(table, "Downloading...")
	p.analyseLine, segments[2] = screen.NewProgressLine(table, "Analysing...")
	p.uploadLine, segments[3] = screen.NewProgressLine(table, "Uploading ...")

	p.screen = screen.NewAutoRefreshScreen(
		screen.RenderingOptions{Width: 180},
		segments...,
	)

	return p
}

func (p *BackupProgress) OnScanComplete(total backupmodel.MediaCounter) {
	if total.Count == 0 {
		p.scanLine.SwapSpinner(1)
		p.scanLine.SetLabel(fmt.Sprintf("Scan complete: no new files found"))

		p.downloadLine.SwapSpinner(1)
		p.downloadLine.SetLabel("Download skipped")

		p.analyseLine.SwapSpinner(1)
		p.analyseLine.SetLabel("Analyse skipped")

		p.uploadLine.SwapSpinner(1)
		p.uploadLine.SetLabel("Upload skipped")
	} else {
		p.scanLine.SwapSpinner(1)
		p.scanLine.SetLabel(fmt.Sprintf("Scan complete: %d files found", total.Count))
	}
}

func (p *BackupProgress) OnDownloaded(done, total backupmodel.MediaCounter) {
	if !total.IsZero() {
		p.downloadLine.SetBar(done.Size, total.Size)
		p.downloadLine.SetExplanation(fmt.Sprintf("%s / %s", byteCountIEC(done.Size), byteCountIEC(total.Size)))

		if done.Count == total.Count {
			p.downloadLine.SwapSpinner(1)
			p.downloadLine.SetLabel("Download complete")
		}
	}
}

func (p *BackupProgress) OnAnalysed(done, total backupmodel.MediaCounter) {
	if !total.IsZero() {
		p.analyseLine.SetBar(done.Count, total.Count)
		p.analyseLine.SetExplanation(fmt.Sprintf("%d / %d files", done.Count, total.Count))

		if done.Count == total.Count {
			p.analyseLine.SwapSpinner(1)
			p.analyseLine.SetLabel("Analyse complete")
		}
	}
}

func (p *BackupProgress) OnUploaded(done, total backupmodel.MediaCounter) {
	if !total.IsZero() {
		p.uploadLine.SetBar(done.Size, total.Size)
		p.uploadLine.SetExplanation(fmt.Sprintf("%s / %s", byteCountIEC(done.Size), byteCountIEC(total.Size)))

		if done.Count == total.Count {
			p.uploadLine.SwapSpinner(1)
			p.uploadLine.SetLabel("Upload complete")
		}
	}
}

func (p *BackupProgress) Stop() {
	p.screen.Stop()
}

// binaryMultiplier returns a next power 2 value above given value
func (p *BackupProgress) binaryMultiplier(value uint) int64 {
	nextBinaryPower := int64(2)
	for nextBinaryPower <= int64(value) {
		nextBinaryPower *= 2
	}

	return nextBinaryPower
}
