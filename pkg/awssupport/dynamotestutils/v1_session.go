package dynamotestutils

import (
	awsv1 "github.com/aws/aws-sdk-go/aws"
	credentialsv1 "github.com/aws/aws-sdk-go/aws/credentials"
	sessionv1 "github.com/aws/aws-sdk-go/aws/session"
	dynamodbv1 "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/internal/localstack"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"testing"
)

func NewLocalstackSession() *sessionv1.Session {
	return sessionv1.Must(sessionv1.NewSession(&awsv1.Config{
		CredentialsChainVerboseErrors: awsv1.Bool(true),
		Endpoint:                      awsv1.String(localstack.Endpoint),
		Credentials:                   credentialsv1.NewStaticCredentials("localstack", "localstack", ""),
		Region:                        awsv1.String(localstack.Region),
	}))
}

// NewClientV1 generate AWS session, DynamoDB client, and a table name that will be deleted at the end of the function.
func NewClientV1(t *testing.T) (*sessionv1.Session, *dynamodbv1.DynamoDB, string) {
	awsSession := NewLocalstackSession()
	tableName := NewTestTableName(t)
	db := dynamodbv1.New(awsSession)

	err := appdynamodb.CreateTableIfNecessary(tableName, db, true)
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		_, _ = db.DeleteTable(&dynamodbv1.DeleteTableInput{TableName: &tableName})
	})

	return awsSession, db, tableName
}
