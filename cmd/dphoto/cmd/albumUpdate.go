package cmd

import (
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

var (
	updateArgs = struct {
		folderName string
	}{}
)

var updateCmd = &cobra.Command{
	Use:   "dates <folderName> <ISO date> <ISO date>",
	Short: "Update the dates of the album, and redispatch the medias accordingly",
	Long: `Update the dates of the album, and redispatch the medias accordingly.

Dates format is: YYYY-MM-DD ; or YYYY-MM-DDTHH:mm:SS
Note: default quarter albums should not be updated unless the 3 months are covered by other albums.
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		updateArgs.folderName = args[0]

		startDate, err := parseDate(args[1])
		printer.FatalWithMessageIfError(err, 2, "Start date is mandatory")
		endDate, err := parseDate(args[2])
		printer.FatalWithMessageIfError(err, 3, "End date is mandatory")

		err = factory.AmendAlbumDatesCase(ctx).AmendAlbumDates(ctx, catalog.NewAlbumIdFromStrings(Owner, updateArgs.folderName), startDate, endDate)
		printer.FatalWithMessageIfError(err, 1, "Couldn't update dates of folder %s", updateArgs.folderName)

		printer.Success("Album %s has been updated.", aurora.Cyan(updateArgs.folderName))
	},
}

func init() {
	albumCmd.AddCommand(updateCmd)
}
