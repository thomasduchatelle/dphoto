package cmd

import (
	"duchatelle.io/dphoto/dphoto/backup"
	tracker2 "duchatelle.io/dphoto/dphoto/backup/interactors/tracker"
	"duchatelle.io/dphoto/dphoto/backup/model"
	"duchatelle.io/dphoto/dphoto/cmd/printer"
	"duchatelle.io/dphoto/dphoto/cmd/screen"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"path/filepath"
)

type ProgressLine struct {
	swapSpinner    func(int)
	setLabel       func(string)
	setBar         func(uint, uint)
	setExplanation func(string)
}

type Progress struct {
	screen       *screen.Screen
	scanLine     *ProgressLine
	downloadLine *ProgressLine
	analyseLine  *ProgressLine
	uploadLine   *ProgressLine
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

		tracker, err := backup.StartBackupRunner(model.VolumeToBackup{
			UniqueId: volumePath,
			Type:     model.VolumeTypeFileSystem,
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

func printBackupStats(tracker *tracker2.Tracker, volumePath string) {
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

	var totals [3]tracker2.MediaCounter
	i := 0
	for folderName, counts := range tracker.CountPerAlbum() {
		newMarker := ""
		if _, ok := newAlbums[folderName]; ok {
			newMarker = "*"
		}

		table.Body.Cells[i] = []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: newMarker},
			{Text: folderName},
			countAndSize(counts.OfType(model.MediaTypeImage)),
			countAndSize(counts.OfType(model.MediaTypeVideo)),
			countAndSize(counts.Total()),
		}

		totals[0] = totals[0].AddCounter(counts.OfType(model.MediaTypeImage))
		totals[1] = totals[1].AddCounter(counts.OfType(model.MediaTypeVideo))
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

func countAndSize(counter tracker2.MediaCounter) *simpletable.Cell {
	if counter.Count == 0 {
		return &simpletable.Cell{Align: simpletable.AlignCenter, Text: "-"}
	}

	return &simpletable.Cell{
		Text: fmt.Sprintf("%d (%s)", counter.Count, byteCountIEC(counter.Size)),
	}
}

func NewProgress() *Progress {
	table := screen.NewTable(" ", 2, 20, 80, 25)

	segments := make([]screen.Segment, 4)
	p := &Progress{}
	p.scanLine, segments[0] = NewProgressLine(table, "Scanning...")
	p.downloadLine, segments[1] = NewProgressLine(table, "Downloading...")
	p.analyseLine, segments[2] = NewProgressLine(table, "Analysing...")
	p.uploadLine, segments[3] = NewProgressLine(table, "Uploading ...")

	p.screen = screen.NewScreen(
		screen.RenderingOptions{Width: 180},
		segments...,
	)

	return p
}

func NewProgressLine(table *screen.TableGenerator, initialLabel string) (*ProgressLine, screen.Segment) {
	spinner, swapSpinner := screen.NewSwitchSegment(screen.NewSpinnerSegment(), screen.NewGreenTickSegment())
	label, setLabel := screen.NewUpdatableSegment(initialLabel)
	bar, setBar := screen.NewProgressBarSegment()
	explanation, setExplanation := screen.NewUpdatableSegment("")

	return &ProgressLine{
		swapSpinner:    swapSpinner,
		setLabel:       setLabel,
		setBar:         setBar,
		setExplanation: setExplanation,
	}, table.NewRow(spinner, label, bar, explanation)
}

func (p *Progress) OnScanComplete(total tracker2.MediaCounter) {
	p.scanLine.swapSpinner(1)
	p.scanLine.setLabel(fmt.Sprintf("Scan complete: %d files found", total.Count))
}

func (p *Progress) OnDownloaded(done, total tracker2.MediaCounter) {
	if !total.IsZero() {
		p.downloadLine.setBar(done.Size, total.Size)
		p.downloadLine.setExplanation(fmt.Sprintf("%s / %s", byteCountIEC(done.Size), byteCountIEC(total.Size)))

		if done.Count == total.Count {
			p.downloadLine.swapSpinner(1)
			p.downloadLine.setLabel("Download complete")
		}
	}
}

func (p *Progress) OnAnalysed(done, total tracker2.MediaCounter) {
	//time.Sleep(330 * time.Millisecond)
	if !total.IsZero() {
		p.analyseLine.setBar(done.Count, total.Count)
		p.analyseLine.setExplanation(fmt.Sprintf("%d / %d files", done.Count, total.Count))

		if done.Count == total.Count {
			p.analyseLine.swapSpinner(1)
			p.analyseLine.setLabel("Analyse complete")
		}
	}
}

func (p *Progress) OnUploaded(done, total tracker2.MediaCounter) {
	//time.Sleep(time.Second)
	if !total.IsZero() {
		p.uploadLine.setBar(done.Size, total.Size)
		p.uploadLine.setExplanation(fmt.Sprintf("%s / %s", byteCountIEC(done.Size), byteCountIEC(total.Size)))

		if done.Count == total.Count {
			p.uploadLine.swapSpinner(1)
			p.uploadLine.setLabel("Upload complete")
		}
	}
}

// binaryMultiplier returns a next power 2 value above given value
func (p *Progress) binaryMultiplier(value uint) int64 {
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
