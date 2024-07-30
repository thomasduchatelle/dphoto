package awsfactory

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

const keysJsonContent = `
{
	"region": {
		"sensitive": false,
		"type": "string",
		"value": "us-east-1"
	},
	"delegate_access_key_id": {
		"sensitive": false,
		"type": [
			"object",
			{
				"2024-01": "string"
			}
		],
		"value": {
			"2024-01": "ACCESS_KEY_ID_2024_01"
		}
	},
	"delegate_secret_access_key": {
		"sensitive": false,
		"type": [
			"object",
			{
				"2024-01": "string"
			}
		],
		"value": {
			"2024-01": "a-very-complex-secret"
		}
	},
	"dynamodb_name": {
		"sensitive": false,
		"type": "string",
		"value": "go-unit-dynamodb-table-name-01"
	}
}
`
const terraformOutputJsonContent = `
{
	"archive_bucket_name": {
		"sensitive": false,
		"type": "string",
		"value": "go-unit-main-bucket"
	},
	"cache_bucket_name": {
		"sensitive": false,
		"type": "string",
		"value": "go-unit-cache-bucket"
	},
	"delegate_access_key_id": {
		"sensitive": false,
		"type": [
			"object",
			{
				"2024-01": "string"
			}
		],
		"value": {
			"2024-01": "ACCESS_KEY_ID_2024_01"
		}
	},
	"delegate_secret_access_key": {
		"sensitive": false,
		"type": [
			"object",
			{
				"2024-01": "string"
			}
		],
		"value": {
			"2024-01": "a-very-complex-secret"
		}
	},
	"delegate_secret_access_key_decrypt_cmd": {
		"sensitive": false,
		"type": "string",
		"value": "terraform output -raw delegate_secret_access_key | base64 --decode | keybase pgp decrypt"
	},
	"dynamodb_name": {
		"sensitive": false,
		"type": "string",
		"value": "go-unit-dynamodb-table-name-01"
	},
	"region": {
		"sensitive": false,
		"type": "string",
		"value": "us-east-1"
	},
	"sns_archive_arn": {
		"sensitive": false,
		"type": "string",
		"value": "arn:aws:sns:eu-west-1:000011112222:go-unit-sns-topic-archive"
	},
	"sqs_archive_url": {
		"sensitive": false,
		"type": "string",
		"value": "https://sqs.us-east-1.amazonaws.com/000011112222/go-unit-sqs-topic-archive.fifo"
	},
	"sqs_async_archive_jobs_arn": {
		"sensitive": false,
		"type": "string",
		"value": "arn:aws:sqs:us-east-1:000011112222:go-unit-sqs-topic-archive.fifo"
	}
}
`

func TestParseJSONTerraformOutput(t *testing.T) {
	addsPrefixDecoder := func(encodedValue string) (string, error) {
		return "decoded-" + encodedValue, nil
	}

	type args struct {
		ctx        context.Context
		jsonOutput string
		decoder    SecretAccessKeyDecoder
	}
	type namesStruct struct {
		DynamoDBMainTable string
	}
	tests := []struct {
		name      string
		args      args
		want      *StaticCredentials
		wantNames namesStruct
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "it should extract latest key id and secret from the JSON output",
			args: args{
				ctx:        context.Background(),
				jsonOutput: keysJsonContent,
				decoder:    addsPrefixDecoder,
			},
			want: &StaticCredentials{
				Region:          "us-east-1",
				AccessKeyID:     "ACCESS_KEY_ID_2024_01",
				SecretAccessKey: "decoded-a-very-complex-secret",
			},
			wantNames: namesStruct{DynamoDBMainTable: "go-unit-dynamodb-table-name-01"},
			wantErr:   assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseJsonContentAndDecode(tt.args.ctx, []byte(tt.args.jsonOutput), tt.args.decoder)
			if tt.wantErr(t, err) {
				assert.Equal(t, *tt.want, got.StaticCredentials)
				assert.Equal(t, tt.wantNames.DynamoDBMainTable, got.DynamoDBMainTable())
			}
		})
	}
}
