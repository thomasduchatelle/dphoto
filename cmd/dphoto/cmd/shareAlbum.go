package cmd

import (
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/internal/printer"
)

var (
	shareAlbumArg = struct {
		owner      string
		folderName string
		userEmail  string
		revoke     bool
	}{}

	ShareAlbumCase   func(owner, folderName, userEmail string) error
	UnShareAlbumCase func(owner, folderName, userEmail string) error
)

var shareAlbumCmd = &cobra.Command{
	Use:   "share-album",
	Short: "Share [or un-share] an album to a user (different from the the owner)",
	Run: func(cmd *cobra.Command, args []string) {
		if !shareAlbumArg.revoke {
			err := ShareAlbumCase(shareAlbumArg.owner, shareAlbumArg.folderName, shareAlbumArg.userEmail)
			printer.FatalIfError(err, 1)

			printer.Success("Album %s/%s has been shared to %s", aurora.Cyan(shareAlbumArg.owner), aurora.Cyan(shareAlbumArg.folderName), aurora.Cyan(shareAlbumArg.userEmail))
		} else {
			err := UnShareAlbumCase(shareAlbumArg.owner, shareAlbumArg.folderName, shareAlbumArg.userEmail)
			printer.FatalIfError(err, 1)

			printer.Success("Access to album %s/%s has been revoked from %s", aurora.Cyan(shareAlbumArg.owner), aurora.Cyan(shareAlbumArg.folderName), aurora.Cyan(shareAlbumArg.userEmail))
		}
	},
}

func init() {
	rootCmd.AddCommand(shareAlbumCmd)

	shareAlbumCmd.Flags().StringVarP(&shareAlbumArg.owner, "owner", "o", "", "owner of the album to be shared")
	shareAlbumCmd.Flags().StringVarP(&shareAlbumArg.folderName, "album", "a", "", "folder name of the album (expected to start with a /)")
	shareAlbumCmd.Flags().StringVarP(&shareAlbumArg.userEmail, "email", "e", "", "email of the user")
	shareAlbumCmd.Flags().BoolVar(&shareAlbumArg.revoke, "revoke", false, "revoke access instead or granting it")
}
