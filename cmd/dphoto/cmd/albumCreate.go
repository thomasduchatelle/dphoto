package cmd

import (
	"context"
	"github.com/logrusorgru/aurora/v3"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"time"
)

var (
	newArgs = struct {
		forcedFolderName string
	}{}
)

var newCmd = &cobra.Command{
	Use:   "create <display name> <ISO date> <ISO date> [--folder-name <forced physical name>]",
	Short: "Create a new album",
	Long: `Create a new album

When not specified, folder name is generated from the pattern 'YYYY-MM_<normalised_display_name>'.
`,
	Aliases: []string{"new"},
	Args:    cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		creationRequest := catalog.CreateAlbumRequest{
			Owner:            ownermodel.Owner(Owner),
			Name:             args[0],
			ForcedFolderName: newArgs.forcedFolderName,
		}

		var err error
		creationRequest.Start, err = parseDate(args[1])
		printer.FatalWithMessageIfError(err, 2, "Start date is mandatory")
		creationRequest.End, err = parseDate(args[2])
		printer.FatalWithMessageIfError(err, 3, "End date is mandatory")

		creationRequest.Owner = ownermodel.Owner(Owner)
		_, err = factory.CreateAlbumCase(ctx).Create(ctx, creationRequest)
		printer.FatalWithMessageIfError(err, 1, "Failed to create the album, or to migrate medias to it.")

		printer.Success("Album %s created.", aurora.Cyan(creationRequest.Name))
	},
}

func init() {
	albumCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&newArgs.forcedFolderName, "folder-name", "f", "", "folder name in which medias will be physically stored (optional)")
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
