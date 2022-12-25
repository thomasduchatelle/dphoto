package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/ephoto/api/lambdas/common"
	"github.com/thomasduchatelle/ephoto/pkg/archive"
	"github.com/thomasduchatelle/ephoto/pkg/archiveadapters/asyncjobadapter"
)

func Handler(request events.SQSEvent) error {
	for _, record := range request.Records {
		mess := &asyncjobadapter.WarmUpCacheByFolderMessageV1{}
		err := json.Unmarshal([]byte(record.Body), mess)
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
		}).Infof("Re-populating cache of media %s at width %d", mess.MissedStoreKey, mess.Width)
		err = archive.WarmUpCacheByFolder(mess.Owner, mess.MissedStoreKey, mess.Width)
		if err != nil {
			log.WithError(err).WithField("Owner", mess.Owner).Errorf("Failed to list cache to reload for storeKey=%s : %s", mess.MissedStoreKey, err.Error())
		}
	}
	return nil
}

func main() {
	common.BootstrapArchiveDomain()

	lambda.Start(Handler)
}
