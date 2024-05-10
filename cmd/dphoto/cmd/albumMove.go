package cmd

import (
	"context"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

var (
	moveArgs = struct {
		folderName       string
		newName          string
		renameFolder     bool
		forcedFolderName string
	}{}
)

var moveCmd = &cobra.Command{
	Use:     "move <folder name> <new display name> [--update-folder-name]",
	Short:   "Rename an album",
	Long:    `Update the display name of a folder, optionally rename the physical folder name`,
	Aliases: []string{"mv"},
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		moveArgs.folderName = args[0]
		moveArgs.newName = args[1]

		err := pkgfactory.RenameAlbumCase(ctx).RenameAlbum(ctx, catalog.RenameAlbumRequest{
			CurrentId:        catalog.NewAlbumIdFromStrings(Owner, moveArgs.folderName),
			NewName:          moveArgs.newName,
			RenameFolder:     moveArgs.renameFolder,
			ForcedFolderName: moveArgs.forcedFolderName,
		})
		printer.FatalWithMessageIfError(err, 1, "renaming %s failed", moveArgs.folderName)

		printer.Success("Album %s has been renamed.", aurora.Cyan(moveArgs.folderName))
	},
}

func init() {
	albumCmd.AddCommand(moveCmd)

	moveCmd.Flags().BoolVarP(&moveArgs.renameFolder, "update-folder-name", "f", false, "rename physical folder (generated from album name and dates), media will be moved to the new folder")
	moveCmd.Flags().StringVarP(&moveArgs.forcedFolderName, "folder-name", "n", "", "rename physical folder with provided name, medias will be moved to the new folder")
}
