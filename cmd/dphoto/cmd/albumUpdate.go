package cmd

import (
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

var (
	updateArgs = struct {
		folderName string
		start      string
		end        string
	}{}
)

var updateCmd = &cobra.Command{
	Use:   "update <folderName> --start <ISO date> --end <ISO date>",
	Short: "Update the date of the albums, and redispatch albums accordingly",
	Long: `Update the date of the albums, and redispatch albums accordingly.

Note: default quarter albums should not be updated unless the 3 months are covered by other albums.
`,
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"up"},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		updateArgs.folderName = args[0]

		startDate, err := parseDate(updateArgs.start)
		printer.FatalWithMessageIfError(err, 2, "Start date is mandatory")
		endDate, err := parseDate(updateArgs.end)
		printer.FatalWithMessageIfError(err, 3, "End date is mandatory")

		err = pkgfactory.AmendAlbumDatesCase(ctx).AmendAlbumDates(ctx, catalog.NewAlbumIdFromStrings(Owner, updateArgs.folderName), startDate, endDate)
		printer.FatalWithMessageIfError(err, 1, "Couldn't update dates of folder %s", updateArgs.folderName)

		printer.Success("Album %s has been updated.", aurora.Cyan(updateArgs.folderName))
	},
}

func init() {
	albumCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updateArgs.start, "start", "s", "", "start date, format: YYYY-MM-DD ; or YYYY-MM-DDTHH:mm:SS")
	updateCmd.Flags().StringVarP(&updateArgs.end, "end", "e", "", "end date, format: YYYY-MM-DD ; or YYYY-MM-DDTHH:mm:SS")
}
