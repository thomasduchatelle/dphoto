package archivedynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

func marshalMediaLocationPK(owner, id string) map[string]*dynamodb.AttributeValue {
	pk := MediaLocationPk(owner, id)
	return map[string]*dynamodb.AttributeValue{
		"PK": {S: aws.String(pk.PK)},
		"SK": {S: aws.String(pk.SK)},
	}
}

func marshalMediaLocation(owner, id, key string) (map[string]*dynamodb.AttributeValue, error) {
	if isBlank(owner) {
		return nil, errors.Errorf("owner is mandatory")
	}
	if isBlank(id) {
		return nil, errors.Errorf("media id is mndatory")
	}

	return dynamodbattribute.MarshalMap(&MediaLocationRecord{
		TablePk:           MediaLocationPk(owner, id),
		LocationKeyPrefix: path.Dir(key),
		LocationId:        id,
		LocationKey:       key,
	})
}

func unmarshalMediaLocation(attributes map[string]*dynamodb.AttributeValue) (string, string, error) {
	location := MediaLocationRecord{}
	err := dynamodbattribute.UnmarshalMap(attributes, &location)
	return location.LocationId, location.LocationKey, errors.Wrapf(err, "MediaLocation cannot be unmarchaled from %+v", attributes)
}

// isBlank returns true is value is empty, or contains only spaces
func isBlank(value string) bool {
	return strings.Trim(value, " ") == ""
}
