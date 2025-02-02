package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/api/lambdas/common"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	catalogarchivesync "github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogarchiveasync"
)

func Handler(request events.SQSEvent) error {
	for _, record := range request.Records {
		payload := &catalogarchivesync.RelocateMediasDTO{}
		err := json.Unmarshal([]byte(record.Body), payload)
		if err != nil {
			log.WithError(err).Errorf("failed to parse SQS message '%s': %s", record.Body, err.Error())
			return nil
		}

		err = archive.Relocate(payload.Owner, payload.Ids, payload.FolderName)

		if err == nil {
			log.WithField("Owner", payload.Owner).Infof("Relocated %d medias to %s", len(payload.Ids), payload.FolderName)
		} else {
			log.WithError(err).WithField("Owner", payload.Owner).Errorf("Failed to relocate %d medias to %s", len(payload.Ids), payload.FolderName)
		}

		return errors.Wrapf(err, "failed to relocate %d images to %s%s", len(payload.Ids), payload.Owner, payload.FolderName)
	}
	return nil
}

func main() {
	common.BootstrapArchiveDomain()

	lambda.Start(Handler)
}
