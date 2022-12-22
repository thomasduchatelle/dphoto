package migrator

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"strings"
)

type TransformationAlbumOwner struct{}

func (t *TransformationAlbumOwner) GeneratePatches(run *TransformationRun, item map[string]*dynamodb.AttributeValue) ([]*dynamodb.WriteRequest, error) {
	pk := *item["PK"].S
	if strings.HasSuffix(pk, "#ALBUM") {
		run.Counter.Inc("ALBUM", 1)

		const albumOwnerKey = "AlbumOwner"
		if owner, ok := item[albumOwnerKey]; !ok || isNull(owner) || isEmpty(owner) {
			run.Counter.Inc("ALBUM_WITHOUT_OWNER", 1)

			item[albumOwnerKey] = &dynamodb.AttributeValue{S: aws.String(strings.TrimSuffix(pk, "#ALBUM"))}
			return []*dynamodb.WriteRequest{
				{
					PutRequest: &dynamodb.PutRequest{
						Item: item,
					},
				},
			}, nil
		}
	}

	return nil, nil
}

func isNull(attr *dynamodb.AttributeValue) bool {
	return attr == nil || attr.NULL != nil && *attr.NULL
}

func isEmpty(attr *dynamodb.AttributeValue) bool {
	return attr == nil || attr.S == nil || *attr.S == ""
}
