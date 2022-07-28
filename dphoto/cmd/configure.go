package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thomasduchatelle/dphoto/dphoto/cmd/ui"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
	"github.com/thomasduchatelle/dphoto/dphoto/printer"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	configureArgs = struct {
		terraformOutput string
	}{}
)

type OutputParam struct {
	Value string `yaml:"value"`
}

type configField struct {
	key         string
	description string
	outputName  string
}

var configureCmd = &cobra.Command{
	Use:   "configure [--terraform-output <JSON output>]",
	Short: "Configuration wizard to grant dphoto access AWS resources.",
	Long: `DPhoto requires specific AWS key and secret, and name of the DynamoDB table and S3 bucket to use.

To set them fro terraform output:

    $ terraform output -json > output.json
    $ terraform output -raw delegate_secret_access_key | base64 --decode | keybase pgp decrypt
    (copy the output, it will be required later)
    $ dphoto configure --terraform-output output.json

The configuration is stored in '~/.dphoto/dphoto.yaml'.
`,
	Run: func(cmd *cobra.Command, args []string) {
		output := make(map[string]OutputParam)

		if configureArgs.terraformOutput != "" {
			content, err := ioutil.ReadFile(configureArgs.terraformOutput)
			printer.FatalIfError(err, 2)

			err = yaml.Unmarshal(content, &output)
			printer.FatalIfError(err, 3)
		}

		fields := []configField{
			{key: config.Owner, description: "Owner of the medias (an email address)"},
			{key: config.AwsRegion, description: "AWS_REGION where dphoto is deployed", outputName: "region"},
			{key: config.AwsKey, description: "AWS_ACCESS_KEY_ID to use with dphoto", outputName: "delegate_access_key_id"},
			{key: config.AwsSecret, description: "AWS_SECRET_ACCESS_KEY to use with dphoto"},
			{key: config.CatalogDynamodbTable, description: "DynamoDB table where catalog is stored", outputName: "dynamodb_name"},
			{key: config.ArchiveDynamodbTable, description: "DynamoDB table where archive index are stored", outputName: "dynamodb_name"},
			{key: config.ArchiveMainBucketName, description: "S3 bucket where medias are archived", outputName: "archive_bucket_name"},
			{key: config.ArchiveCacheBucketName, description: "S3 bucket where medias are cached", outputName: "cache_bucket_name"},
			{key: config.ArchiveJobsSNSARN, description: "SNS ARN where async jobs are dispatched across workers", outputName: "sns_archive_arn"},
			{key: config.ArchiveJobsSQSURL, description: "SQS URL where async jobs are queued and de-duplicated", outputName: "sqs_archive_url"},
		}

		form := ui.NewSimpleForm()
		updated := false
		for _, field := range fields {
			current := viper.GetString(field.key)

			defaultValue := current
			if field.outputName != "" {
				if val, ok := output[field.outputName]; ok {
					defaultValue = val.Value
				}
			}

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

	configureCmd.Flags().StringVarP(&configureArgs.terraformOutput, "terraform-output", "t", "", "File path to terraform output 'terraform output -json > output.json'")
}
