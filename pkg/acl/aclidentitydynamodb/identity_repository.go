package aclidentitydynamodb

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"strings"
)

type IdentityRepository interface {
	aclcore.IdentityDetailsStore
	aclcore.IdentityQueriesIdentityRepository
}

func New(sess *session.Session, tableName string) (IdentityRepository, error) {
	return &repository{
		db:    dynamodb.New(sess),
		table: tableName,
	}, nil
}

func Must(repository IdentityRepository, err error) IdentityRepository {
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
	item, err := r.db.GetItem(&dynamodb.GetItemInput{
		Key:       identityRecordPkAttributes(email),
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

func (r *repository) FindIdentities(emails []string) ([]*aclcore.Identity, error) {
	if len(emails) == 0 {
		return nil, nil
	}

	uniqueEmails := make(map[string]interface{})
	var keys []map[string]*dynamodb.AttributeValue
	for _, email := range emails {
		email = strings.ToLower(email)
		if _, notUnique := uniqueEmails[email]; !notUnique {
			uniqueEmails[email] = nil
			keys = append(keys, identityRecordPkAttributes(email))
		}
	}

	var identities []*aclcore.Identity
	stream := dynamoutils.NewGetStream(dynamoutils.NewGetBatchItem(r.db, r.table, ""), keys, dynamoutils.DynamoReadBatchSize)
	for stream.HasNext() {
		identity, err := unmarshalIdentity(stream.Next())
		if err != nil {
			return nil, err
		}
		identities = append(identities, identity)
	}

	return identities, errors.Wrapf(stream.Error(), "failed to list identities from provided email.")
}
