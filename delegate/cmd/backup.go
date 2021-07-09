package cmd

import (
	"duchatelle.io/dphoto/dphoto/backup"
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"duchatelle.io/dphoto/dphoto/cmd/printer"
	"duchatelle.io/dphoto/dphoto/cmd/screen"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"path/filepath"
)

type BackupProgress struct {
	screen       *screen.AutoRefreshScreen
	scanLine     *screen.ProgressLine
	downloadLine *screen.ProgressLine
	analyseLine  *screen.ProgressLine
	uploadLine   *screen.ProgressLine
}

var (
	backupArgs = struct {
		remote bool
	}{}
)

var backupCmd = &cobra.Command{
	Use:   "backup [--remote] <source path>",
	Short: "Backup photos and videos to personal cloud",
	Long:  `Backup photos and videos to personal cloud`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		volumePath, err := filepath.Abs(args[0])
		printer.FatalWithMessageIfError(err, 2, "provided argument must be a valid file path")

		progress := NewProgress()

		tracker, err := backup.StartBackupRunner(backupmodel.VolumeToBackup{
			UniqueId: volumePath,
			Type:     backupmodel.VolumeTypeFileSystem,
			Path:     volumePath,
			Local:    !backupArgs.remote,
		}, progress)
		printer.FatalIfError(err, 1)

		progress.screen.Stop()

		printBackupStats(tracker, volumePath)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)

	backupCmd.Flags().BoolVarP(&backupArgs.remote, "remote", "r", false, "mark the source as remote ; a local buffer will be used to read files only once")
}

func printBackupStats(tracker backupmodel.BackupReport, volumePath string) {
	if len(tracker.CountPerAlbum()) == 0 {
		printer.Success("\n\nBackup of %s complete: %s.", aurora.Cyan(volumePath), aurora.Bold(aurora.Yellow("no new medias")))
		return
	}

	printer.Success("\n\nBackup of %s complete\n", aurora.Cyan(volumePath))

	table := simpletable.New()
	table.Header = &simpletable.Header{Cells: []*simpletable.Cell{
		{Text: "New", Align: simpletable.AlignCenter},
		{Text: "Album", Align: simpletable.AlignCenter},
		{Text: "Photo", Align: simpletable.AlignCenter},
		{Text: "Video", Align: simpletable.AlignCenter},
		{Text: "Total", Align: simpletable.AlignCenter},
	}}
	table.Body = &simpletable.Body{Cells: make([][]*simpletable.Cell, len(tracker.CountPerAlbum())+1)}

	newAlbums := make(map[string]interface{})
	for _, album := range tracker.NewAlbums() {
		newAlbums[album] = nil
	}

	var totals [3]backupmodel.MediaCounter
	i := 0
	for folderName, counts := range tracker.CountPerAlbum() {
		newMarker := ""
		if _, ok := newAlbums[folderName]; ok {
			newMarker = "*"
		}

		table.Body.Cells[i] = []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: newMarker},
			{Text: folderName},
			countAndSize(counts.OfType(backupmodel.MediaTypeImage)),
			countAndSize(counts.OfType(backupmodel.MediaTypeVideo)),
			countAndSize(counts.Total()),
		}

		totals[0] = totals[0].AddCounter(counts.OfType(backupmodel.MediaTypeImage))
		totals[1] = totals[1].AddCounter(counts.OfType(backupmodel.MediaTypeVideo))
		totals[2] = totals[2].AddCounter(counts.Total())
		i++
	}

	table.Body.Cells[len(table.Body.Cells)-1] = []*simpletable.Cell{
		{Text: "Totals", Align: simpletable.AlignRight, Span: 2},
		countAndSize(totals[0]),
		countAndSize(totals[1]),
		countAndSize(totals[2]),
	}

	fmt.Println(table.String())
}

func countAndSize(counter backupmodel.MediaCounter) *simpletable.Cell {
	if counter.Count == 0 {
		return &simpletable.Cell{Align: simpletable.AlignCenter, Text: "-"}
	}

	return &simpletable.Cell{
		Text: fmt.Sprintf("%d (%s)", counter.Count, byteCountIEC(counter.Size)),
	}
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

// binaryMultiplier returns a next power 2 value above given value
func (p *BackupProgress) binaryMultiplier(value uint) int64 {
	nextBinaryPower := int64(2)
	for nextBinaryPower <= int64(value) {
		nextBinaryPower *= 2
	}

	return nextBinaryPower
}

func byteCountIEC(b uint) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
