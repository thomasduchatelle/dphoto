package cmd

import (
	"duchatelle.io/dphoto/dphoto/cmd/adapters/backupadapter"
	"duchatelle.io/dphoto/dphoto/cmd/printer"
	"duchatelle.io/dphoto/dphoto/cmd/ui"
	"github.com/spf13/cobra"
)

const ukDateLayout = "02 Jan 06 15:04"

var (
	listArgs = struct {
		interactive bool
	}{}
)

var albumCmd = &cobra.Command{
	Use:     "album [--stats]",
	Aliases: []string{"albums", "alb"},
	Short:   "Organise your collection into albums",
	Long:    `Organise your collection into albums.`,
	Run: func(cmd *cobra.Command, args []string) {
		if listArgs.interactive {
			err := ui.NewInteractiveSession(new(uiCatalogAdapter), backupadapter.NewAlbumRepository()).Start()
			printer.FatalIfError(err, 1)
		} else {
			err := ui.NewSimpleSession(backupadapter.NewAlbumRepository()).Render()
			printer.FatalIfError(err, 1)
		}
	},
}

func init() {
	rootCmd.AddCommand(albumCmd)

	albumCmd.Flags().BoolVarP(&listArgs.interactive, "interactive", "i", false, "start an interactive session where albums can be added, deleted, renamed, ...")
}
