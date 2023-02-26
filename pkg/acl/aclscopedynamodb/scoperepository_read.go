package aclscopedynamodb

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	dynamoutils2 "github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
)

func (r *repository) ListUserScopes(email string, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	if len(types) == 0 {
		return nil, nil
	}

	var queries []*dynamodb.QueryInput
	for _, scopeType := range types {
		queries = append(queries, &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":user":     {S: aws.String(appdynamodb.UserPk(email))},
				":skPrefix": {S: aws.String(fmt.Sprintf("%s%s", scopePrefix, scopeType))},
			},
			KeyConditionExpression: aws.String("PK = :user AND begins_with(SK, :skPrefix)"),
			TableName:              &r.table,
		})
	}

	var scopes []*aclcore.Scope
	stream := dynamoutils2.NewQueryStream(r.db, queries)
	for stream.HasNext() {
		scope, err := UnmarshalScope(stream.Next())
		if err != nil {
			return nil, err
		}

		scopes = append(scopes, scope)
	}

	return scopes, stream.Error()
}

func (r *repository) ListOwnerScopes(owner string, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	if len(types) == 0 {
		return nil, nil
	}

	var queries []*dynamodb.QueryInput
	for _, mediaType := range types {
		queries = append(queries, &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":owner":    {S: aws.String(owner)},
				":skPrefix": {S: aws.String(fmt.Sprintf("%s%s", scopePrefix, mediaType))},
			},
			IndexName:              aws.String("ReverseGrantIndex"),
			KeyConditionExpression: aws.String("ResourceOwner = :owner AND begins_with(SK, :skPrefix)"),
			TableName:              &r.table,
		})
	}

	var scopes []*aclcore.Scope
	stream := dynamoutils2.NewQueryStream(r.db, queries)
	for stream.HasNext() {
		scope, err := UnmarshalScope(stream.Next())
		if err != nil {
			return nil, err
		}

		scopes = append(scopes, scope)
	}

	return scopes, stream.Error()
}

func (r *repository) FindScopesById(ids ...aclcore.ScopeId) ([]*aclcore.Scope, error) {
	keys := make([]map[string]*dynamodb.AttributeValue, len(ids), len(ids))
	for i, id := range ids {
		keys[i] = MarshalScopeId(id)
	}

	var scopes []*aclcore.Scope
	stream := dynamoutils2.NewGetStream(dynamoutils2.NewGetBatchItem(r.db, r.table, ""), keys, dynamoutils2.DynamoReadBatchSize)
	for stream.HasNext() {
		scope, err := UnmarshalScope(stream.Next())
		if err != nil {
			return nil, err
		}

		scopes = append(scopes, scope)
	}

	return scopes, stream.Error()
}
