package aclrefreshdynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	dynamoutils "github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutilsv2"
)

func New(cfg aws.Config, tableName string) (aclcore.RefreshTokenRepository, error) {
	return &repository{
		client: dynamodb.NewFromConfig(cfg),
		table:  tableName,
	}, nil
}

func Must(repository aclcore.RefreshTokenRepository, err error) aclcore.RefreshTokenRepository {
	if err != nil {
		panic(err)
	}
	return repository
}

type repository struct {
	client *dynamodb.Client
	table  string
}

func (r *repository) StoreRefreshToken(token string, spec aclcore.RefreshTokenSpec) error {
	if token == "" {
		return errors.Errorf("refresh token must not be empty")
	}

	item, err := marshalToken(token, spec)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
		Item:                item,
		TableName:           &r.table,
	})
	return errors.Wrapf(err, "failed to put refresh token %+v", spec)
}

func (r *repository) FindRefreshToken(token string) (*aclcore.RefreshTokenSpec, error) {
	item, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
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
	_, err := r.client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key:       r.refreshTokenKeyAsAttributes(token),
		TableName: &r.table,
	})
	return errors.Wrapf(err, "couldn't delete RefreshToken")
}

func (r *repository) HouseKeepRefreshToken() (int, error) {
	ctx := context.TODO()

	now, err := attributevalue.Marshal(aclcore.TimeFunc())
	if err != nil {
		return 0, errors.Wrapf(err, "couldn't marshal NOW timestamp")
	}

	expr, err := expression.NewBuilder().WithKeyCondition(
		expression.Key("SK").Equal(expression.Value(skValue)).
			And(expression.Key("AbsoluteExpiryTime").LessThanEqual(expression.Value(now))),
	).Build()
	if err != nil {
		return 0, err
	}

	stream := dynamoutils.NewQueryStream(ctx, r.client, []*dynamodb.QueryInput{{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		IndexName:                 aws.String("RefreshTokenExpiration"),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 &r.table,
	}})

	var requests []types.WriteRequest
	for stream.HasNext() {
		item := stream.Next()
		requests = append(requests, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{Key: map[string]types.AttributeValue{
				"PK": item["PK"],
				"SK": item["SK"],
			}},
		})
	}

	if stream.Error() == nil && len(requests) > 0 {
		log.Infof("Removing %d expired refresh tokens...", len(requests))
		err = dynamoutils.BufferedWriteItems(ctx, r.client, requests, r.table, dynamoutils.DynamoWriteBatchSize)
		return len(requests), errors.Wrapf(err, "couldn't delete expired refresh tokens")
	}

	return 0, errors.Wrapf(stream.Error(), "couldn't find expired refresh tokens")
}

func (r *repository) refreshTokenKeyAsAttributes(token string) map[string]types.AttributeValue {
	pk := RefreshTokenRecordPk(token)
	return map[string]types.AttributeValue{
		"PK": dynamoutils.AttributeValueMemberS(pk.PK),
		"SK": dynamoutils.AttributeValueMemberS(pk.SK),
	}
}
