package scoperepository

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/domain/accesscontrol"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/catalogdynamo"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/dynamoutils"
)

type GrantRepository interface {
	accesscontrol.ScopesReader
}

func New(sess *session.Session, tableName string, createTable bool) (GrantRepository, error) {
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

func Must(repository GrantRepository, err error) GrantRepository {
	if err != nil {
		panic(err)
	}
	return repository
}

type repository struct {
	db    *dynamodb.DynamoDB
	table string
}

func (r *repository) ListUserScopes(email string, types ...accesscontrol.ScopeType) ([]*accesscontrol.Scope, error) {
	if len(types) == 0 {
		return nil, nil
	}

	var queries []*dynamodb.QueryInput
	for _, mediaType := range types {
		queries = append(queries, &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":user":     {S: aws.String(userPk(email))},
				":skPrefix": {S: aws.String(fmt.Sprintf("%s%s", scopePrefix, mediaType))},
			},
			KeyConditionExpression: aws.String("PK = :user AND begins_with(SK, :skPrefix)"),
			TableName:              &r.table,
		})
	}

	var scopes []*accesscontrol.Scope
	stream := dynamoutils.NewQueryStream(r.db, queries)
	for stream.HasNext() {
		scope, _, err := UnmarshalScope(stream.Next())
		if err != nil {
			return nil, err
		}

		scopes = append(scopes, scope)
	}

	return scopes, stream.Error()
}
