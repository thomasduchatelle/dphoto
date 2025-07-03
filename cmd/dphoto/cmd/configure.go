package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/cmd/ui"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/internal/printer"
)

type configField struct {
	key         string
	description string
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configuration wizard to grant dphoto access AWS resources.",
	Long: `DPhoto requires specific AWS key and secret, and name of the DynamoDB table and S3 bucket to use.

The configuration is stored in '~/.dphoto/dphoto.yaml'.
`,
	Run: func(cmd *cobra.Command, args []string) {
		defaultValues := make(map[string]string)

		fields := []configField{
			{key: config.Owner, description: "Owner of the medias (an email address)"},
			{key: config.AwsRegion, description: "AWS_REGION where dphoto is deployed"},
			{key: config.AwsKey, description: "AWS_ACCESS_KEY_ID to use with dphoto"},
			{key: config.AwsSecret, description: "AWS_SECRET_ACCESS_KEY to use with dphoto"},
			{key: config.CatalogDynamodbTable, description: "DynamoDB table where catalog is stored"},
			{key: config.ArchiveDynamodbTable, description: "DynamoDB table where archive index are stored"},
			{key: config.ArchiveMainBucketName, description: "S3 bucket where medias are archived"},
			{key: config.ArchiveCacheBucketName, description: "S3 bucket where medias are cached"},
			{key: config.ArchiveJobsSNSARN, description: "SNS ARN where async jobs are dispatched across workers"},
			{key: config.ArchiveJobsSQSURL, description: "SQS URL where async jobs are queued and de-duplicated"},
		}

		form := ui.NewSimpleForm()
		updated := false
		for _, field := range fields {
			current := viper.GetString(field.key)

			defaultValue, _ := defaultValues[field.key]

			if read, ok := form.ReadString(field.description, defaultValue); ok && read != current {
				viper.Set(field.key, read)
				updated = true
			}
		}

		if updated {
			err := viper.WriteConfig()
			printer.FatalIfError(err, 1)
		}
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
