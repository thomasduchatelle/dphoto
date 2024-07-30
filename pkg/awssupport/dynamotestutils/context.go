package dynamotestutils

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/awsfactory"
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
	factory, err := awsfactory.LocalstackAWSFactory(ctx, awsfactory.LocalstackEndpoint)
	client := factory.GetDynamoDBClient()
	tableName := NewTestTableName(t)

	err = appdynamodb.CreateTableIfNecessary(ctx, tableName, client, true)
	if err != nil {
		panic(err)
	}

	t.Cleanup(func() {
		_, _ = client.DeleteTable(ctx, &dynamodb.DeleteTableInput{TableName: &tableName})
	})

	return &DynamodbTestContext{
		T:      t,
		Ctx:    ctx,
		Cfg:    factory.GetCfg(),
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
