package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/cfn"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/pkgfactory"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Handler(event *cfn.Event) error {
	ctx := context.Background()

	table := viper.GetString(common.DynamoDBTableName)
	if table == "" {
		return handleError(event, "DBStructureNoID", errors.Errorf("'%s' environment variable is required with the name of the table.", common.DynamoDBTableName))
	}
	physicalResourceID := fmt.Sprintf("DBStructure%s", strings.ReplaceAll(cases.Title(language.English).String(table), "-", ""))

	if event.RequestType == cfn.RequestCreate || event.RequestType == cfn.RequestUpdate {
		log.WithField("PhysicalResourceID", physicalResourceID).
			WithField("event.PhysicalResourceID", event.PhysicalResourceID).
			Infof("Handling %s request.", event.RequestType)

		err := appdynamodb.CreateTableIfNecessary(ctx, table, pkgfactory.AWSFactory(ctx).GetDynamoDBClient(), false)
		if err != nil {
			return handleError(event, physicalResourceID, errors.Wrapf(err, "table structure update failed"))
		}

		return handleSuccess(event, physicalResourceID)

	} else {
		log.WithField("event.PhysicalResourceID", event.PhysicalResourceID).
			Infof("%s RequestType ignored", event.RequestType)
		return handleSuccess(event, physicalResourceID)
	}
}

func main() {
	lambda.Start(Handler)
}

func handleSuccess(event *cfn.Event, physicalResourceID string) error {
	response := cfn.NewResponse(event)
	response.PhysicalResourceID = physicalResourceID
	response.Status = cfn.StatusSuccess
	response.Reason = "Update completed."

	return response.Send()
}

func handleError(event *cfn.Event, physicalResourceID string, err error) error {
	response := cfn.NewResponse(event)
	response.PhysicalResourceID = physicalResourceID
	response.Status = cfn.StatusFailed
	response.Reason = err.Error()

	sendErr := response.Send()
	if sendErr != nil {
		log.WithError(err).Error("Failed to ack the Cloud Formation event")
		return errors.Wrapf(err, "failed to send response dur to [%s], migration original failed because", sendErr.Error())
	}

	return err
}
