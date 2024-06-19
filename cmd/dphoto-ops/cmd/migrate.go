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

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply migration script on storages (DynamoDB, ...)",
	Args:  cobra.NoArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		if migrateCmdArgs.localstack {
			pkgfactory.AWSConfigFactory = awsfactory.NewLocalstackConfig(ctx)
			pkgfactory.AWSNames = &StaticConfigNames{DynamoDBNameValue: "dphoto-local"}

		} else if migrateCmdArgs.terraformOutput != "" {
			terraformCfg, err := awsfactory.ReadTerraformOutputFile(ctx, migrateCmdArgs.terraformOutput)
			printer.FatalIfError(err, 10)

			pkgfactory.AWSConfigFactory = terraformCfg
			pkgfactory.AWSNames = &StaticConfigNames{DynamoDBNameValue: terraformCfg.DynamoDBMainTable()}

		} else {
			printer.ErrorText("Either terraform-output or localstack must be set")
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.PersistentFlags().StringVarP(&migrateCmdArgs.terraformOutput, "terraform-output", "f", "", "(required) output file obtained with 'terraform output -json > output.json'")
	migrateCmd.PersistentFlags().BoolVar(&migrateCmdArgs.localstack, "localstack", false, "set to true when running against localstack DynamoDB")
}

type StaticConfigNames struct {
	DynamoDBNameValue string
}

func (s StaticConfigNames) DynamoDBName() string {
	return s.DynamoDBNameValue
}

func (s StaticConfigNames) ArchiveMainBucketName() string {
	panic("implement me")
}

func (s StaticConfigNames) ArchiveCacheBucketName() string {
	panic("implement me")
}

func (s StaticConfigNames) ArchiveJobsSNSARN() string {
	panic("implement me")
}

func (s StaticConfigNames) ArchiveJobsSQSURL() string {
	panic("implement me")
}
