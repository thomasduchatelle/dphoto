package cmd

import (
	"github.com/logrusorgru/aurora/v3"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"github.com/thomasduchatelle/dphoto/dphoto/cmd/printer"
	"time"
)

var (
	creationRequest catalog.CreateAlbum
	newArgs         = struct {
		startDate string
		endDate   string
	}{}
)

var newCmd = &cobra.Command{
	Use:   "new --name <display name> --start <ISO date> --end <ISO date> [--folder-name <forced physical name>]",
	Short: "Create a new album",
	Long: `Create a new album

When not specified, folder name is generated from the pattern 'YYYY-MM_<normalised_display_name>'.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		creationRequest.Start, err = parseDate(newArgs.startDate)
		printer.FatalWithMessageIfError(err, 2, "Start date is mandatory")
		creationRequest.End, err = parseDate(newArgs.endDate)
		printer.FatalWithMessageIfError(err, 3, "End date is mandatory")

		err = catalog.Create(creationRequest)
		printer.FatalWithMessageIfError(err, 1, "Failed to create the album, or to migrate medias to it.")

		printer.Success("Album %s created.", aurora.Cyan(creationRequest.Name))
	},
}

func init() {
	albumCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&creationRequest.Name, "name", "n", "", "name of the album, mandatory")
	newCmd.Flags().StringVarP(&creationRequest.ForcedFolderName, "folder-name", "f", "", "folder name in which medias will be physically stored (optional)")

	newCmd.Flags().StringVarP(&newArgs.startDate, "start", "s", "", "start date, format: YYYY-MM-DD ; or YYYY-MM-DDTHH:mm:SS")
	newCmd.Flags().StringVarP(&newArgs.endDate, "end", "e", "", "end date, format: YYYY-MM-DD ; or YYYY-MM-DDTHH:mm:SS")
}

func parseDate(value string) (time.Time, error) {
	for _, layout := range []string{"2006-01-02T15:04:05", "2006-01-02"} {
		parse, err := time.Parse(layout, value)
		if err == nil {
			return parse, nil
		}
	}

	return time.Time{}, errors.Errorf("'%s' is not a valid date, or datetime, format.", value)
}
