package archivedynamo

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	dynamodbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"path"
	"strings"
)

type MediaLocationRecord struct {
	appdynamodb.TablePk
	LocationKeyPrefix string // LocationKeyPrefix is used for indexing
	LocationId        string // LocationId is also part of the primary key
	LocationKey       string // LocationKey is the physical location
}

func MediaLocationPk(owner, id string) appdynamodb.TablePk {
	return appdynamodb.TablePk{
		PK: appdynamodb.MediaPrimaryKeyPK(owner, id),
		SK: "LOCATION#",
	}
}

func marshalMediaLocationPK(owner, id string) map[string]dynamodbtypes.AttributeValue {
	pk := MediaLocationPk(owner, id)
	return map[string]dynamodbtypes.AttributeValue{
		"PK": &dynamodbtypes.AttributeValueMemberS{Value: pk.PK},
		"SK": &dynamodbtypes.AttributeValueMemberS{Value: pk.SK},
	}
}

func marshalMediaLocation(owner, id, key string) (map[string]dynamodbtypes.AttributeValue, error) {
	if isBlank(owner) {
		return nil, errors.Errorf("owner is mandatory")
	}
	if isBlank(id) {
		return nil, errors.Errorf("media id is mndatory")
	}

	return attributevalue.MarshalMap(&MediaLocationRecord{
		TablePk:           MediaLocationPk(owner, id),
		LocationKeyPrefix: path.Dir(key),
		LocationId:        id,
		LocationKey:       key,
	})
}

func unmarshalMediaLocation(attributes map[string]dynamodbtypes.AttributeValue) (string, string, error) {
	location := MediaLocationRecord{}
	err := attributevalue.UnmarshalMap(attributes, &location)
	return location.LocationId, location.LocationKey, errors.Wrapf(err, "MediaLocation cannot be unmarchaled from %+v", attributes)
}

// isBlank returns true is value is empty, or contains only spaces
func isBlank(value string) bool {
	return strings.Trim(value, " ") == ""
}
