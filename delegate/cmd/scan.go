package cmd

import (
	"duchatelle.io/dphoto/dphoto/catalog"
	"duchatelle.io/dphoto/dphoto/cmd/adapters/backupadapter"
	"duchatelle.io/dphoto/dphoto/cmd/printer"
	"duchatelle.io/dphoto/dphoto/cmd/ui"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var (
	scanArgs = struct {
		nonInteractive bool
	}{}
)

var scan = &cobra.Command{
	Use:   "scan <folder to scan>",
	Short: "Discover directory structure to suggest new albums to create",
	Long:  "Discover directory structure to suggest new album to create",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		volume := args[0]

		recordRepository, count, err := backupadapter.ScanWithCache(volume)
		printer.FatalIfError(err, 2)

		if count == 0 {
			fmt.Println(aurora.Red(fmt.Sprintf("No media found on path %s .", volume)))
		} else if scanArgs.nonInteractive {
			err = ui.NewSimpleSession(recordRepository, backupadapter.NewAlbumRepository()).Render()
			printer.FatalIfError(err, 1)
		} else {
			err = ui.NewInteractiveSession(new(uiCatalogAdapter), recordRepository, backupadapter.NewAlbumRepository()).Start()
			printer.FatalIfError(err, 1)
		}

		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(scan)

	scan.Flags().BoolVarP(&scanArgs.nonInteractive, "non-interactive", "I", false, "Disable interactive output and only display the scan results.")
}

type uiCatalogAdapter struct{}

func (o uiCatalogAdapter) Create(request ui.RecordCreation) error {
	return catalog.Create(catalog.CreateAlbum{
		Name:             request.Name,
		Start:            request.Start,
		End:              request.End,
		ForcedFolderName: request.FolderName,
	})
}

func (o *uiCatalogAdapter) RenameAlbum(folderName, newName string, renameFolder bool) error {
	return catalog.RenameAlbum(folderName, newName, renameFolder)
}

func (o *uiCatalogAdapter) UpdateAlbum(folderName string, start, end time.Time) error {
	return catalog.UpdateAlbum(folderName, start, end)
}

func (o *uiCatalogAdapter) DeleteAlbum(folderName string) error {
	return catalog.DeleteAlbum(folderName, false)
}
