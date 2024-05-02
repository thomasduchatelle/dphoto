package awsfactory

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
)

type stringValue struct {
	Value string `json:"value,omitempty"`
}

type mapValue struct {
	Value map[string]string `json:"value,omitempty"`
}

// TerraformOutput is the content of `infra-data` output and can be used as Names and as ConfigFactory.
type TerraformOutput struct {
	StaticConfig
	DynamoDBName           string
	ArchiveMainBucketName  string
	ArchiveCacheBucketName string
	ArchiveJobsSNSARN      string
	ArchiveJobsSQSURL      string
}

type output struct {
	Region                  stringValue `json:"region,omitempty"`
	DelegateAccessKeyId     mapValue    `json:"delegate_access_key_id,omitempty"`
	DelegateSecretAccessKey mapValue    `json:"delegate_secret_access_key,omitempty"`
	DynamoDBName            stringValue `json:"dynamodb_name,omitempty"`
	ArchiveMainBucketName   stringValue `json:"archive_bucket_name,omitempty"`
	ArchiveCacheBucketName  stringValue `json:"cache_bucket_name,omitempty"`
	ArchiveJobsSNSARN       stringValue `json:"sns_archive_arn,omitempty"`
	ArchiveJobsSQSURL       stringValue `json:"sqs_archive_url,omitempty"`
}

func (t *TerraformOutput) DynamoDBMainTable() string {
	return t.DynamoDBName
}

type SecretAccessKeyDecoder func(encodedValue string) (string, error)

// ReadTerraformOutputFile read the `terraform output -json` file as base to configure the application.
func ReadTerraformOutputFile(ctx context.Context, outputFile string) (*TerraformOutput, error) {
	content, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, errors.Wrapf(err, "Terraform output '%s' is not readable", outputFile)
	}

	return parseJsonContentAndDecode(ctx, content, keybaseDecode)
}

func parseJsonContentAndDecode(ctx context.Context, jsonOutput []byte, decoder SecretAccessKeyDecoder) (*TerraformOutput, error) {
	tf := output{}
	err := json.Unmarshal(jsonOutput, &tf)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid JSON: [%s]", jsonOutput)
	}

	latestKey := ""
	for key := range tf.DelegateAccessKeyId.Value {
		if key > latestKey {
			latestKey = key
		}
	}

	if latestKey == "" {
		return nil, errors.Errorf("no AWS keys are defined in JSON output.")
	}
	if _, found := tf.DelegateSecretAccessKey.Value[latestKey]; !found {
		return nil, errors.Errorf("key '%s' is not defined on both delegate_access_key_id and delegate_secret_access_key.", latestKey)
	}

	decodedSecretAccessKey, err := decoder(tf.DelegateSecretAccessKey.Value[latestKey])

	return &TerraformOutput{
		StaticConfig: StaticConfig{
			Region:          tf.Region.Value,
			AccessKeyID:     tf.DelegateAccessKeyId.Value[latestKey],
			SecretAccessKey: decodedSecretAccessKey,
		},
		DynamoDBName:           tf.DynamoDBName.Value,
		ArchiveMainBucketName:  tf.ArchiveMainBucketName.Value,
		ArchiveCacheBucketName: tf.ArchiveCacheBucketName.Value,
		ArchiveJobsSNSARN:      tf.ArchiveJobsSNSARN.Value,
		ArchiveJobsSQSURL:      tf.ArchiveJobsSQSURL.Value,
	}, errors.Wrapf(err, "failed to decode secret key")
}

func keybaseDecode(b64Secret string) (string, error) {
	cmd := fmt.Sprintf("echo '%s' | base64 --decode | keybase pgp decrypt", b64Secret)
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	return string(out), errors.Wrapf(err, "keybase failed to decrypt the secret key [%s]", out)
}
