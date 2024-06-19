package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

var indexesCmd = &cobra.Command{
	Use:   "index",
	Short: "Check the version of the DynamoDB structure (indexes) and update them if necessary",
	Long: `Check the version of the DynamoDB structure (indexes) and update them if necessary

The script requires direct access to S3 and DynamoDB. It can be done:

    cd deployments/infra-data
	terraform output -json > /tmp/output.json

Then run the command:

	dphoto-ops dynamodb upgrade -f /tmp/output.json

`,
	Aliases: []string{"indexes"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		err := appdynamodb.CreateTableIfNecessary(ctx, pkgfactory.AWSNames.DynamoDBName(), pkgfactory.AWSFactory(ctx).GetDynamoDBClient(), migrateCmdArgs.localstack)
		printer.FatalIfError(err, 1)
	},
}

func init() {
	migrateCmd.AddCommand(indexesCmd)
}
