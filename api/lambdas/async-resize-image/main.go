package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/ephoto/api/lambdas/common"
	"github.com/thomasduchatelle/ephoto/pkg/archive"
	"github.com/thomasduchatelle/ephoto/pkg/archiveadapters/asyncjobadapter"
	"time"
)

var (
	asyncJobAdapter archive.AsyncJobAdapter
)

func Handler(ctx context.Context, request events.SNSEvent) {
	record := request.Records[0]

	var mess []asyncjobadapter.ImageToResizeMessageV1
	err := json.Unmarshal([]byte(record.SNS.Message), &mess)
	if err != nil {
		log.WithError(err).Errorf("failed to parse SQS message '%s': %s", record.SNS.Message, err.Error())
		return
	}

	images := make([]*archive.ImageToResize, len(mess), len(mess))
	for i, img := range mess {
		images[i] = &archive.ImageToResize{
			Owner:    img.Owner,
			MediaId:  img.MediaId,
			StoreKey: img.StoreKey,
			Widths:   img.Widths,
		}
	}

	cachingContext := ctx
	if deadline, ok := ctx.Deadline(); ok {
		var cancel context.CancelFunc
		newDeadline := deadline.Add(-30 * time.Second)
		cachingContext, cancel = context.WithDeadline(cachingContext, newDeadline)
		defer cancel()

		log.Infof("resize %d images with timebox %ds", len(images), int(newDeadline.Sub(time.Now()).Seconds()))
	}

	processed, err := archive.LoadImagesInCache(cachingContext, images...)
	if err != nil {
		log.WithError(err).Errorf("Failed to load image in caches: %s", err.Error())
		return
	}
	if processed < len(images) {
		err = asyncJobAdapter.LoadImagesInCache(images[processed:]...)
		if err != nil {
			log.WithError(err).Errorf("Failed to send %d images back to the queue", len(images)-processed)
		}
	}

	return
}

func main() {
	asyncJobAdapter = common.BootstrapArchiveDomain()

	lambda.Start(Handler)
}
