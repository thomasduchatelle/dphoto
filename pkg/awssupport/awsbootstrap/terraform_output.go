package awsbootstrap

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"os"
)

type stringValue struct {
	Value string `json:"value,omitempty"`
}

type mapValue struct {
	Value map[string]string `json:"value,omitempty"`
}

type terraformOutput struct {
	Region                  stringValue `json:"region,omitempty"`
	DelegateAccessKeyId     mapValue    `json:"delegate_access_key_id,omitempty"`
	DelegateSecretAccessKey mapValue    `json:"delegate_secret_access_key,omitempty"`
	DynamoDBName            stringValue `json:"dynamodb_name,omitempty"`
}

func (t *terraformOutput) DynamoDBMainTable() string {
	return t.DynamoDBName.Value
}

type SecretAccessKeyDecoder func(encodedValue string) (string, error)

func FromTerraformOutputFile(ctx context.Context, outputFile string) (*DynamoDBFactory, Names, error) {
	content, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Terraform output '%s' is not readable", outputFile)
	}

	cfgFactory, names, err := ParseJSONTerraformOutput(ctx, content, KeybaseDecode)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "terraform output '%s' is invalid", outputFile)
	}

	return &DynamoDBFactory{
		ConfigFactory: cfgFactory,
	}, names, nil
}

func ParseJSONTerraformOutput(ctx context.Context, jsonOutput []byte, decoder SecretAccessKeyDecoder) (*StaticConfig, Names, error) {
	tf := terraformOutput{}
	err := json.Unmarshal(jsonOutput, &tf)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "invalid JSON: [%s]", jsonOutput)
	}

	latestKey := ""
	for key := range tf.DelegateAccessKeyId.Value {
		if key > latestKey {
			latestKey = key
		}
	}

	if latestKey == "" {
		return nil, nil, errors.Errorf("no AWS keys are defined in JSON output.")
	}
	if _, found := tf.DelegateSecretAccessKey.Value[latestKey]; !found {
		return nil, nil, errors.Errorf("key '%s' is not defined on both delegate_access_key_id and delegate_secret_access_key.", latestKey)
	}

	decodedSecretKey, err := decoder(tf.DelegateSecretAccessKey.Value[latestKey])
	return &StaticConfig{
		Region:          tf.Region.Value,
		AccessKeyID:     tf.DelegateAccessKeyId.Value[latestKey],
		SecretAccessKey: decodedSecretKey,
	}, &tf, errors.Wrapf(err, "failed to decode secret key")
}
