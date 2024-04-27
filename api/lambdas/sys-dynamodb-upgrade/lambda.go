package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/awsfactory"
)

func Handler() error {
	ctx := context.Background()
	factory, err := awsfactory.NewAWSFactory(ctx, awsfactory.DefaultConfigFactory)
	if err != nil {
		return err
	}

	table := viper.GetString(common.DynamoDBTableName)
	if table == "" {
		return errors.Errorf("'%s' environment variable is required with the name of the table.", common.DynamoDBTableName)
	}
	return appdynamodb.CreateTableIfNecessary(ctx, table, factory.GetDynamoDBClient(), false)
}

func main() {
	lambda.Start(Handler)
}
