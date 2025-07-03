package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/thomasduchatelle/dphoto/internal/printer"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/awsfactory"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
)

var (
	migrateCmdArgs = struct {
		terraformOutput string
		localstack      bool
	}{}
)

// FIXME - NOTICE FOR RETIREMENT !!
// FIXME - THIS COMMAND TO UPDATE THE DYNAMODB TABLE STRUCTURE IS REDUNDANT WITH CDK AND WOULD BE RETIRED ONCE CDK IS FULLY MIGRATED.

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply migration script on storages (DynamoDB, ...)",
	Args:  cobra.NoArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		var builder *pkgfactory.AWSCloudBuilder
		if migrateCmdArgs.localstack {
			pkgfactory.StartAWSCloudBuilder(&pkgfactory.StaticAWSAdapterNames{DynamoDBNameValue: "dphoto-local"}).
				OverridesAWSFactory(awsfactory.LocalstackAWSFactory(ctx, awsfactory.LocalstackEndpoint))

		} else if migrateCmdArgs.terraformOutput != "" {
			terraformCfg, err := awsfactory.ReadTerraformOutputFile(ctx, migrateCmdArgs.terraformOutput)
			printer.FatalIfError(err, 10)

			pkgfactory.StartAWSCloudBuilder(&pkgfactory.StaticAWSAdapterNames{DynamoDBNameValue: terraformCfg.DynamoDBMainTable()}).
				OverridesAWSFactory(awsfactory.StaticAWSFactory(ctx, terraformCfg.StaticCredentials))

		} else {
			printer.ErrorText("Either terraform-output or localstack must be set")
		}

		_, err := builder.Build(ctx)
		if err != nil {
			printer.FatalIfError(err, 12)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.PersistentFlags().StringVarP(&migrateCmdArgs.terraformOutput, "terraform-output", "f", "", "(required) output file obtained with 'terraform output -json > output.json'")
	migrateCmd.PersistentFlags().BoolVar(&migrateCmdArgs.localstack, "localstack", false, "set to true when running against localstack DynamoDB")
}
