package aclrefreshdynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
)

func New(sess *session.Session, tableName string) (aclcore.RefreshTokenRepository, error) {
	return &repository{
		db:    dynamodb.New(sess),
		table: tableName,
	}, nil
}

func Must(repository aclcore.RefreshTokenRepository, err error) aclcore.RefreshTokenRepository {
	if err != nil {
		panic(err)
	}
	return repository
}

type repository struct {
	db    *dynamodb.DynamoDB
	table string
}

func (r *repository) StoreRefreshToken(token string, spec aclcore.RefreshTokenSpec) error {
	if token == "" {
		return errors.Errorf("refresh token must not be empty")
	}

	item, err := marshalToken(token, spec)
	if err != nil {
		return err
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
		Item:                item,
		TableName:           &r.table,
	})
	return errors.Wrapf(err, "failed to put refresh token %+v", spec)
}

func (r *repository) FindRefreshToken(token string) (*aclcore.RefreshTokenSpec, error) {
	item, err := r.db.GetItem(&dynamodb.GetItemInput{
		Key:       r.refreshTokenKeyAsAttributes(token),
		TableName: &r.table,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find refresh token by its primary key")
	}

	if len(item.Item) == 0 {
		return nil, aclcore.InvalidRefreshTokenError
	}

	return unmarshalToken(item.Item)
}

func (r *repository) DeleteRefreshToken(token string) error {
	_, err := r.db.DeleteItem(&dynamodb.DeleteItemInput{
		Key:       r.refreshTokenKeyAsAttributes(token),
		TableName: &r.table,
	})
	return errors.Wrapf(err, "couldn't delete RefreshToken")
}

func (r *repository) HouseKeepRefreshToken() (int, error) {
	now, err := dynamodbattribute.Marshal(aclcore.TimeFunc())
	if err != nil {
		return 0, errors.Wrapf(err, "couldn't marshal NOW timestamp")
	}
	stream := dynamoutils.NewQueryStream(r.db, []*dynamodb.QueryInput{{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":sk":  {S: aws.String(skValue)},
			":now": now,
		},
		IndexName:              aws.String("RefreshTokenExpiration"),
		KeyConditionExpression: aws.String("SK = :sk AND AbsoluteExpiryTime <= :now"),
		TableName:              &r.table,
	}})

	var requests []*dynamodb.WriteRequest
	for stream.HasNext() {
		item := stream.Next()
		requests = append(requests, &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{Key: map[string]*dynamodb.AttributeValue{
				"PK": item["PK"],
				"SK": item["SK"],
			}},
		})
	}

	if stream.Error() == nil && len(requests) > 0 {
		log.Infof("Removing %d expired refresh tokens...", len(requests))
		err = dynamoutils.BufferedWriteItems(r.db, requests, r.table, dynamoutils.DynamoWriteBatchSize)
		return len(requests), errors.Wrapf(err, "couldn't delete expired refresh tokens")
	}

	return 0, errors.Wrapf(stream.Error(), "couldn't find expired refresh tokens")
}

func (r *repository) refreshTokenKeyAsAttributes(token string) map[string]*dynamodb.AttributeValue {
	pk := RefreshTokenRecordPk(token)
	return map[string]*dynamodb.AttributeValue{
		"PK": {S: &pk.PK},
		"SK": {S: &pk.SK},
	}
}
