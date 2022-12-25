package cmd

import (
	"github.com/logrusorgru/aurora/v3"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/tools/dphotoops/migrator"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var migrateArg = struct {
	tableName                string
	snsARN                   string
	repopulate               bool
	indexTransformation      bool
	albumOwnerTransformation bool
}{}

// rootCmd represents the base command when called without any subcommands
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply migration scripts on a DynamoDB table (index upgrades, v1 -> v2, album's owner fix, ...)",
	Run: func(cmd *cobra.Command, args []string) {
		if migrateArg.tableName == "" {
			printer.Error(errors.Errorf("--table is mandatory"), "")
			os.Exit(10)
		}

		printer.Info("Start to migrate table %s ...", aurora.Cyan(migrateArg.tableName))
		start := time.Now()
		count, err := migrator.Migrate(migrateArg.tableName, migrateArg.snsARN, migrateArg.repopulate, buildTransformationsSlice())
		printer.FatalIfError(err, 1)

		elapsed := int(time.Now().Sub(start).Seconds())
		printer.Success("table %s with %d records has been migrated in %d seconds", aurora.Cyan(migrateArg.tableName), count, elapsed)
	},
}

func buildTransformationsSlice() (transformations []interface{}) {
	if migrateArg.indexTransformation {
		transformations = append(transformations, new(migrator.TransformationUpDateIndex))
	}

	if migrateArg.albumOwnerTransformation {
		transformations = append(transformations, new(migrator.TransformationAlbumOwner))
	}

	return
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().StringVarP(&migrateArg.tableName, "table", "t", "", "Table name")
	migrateCmd.Flags().StringVar(&migrateArg.snsARN, "sns", "", "SNS ARN")
	migrateCmd.Flags().BoolVar(&migrateArg.repopulate, "caching", false, "TRUE to populate the cache")

	migrateCmd.Flags().BoolVar(&migrateArg.indexTransformation, "index", false, "update DynamoDB indexes")
	migrateCmd.Flags().BoolVar(&migrateArg.albumOwnerTransformation, "album-owner", false, "add Owner field to albums missing it")
}
