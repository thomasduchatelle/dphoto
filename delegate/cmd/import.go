package cmd

import (
	"duchatelle.io/dphoto/dphoto/backup"
	"duchatelle.io/dphoto/dphoto/backup/model"
	"duchatelle.io/dphoto/dphoto/cmd/printer"
	"duchatelle.io/dphoto/dphoto/cmd/screen"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
)

type ScanProgress struct {
	screen       *screen.Screen
	scanningLine *screen.ProgressLine
	analysedLine *screen.ProgressLine
}

const layout = "02/01/2006"

var scan = &cobra.Command{
	Use:   "scan <folder to scan>",
	Short: "Discover directory structure to suggest new album to create",
	Long:  "Discover directory structure to suggest new album to create",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		volume := args[0]

		progress := newScanProgress()
		suggestions, err := backup.DiscoverAlbumFromSource(model.VolumeToBackup{
			UniqueId: volume,
			Type:     model.VolumeTypeFileSystem,
			Path:     volume,
			Local:    true,
		}, progress)
		printer.FatalIfError(err, 1)
		progress.screen.Stop()

		if len(suggestions) == 0 {
			fmt.Println(aurora.Red(fmt.Sprintf("No media found on path %s .", volume)))
		} else {
			table := simpletable.New()
			table.Header = &simpletable.Header{Cells: []*simpletable.Cell{
				{Text: "Name"},
				{Text: "Start"},
				{Text: "End"},
			}}
			table.Body = &simpletable.Body{Cells: make([][]*simpletable.Cell, len(suggestions))}

			for i, album := range suggestions {
				table.Body.Cells[i] = []*simpletable.Cell{
					{Text: aurora.Cyan(album.Name).String()},
					{Text: album.Start.Format(layout)},
					{Text: album.End.Format(layout)},
				}
			}

			fmt.Printf("%d albums suggested:\n", len(suggestions))
			fmt.Println(table)
		}
	},
}

func init() {
	rootCmd.AddCommand(scan)
}

func newScanProgress() *ScanProgress {
	table := screen.NewTable(" ", 2, 20, 80, 25)
	scanningLine, scanningSegment := screen.NewProgressLine(table, "Scanning...")
	analysedLine, analysedSegment := screen.NewProgressLine(table, "Analysed...")

	progressScreen := screen.NewScreen(
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
