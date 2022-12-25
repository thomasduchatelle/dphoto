package dynamotestutils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"regexp"
	"strings"
	"testing"
	"time"
)

// NewLocalstackSession creates a session to connect to localstack (must be started beforehand)
func NewLocalstackSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(true),
		Endpoint:                      aws.String("http://localhost:4566"),
		Credentials:                   credentials.NewStaticCredentials("localstack", "localstack", ""),
		Region:                        aws.String("eu-west-1"),
	}))
}

func NewTestTableName(t *testing.T) string {
	notValidChar := regexp.MustCompile("[^A-Za-z0-9-]+")
	return strings.ToLower(notValidChar.ReplaceAllString(fmt.Sprintf("%s-%s", t.Name(), time.Now().Format("20060102150405.000")), "-"))
}

// NewDbContext generate AWS session, DynamoDB client, and a table name that will be deleted at the end of the function.
func NewDbContext(t *testing.T) (*session.Session, *dynamodb.DynamoDB, string) {
	awsSession := NewLocalstackSession()
	tableName := NewTestTableName(t)
	db := dynamodb.New(awsSession)

	t.Cleanup(func() {
		_, _ = db.DeleteTable(&dynamodb.DeleteTableInput{TableName: &tableName})
	})

	return awsSession, db, tableName
}
