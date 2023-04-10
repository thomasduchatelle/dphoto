package aclidentitydynamodb

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
)

const (
	identityPrefix = "IDENTITY#"
)

type IdentityRecord struct {
	appdynamodb.TablePk
	Email   string
	Name    string
	Picture string
}

func IdentityRecordPk(user string) appdynamodb.TablePk {
	return appdynamodb.TablePk{
		PK: appdynamodb.UserPk(user),
		SK: identityPrefix,
	}
}

func identityRecordPkAttributes(email string) map[string]*dynamodb.AttributeValue {
	pk := IdentityRecordPk(email)

	key := map[string]*dynamodb.AttributeValue{
		"PK": {S: &pk.PK},
		"SK": {S: &pk.SK},
	}
	return key
}

func marshalIdentity(identity *aclcore.Identity) (map[string]*dynamodb.AttributeValue, error) {
	item, err := dynamodbattribute.MarshalMap(IdentityRecord{
		TablePk: IdentityRecordPk(identity.Email),
		Email:   identity.Email,
		Name:    identity.Name,
		Picture: identity.Picture,
	})
	return item, errors.Wrapf(err, "failed to serialise identity for email %s", identity.Email)
}

func unmarshalIdentity(item map[string]*dynamodb.AttributeValue) (*aclcore.Identity, error) {
	record := new(IdentityRecord)
	err := dynamodbattribute.UnmarshalMap(item, record)

	return &aclcore.Identity{
		Email:   record.Email,
		Name:    record.Name,
		Picture: record.Picture,
	}, errors.Wrapf(err, "failed to unmarshap identity %+v", item)
}
