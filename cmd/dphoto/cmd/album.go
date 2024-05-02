package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/scanui"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/ui"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
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
		ctx := context.TODO()

		if listArgs.interactive {
			err := ui.NewInteractiveSession(&uiCatalogAdapter{
				BackupSuggestionPort: nil,
				CreateAlbum:          pkgfactory.CreateAlbumCase(ctx),
			}, scanui.NewAlbumRepository(Owner), ui.NewNoopRepository(), Owner).Start()
			printer.FatalIfError(err, 1)
		} else {
			err := ui.NewSimpleSession(scanui.NewAlbumRepository(Owner), ui.NewNoopRepository()).Render()
			printer.FatalIfError(err, 1)
		}
	},
}

func init() {
	rootCmd.AddCommand(albumCmd)

	albumCmd.Flags().BoolVarP(&listArgs.interactive, "interactive", "i", false, "start an interactive session where albums can be added, deleted, renamed, ...")
}
