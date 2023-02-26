package aclidentitydynamodb

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
)

func New(sess *session.Session, tableName string) (aclcore.IdentityDetailsStore, error) {
	return &repository{
		db:    dynamodb.New(sess),
		table: tableName,
	}, nil
}

func Must(repository aclcore.IdentityDetailsStore, err error) aclcore.IdentityDetailsStore {
	if err != nil {
		panic(err)
	}
	return repository
}

type repository struct {
	db    *dynamodb.DynamoDB
	table string
}

func (r *repository) StoreIdentity(identity aclcore.Identity) error {
	if identity.Email == "" {
		return errors.Errorf("email is required to store identity details")
	}

	item, err := marshalIdentity(&identity)
	if err != nil {
		return err
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: &r.table,
	})
	return errors.Wrapf(err, "failed to store %s identity details", identity.Email)
}

func (r *repository) FindIdentity(email string) (*aclcore.Identity, error) {
	pk := IdentityRecordPk(email)

	item, err := r.db.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {S: &pk.PK},
			"SK": {S: &pk.SK},
		},
		TableName: &r.table,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find identity details for %s", email)
	}
	if len(item.Item) == 0 {
		return nil, aclcore.IdentityDetailsNotFoundError
	}

	return unmarshalIdentity(item.Item)
}
