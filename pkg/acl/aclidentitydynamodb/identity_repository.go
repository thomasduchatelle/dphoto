package aclidentitydynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type IdentityRepository interface {
	aclcore.IdentityDetailsStore
	aclcore.IdentityQueriesIdentityRepository
}

func New(cfg aws.Config, tableName string) (IdentityRepository, error) {
	return &repository{
		db:    dynamodb.NewFromConfig(cfg),
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
	db    *dynamodb.Client
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

	_, err = r.db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      item,
		TableName: &r.table,
	})
	return errors.Wrapf(err, "failed to store %s identity details", identity.Email)
}

func (r *repository) FindIdentity(email usermodel.UserId) (*aclcore.Identity, error) {
	item, err := r.db.GetItem(context.TODO(), &dynamodb.GetItemInput{
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

func (r *repository) FindIdentities(emails []usermodel.UserId) ([]*aclcore.Identity, error) {
	if len(emails) == 0 {
		return nil, nil
	}

	uniqueEmails := make(map[usermodel.UserId]interface{})
	var keys []map[string]types.AttributeValue
	for _, userId := range emails {
		if _, notUnique := uniqueEmails[userId]; !notUnique {
			uniqueEmails[userId] = nil
			keys = append(keys, identityRecordPkAttributes(userId))
		}
	}

	var identities []*aclcore.Identity
	stream := dynamoutils.NewGetStream(context.TODO(), dynamoutils.NewGetBatchItem(r.db, r.table, ""), keys, dynamoutils.DynamoReadBatchSize)
	for stream.HasNext() {
		identity, err := unmarshalIdentity(stream.Next())
		if err != nil {
			return nil, err
		}
		identities = append(identities, identity)
	}

	return identities, errors.Wrapf(stream.Error(), "failed to list identities from provided email.")
}
