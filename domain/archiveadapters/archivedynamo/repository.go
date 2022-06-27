// Package archivedynamo extends catalogdynamo to add media locations to the main table
package archivedynamo

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/catalogdynamo"
)

func New(sess *session.Session, tableName string, createTable bool) (archive.ARepositoryAdapter, error) {
	if createTable {
		_, err := catalogdynamo.NewRepository(sess, tableName)
		if err != nil {
			return nil, err
		}
	}

	return &repository{
		db:    dynamodb.New(sess),
		table: tableName,
	}, nil
}

func Must(a archive.ARepositoryAdapter, err error) archive.ARepositoryAdapter {
	if err != nil {
		panic(err)
	}

	return a
}

type repository struct {
	db    *dynamodb.DynamoDB
	table string
}

func (r *repository) FindById(owner, id string) (string, error) {
	item, err := r.db.GetItem(&dynamodb.GetItemInput{
		Key:       marshalMediaLocationPK(owner, id),
		TableName: &r.table,
	})
	if err != nil {
		return "", errors.Wrapf(err, "FindById %s, %s failed", owner, id)
	}

	if len(item.Item) == 0 {
		return "", archive.NotFoundError
	}

	_, location, err := unmarshalMediaLocation(item.Item)
	return location, err
}

func (r *repository) AddLocation(owner, id, key string) error {
	location, err := marshalMediaLocation(owner, id, key)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal location")
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		Item:      location,
		TableName: &r.table,
	})
	return errors.Wrapf(err, "failed to upsert media location %s - %s - %s", owner, id, key)
}
