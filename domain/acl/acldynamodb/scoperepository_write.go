package acldynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/acl/aclcore"
	"github.com/thomasduchatelle/dphoto/domain/catalogadapters/dynamoutils"
)

func (r *repository) DeleteScopes(ids ...aclcore.ScopeId) error {
	requests := make([]*dynamodb.WriteRequest, len(ids), len(ids))
	for i, id := range ids {
		requests[i] = &dynamodb.WriteRequest{DeleteRequest: &dynamodb.DeleteRequest{
			Key: MarshalScopeId(id),
		}}
	}

	return errors.Wrapf(
		dynamoutils.BufferedWriteItems(r.db, requests, r.table, dynamoutils.DynamoWriteBatchSize),
		"failed to delete scopes %+v",
		ids,
	)
}

func (r *repository) SaveIfNewScope(scope aclcore.Scope) error {
	attributes, err := MarshalScope(scope)
	if err != nil {
		return err
	}

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
		Item:                attributes,
		TableName:           &r.table,
	})
	if aerr, ok := err.(awserr.Error); ok && aerr.Code() == dynamodb.ErrCodeConditionalCheckFailedException {
		return nil
	}

	return errors.Wrapf(err, "failed to insert scope %+v", scope)
}
