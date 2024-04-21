package dynamotestutils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamodbv1 "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"regexp"
	"strings"
	"testing"
	"time"
)

type DynamodbTestContext struct {
	T      *testing.T
	Ctx    context.Context
	Cfg    aws.Config
	Client *dynamodb.Client
	Table  string
}

func NewTestContext(ctx context.Context, t *testing.T) *DynamodbTestContext {
	cfg := localstackConfig(ctx)
	client := dynamodb.NewFromConfig(cfg)
	tableName := NewTestTableName(t)

	err := appdynamodb.CreateTableIfNecessary(tableName, dynamodbv1.New(NewLocalstackSession()), true)
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		_, _ = client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: &tableName})
	})

	return &DynamodbTestContext{
		T:      t,
		Ctx:    ctx,
		Cfg:    cfg,
		Client: client,
		Table:  tableName,
	}
}

func NewTestTableName(t *testing.T) string {
	notValidChar := regexp.MustCompile("[^A-Za-z0-9-]+")
	return strings.ToLower(notValidChar.ReplaceAllString(fmt.Sprintf("%s-%s", t.Name(), time.Now().Format("20060102150405.000")), "-"))
}

func (d *DynamodbTestContext) Subtest(t *testing.T) *DynamodbTestContext {
	return &DynamodbTestContext{
		T:      t,
		Ctx:    d.Ctx,
		Cfg:    d.Cfg,
		Client: d.Client,
		Table:  d.Table,
	}
}

func localstackConfig(ctx context.Context) aws.Config {
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
