package cmd

import (
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/dphoto/cmd/printer"
	"strings"
)

var (
	rmArgs = struct {
		folderName string
		notEmpty   bool
	}{}
)

var removeCmd = &cobra.Command{
	Use:   "remove <folder name> [--not-empty]",
	Short: "Delete an album",
	Long: `Delete an album only if there is no media in it.

Use --not-empty to redispatch medias to other albums before deletion. Default albums (quarters) can not be deleted 
unless different albums fully cover the 3 months.`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"rm"},
	Run: func(cmd *cobra.Command, args []string) {
		rmArgs.folderName = args[0]

		err := catalog.DeleteAlbum(strings.Trim(rmArgs.folderName, "/"), !rmArgs.notEmpty)
		printer.FatalWithMessageIfError(err, 1, "Album %s couldn't be deleted", rmArgs.folderName)

		printer.Success("Album %s has been deleted", aurora.Cyan(rmArgs.folderName))
	},
}

func init() {
	albumCmd.AddCommand(removeCmd)

	removeCmd.Flags().BoolVarP(&rmArgs.notEmpty, "not-empty", "n", false, "redispatch medias to different albums before deletion")
}
