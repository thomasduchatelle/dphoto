package aclscopedynamodb

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	dynamoutils "github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutilsv2"
)

func (r *repository) ListUserScopes(email string, scopeTypes ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	ctx := context.TODO()

	if len(scopeTypes) == 0 {
		return nil, nil
	}

	var queries []*dynamodb.QueryInput
	for _, scopeType := range scopeTypes {
		expr, err := expression.NewBuilder().WithKeyCondition(
			expression.Key("PK").Equal(expression.Value(appdynamodb.UserPk(email))).
				And(expression.Key("SK").BeginsWith(fmt.Sprintf("%s%s", scopePrefix, scopeType))),
		).Build()
		if err != nil {
			return nil, err
		}

		queries = append(queries, &dynamodb.QueryInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			KeyConditionExpression:    expr.KeyCondition(),
			TableName:                 &r.table,
		})
	}

	var scopes []*aclcore.Scope
	stream := dynamoutils.NewQueryStream(ctx, r.client, queries)
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
	return r.ListScopesByOwners([]string{owner}, types...)
}

func (r *repository) ListScopesByOwners(owners []string, scopeTypes ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	ctx := context.TODO()

	if len(scopeTypes) == 0 {
		return nil, nil
	}

	var queries []*dynamodb.QueryInput
	for _, owner := range owners {
		for _, scopeType := range scopeTypes {
			expr, err := expression.NewBuilder().WithKeyCondition(
				expression.Key("ResourceOwner").Equal(expression.Value(owner)).And(expression.Key("SK").BeginsWith(fmt.Sprintf("%s%s", scopePrefix, scopeType))),
			).Build()
			if err != nil {
				return nil, err
			}

			queries = append(queries, &dynamodb.QueryInput{
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				IndexName:                 aws.String("ReverseGrantIndex"),
				KeyConditionExpression:    expr.KeyCondition(),
				TableName:                 &r.table,
			})
		}
	}

	var scopes []*aclcore.Scope
	stream := dynamoutils.NewQueryStream(ctx, r.client, queries)
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
	ctx := context.TODO()

	keys := make([]map[string]types.AttributeValue, len(ids), len(ids))
	for i, id := range ids {
		keys[i] = MarshalScopeId(id)
	}

	var scopes []*aclcore.Scope
	stream := dynamoutils.NewGetStream(ctx, dynamoutils.NewGetBatchItem(r.client, r.table, ""), keys, dynamoutils.DynamoReadBatchSize)
	for stream.HasNext() {
		scope, err := UnmarshalScope(stream.Next())
		if err != nil {
			return nil, err
		}

		scopes = append(scopes, scope)
	}

	return scopes, stream.Error()
}
