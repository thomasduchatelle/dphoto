package dynamotestutils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	awsv1 "github.com/aws/aws-sdk-go/aws"
	credentialsv1 "github.com/aws/aws-sdk-go/aws/credentials"
	sessionv1 "github.com/aws/aws-sdk-go/aws/session"
	dynamodbv1 "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"regexp"
	"strings"
	"testing"
	"time"
)

const (
	region   = "us-east-1"
	endpoint = "http://localhost:4566"
)

func NewLocalstackSession() *sessionv1.Session {
	return sessionv1.Must(sessionv1.NewSession(&awsv1.Config{
		CredentialsChainVerboseErrors: awsv1.Bool(true),
		Endpoint:                      awsv1.String(endpoint),
		Credentials:                   credentialsv1.NewStaticCredentials("localstack", "localstack", ""),
		Region:                        awsv1.String(region),
	}))
}

func NewLocalstackConfig(ctx context.Context) aws.Config {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("localstack", "localstack", "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				PartitionID:   "aws",
				SigningRegion: region,
			}, nil
		})),
	)
	if err != nil {
		panic(err)
	}

	return cfg
}

func NewTestTableName(t *testing.T) string {
	notValidChar := regexp.MustCompile("[^A-Za-z0-9-]+")
	return strings.ToLower(notValidChar.ReplaceAllString(fmt.Sprintf("%s-%s", t.Name(), time.Now().Format("20060102150405.000")), "-"))
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

func NewClientV2(t *testing.T) (aws.Config, *dynamodb.Client, string) {
	cfg := NewLocalstackConfig(context.TODO())
	client := dynamodb.NewFromConfig(cfg)
	tableName := NewTestTableName(t)

	err := appdynamodb.CreateTableIfNecessary(tableName, dynamodbv1.New(NewLocalstackSession()), true)
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		_, _ = client.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{TableName: &tableName})
	})

	return cfg, client, tableName
}
