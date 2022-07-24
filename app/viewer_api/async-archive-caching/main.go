package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/app/viewer_api/common"
	"github.com/thomasduchatelle/dphoto/domain/archiveadapters/jobqueuesns"
)

func Handler(request events.SQSEvent) error {
	for _, record := range request.Records {
		mess, err := unpackMessage(record)
		if err != nil {
			log.WithError(err).Errorf("failed to parse SQS message '%s': %s", record.Body, err.Error())
			return nil
		}

		if mess.Owner == "" || mess.Width == 0 {
			log.Errorf("Invalid message, unmarshaling have failed with nessage '%s'", record.Body)
			return nil
		}

		log.WithFields(log.Fields{
			"Owner": mess.Owner,
			"Width": mess.Width,
		}).Infof("Re-populating cache of media %s at width %d", mess.StoreKey, mess.Width)

	}
	return nil
}

func unpackMessage(record events.SQSMessage) (*jobqueuesns.MisfiredCacheMessageV1, error) {
	snsEntity := &events.SNSEntity{}
	err := json.Unmarshal([]byte(record.Body), snsEntity)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshaling events.SNSEntity '%s'", record.Body)
	}

	mess := &jobqueuesns.MisfiredCacheMessageV1{}
	err = json.Unmarshal([]byte(snsEntity.Message), mess)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshaling jobqueuesns.MisfiredCacheMessageV1 '%s'", snsEntity.Message)
	}

	return mess, nil
}

func main() {
	common.BootstrapCatalogAndArchiveDomains()

	lambda.Start(Handler)
}
