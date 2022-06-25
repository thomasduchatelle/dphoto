package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/dphoto/cmd/adapters/backupproxy"
	"github.com/thomasduchatelle/dphoto/dphoto/cmd/ui"
	"github.com/thomasduchatelle/dphoto/dphoto/printer"
)

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
			err := ui.NewInteractiveSession(new(uiCatalogAdapter), backupproxy.NewAlbumRepository(Owner), ui.NewNoopRepository(), Owner).Start()
			printer.FatalIfError(err, 1)
		} else {
			err := ui.NewSimpleSession(backupproxy.NewAlbumRepository(Owner), ui.NewNoopRepository()).Render()
			printer.FatalIfError(err, 1)
		}
	},
}

func init() {
	rootCmd.AddCommand(albumCmd)

	albumCmd.Flags().BoolVarP(&listArgs.interactive, "interactive", "i", false, "start an interactive session where albums can be added, deleted, renamed, ...")
}
