package cmd

import (
	"context"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

var removeCmd = &cobra.Command{
	Use:   "remove <folder name> [--not-empty]",
	Short: "Delete an album",
	Long: `Delete an album and re-distribute medias to other albums.

Albums can only be deleted if all medias can be assigned to a different album.
`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"rm"},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		folderName := args[0]

		deleteCase := pkgfactory.CreateAlbumDeleteCase(ctx)
		err := deleteCase.DeleteAlbum(ctx, catalog.NewAlbumIdFromStrings(Owner, folderName))
		printer.FatalWithMessageIfError(err, 1, "Album %s couldn't be deleted", folderName)

		printer.Success("Album %s has been deleted", aurora.Cyan(folderName))
	},
}

func init() {
	albumCmd.AddCommand(removeCmd)
}
