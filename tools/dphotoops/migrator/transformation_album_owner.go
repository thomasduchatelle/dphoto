package migrator

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strings"
)

type TransformationAlbumOwner struct{}

func (t *TransformationAlbumOwner) GeneratePatches(run *TransformationRun, item map[string]types.AttributeValue) ([]*types.WriteRequest, error) {
	if pk, ok := item["PK"].(*types.AttributeValueMemberS); ok && strings.HasSuffix(pk.Value, "#ALBUM") {
		run.Counter.Inc("ALBUM", 1)

		const albumOwnerKey = "AlbumOwner"
		if owner, ok := item[albumOwnerKey]; !ok || isNullOrEmpty(owner) {
			run.Counter.Inc("ALBUM_WITHOUT_OWNER", 1)

			item[albumOwnerKey] = &types.AttributeValueMemberS{Value: strings.TrimSuffix(pk.Value, "#ALBUM")}
			return []*types.WriteRequest{
				{
					PutRequest: &types.PutRequest{
						Item: item,
					},
				},
			}, nil
		}
	}

	return nil, nil
}

func isNullOrEmpty(attr types.AttributeValue) bool {
	_, isNullAttribute := attr.(*types.AttributeValueMemberNULL)
	stringValue, isString := attr.(*types.AttributeValueMemberS)
	return isNullAttribute || isString && stringValue.Value == ""
}
