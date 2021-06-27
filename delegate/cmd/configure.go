package cmd

import (
	"duchatelle.io/dphoto/dphoto/cmd/printer"
	"duchatelle.io/dphoto/dphoto/cmd/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			{key: "aws.key", description: "AWS_ACCESS_KEY_ID to use with dphoto", outputName: "delegate_access_key_id"},
			{key: "aws.secret", description: "AWS_SECRET_ACCESS_KEY to use with dphoto"},
			{key: "aws.region", description: "AWS_REGION where dphoto is deployed"},
			{key: "catalog.dynamodb.table", description: "DynamoDB table where catalog is stored", outputName: "dynamodb_name"},
			{key: "backup.s3.bucket", description: "S3 bucket where medias are stored", outputName: "bucket_name"},
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
	//{
	//	"bucket_name": {
	//	"sensitive": false,
	//		"type": "string",
	//		"value": "dphoto-dev-storage"
	//},
	//	"delegate_access_key_id": {
	//	"sensitive": false,
	//		"type": "string",
	//		"value": "AKIAW32ATO5LHFCC6UWF"
	//},
	//	"delegate_secret_access_key": {
	//	"sensitive": false,
	//		"type": "string",
	//		"value": "wcFMA+4kFetGXavOARAA0g20iaNiIXG1kOeFrYlyGxRqE+xTXX/V6TPU4HpEK7uNNzclYqb7OR810H0S6CzmPHo1tpmXMJi4eVNpee8oilnseGlkZjsMIHTMM52pQ0GZCox4gMbcPDN+FYVR1pp3jatBoQda49lvKhWXTZjW4LzOw9uM7HhGXKIMr9VThC30nFFLJttgPhmtQR/oEZMtC36gnJMH+QuzZ0GnWZ1tuXadaqW5u0PpY+ZlfJV/+qHkRBpqOspp65OSpfR1ctA4g0dKa7hfVNNSyN7LbVLtRGHtWtXmTAYqny/nA6u3+DHHhNj79EH3Sd3mVCi4SC7QaATmi9k97knL7QqxN/+uKuc2Yi3Yn6sW+8aXPr9LQg0n+SBS7Ic8v04llLWQbKe0QoVw9TVvSr0DkY2P+yF3GmO3hOW2uQgcI4ARjI4BA8bBu5iIJxIJfL8DNxv3gQig2QEB5Y1u5m57z6RwQ+GfcsP8npgdPiY+wRx3XIc2fl3t/fH/H3G29urN9Si6/tGqJmg8mFxu7Qvv0mMg39KDZBIPnUfa9lRoGoHXRYKnS8GyNb4zI/7WIpIfFje1wHT/jpuxrs9zO+xkNETZ16iNapjckb0gOBm15Qwzn6RKTRKQAaGUevar4KpUHsUermvs8raJuNY0k/VgDG8gt1Wofb2N7ebmk3gpGYivMMAxuijS4AHkt6U5PbL/IYMfmqT5WgHi1eHdfODO4L3ho1fg5OJj46JW4Nbly3Fs1gSBjH4uE0Xuv9Vi8l88GLmSnAeRcGn8VtbfAqbg7+OKoTiFXwbmEOCc5P14G1AyTAdrfrOMrN4nf+DiNfm2kOEm/gA="
	//},
	//	"delegate_secret_access_key_decrypt_cmd": {
	//	"sensitive": false,
	//		"type": "string",
	//		"value": "terraform output -raw delegate_secret_access_key | base64 --decode | keybase pgp decrypt"
	//},
	//	"dynamodb_name": {
	//	"sensitive": false,
	//		"type": "string",
	//		"value": "dphoto-dev-index"
	//}
	//}

}
