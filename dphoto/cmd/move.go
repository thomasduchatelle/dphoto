package cmd

import (
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/dphoto/printer"
	"strings"
)

var (
	moveArgs = struct {
		folderName   string
		newName      string
		renameFolder bool
	}{}
)

var moveCmd = &cobra.Command{
	Use:     "move <folder name> <new display name> [--update-folder-name]",
	Short:   "Rename an album",
	Long:    `Update the display name of a folder, optionally rename the physical folder name`,
	Aliases: []string{"mv"},
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		moveArgs.folderName = args[0]
		moveArgs.newName = args[1]

		err := catalog.RenameAlbum(Owner, strings.Trim(moveArgs.folderName, "/"), moveArgs.newName, moveArgs.renameFolder)
		printer.FatalWithMessageIfError(err, 1, "renaming %s failed", moveArgs.folderName)

		printer.Success("Album %s has been renamed.", aurora.Cyan(moveArgs.folderName))
	},
}

func init() {
	albumCmd.AddCommand(moveCmd)

	moveCmd.Flags().BoolVarP(&moveArgs.renameFolder, "update-folder-name", "f", false, "rename physical folder as well, medias will be moved in the next housekeeping process")
}
