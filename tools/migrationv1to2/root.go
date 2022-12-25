package main

import (
	"github.com/logrusorgru/aurora/v3"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/tools/migrationv1to2/migrator"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	tableName  string
	snsARN     string
	repopulate bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate dynamodb table from its v1 structure to its v2",
	Run: func(cmd *cobra.Command, args []string) {
		if tableName == "" {
			printer.Error(errors.Errorf("--table is mandatory"), "")
			os.Exit(10)
		}

		printer.Info("Start to migrate table %s ...", aurora.Cyan(tableName))
		start := time.Now()
		count, err := migrator.Migrate(tableName, snsARN, repopulate)
		printer.FatalIfError(err, 1)

		elapsed := int(time.Now().Sub(start).Seconds())
		printer.Success("table %s with %d records has been migrated in %d seconds", aurora.Cyan(tableName), count, elapsed)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&tableName, "table", "", "Table name")
	rootCmd.Flags().StringVar(&snsARN, "sns", "", "SNS ARN")
	rootCmd.Flags().BoolVar(&repopulate, "caching", false, "TRUE to populate the cache")
}
