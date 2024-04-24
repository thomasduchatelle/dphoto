package aclscopedynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
)

func (r *repository) DeleteScopes(ids ...aclcore.ScopeId) error {
	ctx := context.TODO()

	requests := make([]types.WriteRequest, len(ids), len(ids))
	for i, id := range ids {
		requests[i] = types.WriteRequest{DeleteRequest: &types.DeleteRequest{
			Key: MarshalScopeId(id),
		}}
	}

	return errors.Wrapf(
		dynamoutils.BufferedWriteItems(ctx, r.client, requests, r.table, dynamoutils.DynamoWriteBatchSize),
		"failed to delete scopes %+v",
		ids,
	)
}

func (r *repository) SaveIfNewScope(scope aclcore.Scope) error {
	ctx := context.TODO()

	attributes, err := MarshalScope(scope)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
		Item:                attributes,
		TableName:           &r.table,
	})
	var conditionalCheckFailedErr *types.ConditionalCheckFailedException
	if errors.As(err, &conditionalCheckFailedErr) {
		return nil
	}

	return errors.Wrapf(err, "failed to insert scope %+v", scope)
}
