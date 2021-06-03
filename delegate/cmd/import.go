package cmd

import (
	"duchatelle.io/dphoto/dphoto/backup"
	"duchatelle.io/dphoto/dphoto/backup/model"
	"duchatelle.io/dphoto/dphoto/catalog"
	"duchatelle.io/dphoto/dphoto/cmd/interactive"
	"duchatelle.io/dphoto/dphoto/cmd/printer"
	"duchatelle.io/dphoto/dphoto/cmd/screen"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"time"
)

type ScanProgress struct {
	screen       *screen.AutoRefreshScreen
	scanningLine *screen.ProgressLine
	analysedLine *screen.ProgressLine
}

type operations struct{}

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
			fmt.Println()
			printer.Info("%d albums suggested:\n", len(suggestions))

			interactiveSuggestion := make([]*interactive.AlbumRecord, len(suggestions))
			for i, s := range suggestions {
				interactiveSuggestion[i] = &interactive.AlbumRecord{
					Suggestion: true,
					FolderName: "",
					Name:       s.Name,
					Start:      s.Start,
					End:        s.End,
					Count:      0,
					Size:       0,
				}
			}

			interactive.StartInteractiveSession(interactiveSuggestion, new(operations))
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

func (o *operations) FindAllAlbumsWithStats() ([]*interactive.AlbumRecord, error) {
	albums, err := catalog.FindAllAlbumsWithStats()

	records := make([]*interactive.AlbumRecord, len(albums))
	for i, album := range albums {
		records[i] = &interactive.AlbumRecord{
			Suggestion: false,
			FolderName: album.Album.FolderName,
			Name:       album.Album.Name,
			Start:      album.Album.Start,
			End:        album.Album.End,
			Count:      uint(album.TotalCount()),
		}
	}
	return records, err
}

func (o *operations) Create(request interactive.AlbumRecord) error {
	return catalog.Create(catalog.CreateAlbum{
		Name:             request.Name,
		Start:            request.Start,
		End:              request.End,
		ForcedFolderName: request.FolderName,
	})
}

func (o *operations) RenameAlbum(folderName, newName string, renameFolder bool) error {
	return catalog.RenameAlbum(folderName, newName, renameFolder)
}

func (o *operations) UpdateAlbum(folderName string, start, end time.Time) error {
	return catalog.UpdateAlbum(folderName, start, end)
}
