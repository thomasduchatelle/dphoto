package aclidentitydynamodb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func identityRecordPkAttributes(email string) map[string]types.AttributeValue {
	pk := IdentityRecordPk(email)

	key := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: pk.PK},
		"SK": &types.AttributeValueMemberS{Value: pk.SK},
	}
	return key
}

func marshalIdentity(identity *aclcore.Identity) (map[string]types.AttributeValue, error) {
	item, err := attributevalue.MarshalMap(IdentityRecord{
		TablePk: IdentityRecordPk(identity.Email),
		Email:   identity.Email,
		Name:    identity.Name,
		Picture: identity.Picture,
	})
	return item, errors.Wrapf(err, "failed to serialise identity for email %s", identity.Email)
}

func unmarshalIdentity(item map[string]types.AttributeValue) (*aclcore.Identity, error) {
	record := new(IdentityRecord)
	err := attributevalue.UnmarshalMap(item, record)

	return &aclcore.Identity{
		Email:   record.Email,
		Name:    record.Name,
		Picture: record.Picture,
	}, errors.Wrapf(err, "failed to unmarshap identity %+v", item)
}
