package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/awsbootstrap"
)

var (
	upgradeCmdArgs = struct {
		terraformOutput string
		localstack      bool
	}{}
)
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Check the version of the DynamoDB structure (indexes) and update them if necessary",
	Long: `Check the version of the DynamoDB structure (indexes) and update them if necessary

The script requires direct access to S3 and DynamoDB. It can be done:

    cd deployments/infra-data
	terraform output -json > /tmp/output.json

Then run the command:

	dphoto-ops dynamodb upgrade -f /tmp/output.json

`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		factory, names, err := awsbootstrap.FromTerraformOutputFile(ctx, upgradeCmdArgs.terraformOutput)
		printer.FatalIfError(err, 10)

		client, err := factory.GetDynamoDBClient(ctx)
		printer.FatalIfError(err, 11)

		err = appdynamodb.CreateTableIfNecessary(ctx, names.DynamoDBMainTable(), client, upgradeCmdArgs.localstack)
		printer.FatalIfError(err, 1)
	},
}

func init() {
	dynamodbCmd.AddCommand(upgradeCmd)

	upgradeCmd.Flags().StringVarP(&upgradeCmdArgs.terraformOutput, "terraform-output", "f", "", "(required) output file obtained with 'terraform output -json > output.json'")
	upgradeCmd.Flags().BoolVar(&upgradeCmdArgs.localstack, "localstack", false, "set to true when running against localstack DynamoDB")

}
