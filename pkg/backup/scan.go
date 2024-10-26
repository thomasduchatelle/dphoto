package backup

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"regexp"
	"strings"
	"time"
)

// Scan a source to discover albums based on original folder structure. Use listeners will be notified on the progress of the scan.
func Scan(owner string, volume SourceVolume, optionSlice ...Options) ([]*ScannedFolder, error) {
	unsafeChar := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	scanId := fmt.Sprintf("%s_%s", strings.Trim(unsafeChar.ReplaceAllString(volume.String(), "_"), "_"), time.Now().Format("20060102_150405"))
	mdc := log.WithFields(log.Fields{
		"ScanId": scanId,
		"Volume": volume.String(),
	})

	options := ReduceOptions(append(optionSlice, Options{ConcurrentUploader: 1})...)

	publisher, hintSize, err := newPublisher(volume)
	if err != nil {
		return nil, err
	}

	referencer, err := NewReferencer(ownermodel.Owner(owner), true)
	if err != nil {
		return nil, err
	}

	receiver := newScanReceiver(mdc, volume)
	run := runner{
		MDC:               mdc,
		Options:           options,
		Publisher:         publisher,
		Analyser:          getDefaultAnalyser(),
		CatalogReferencer: referencer,
		UniqueFilter:      newUniqueFilter(),
		Uploader: RunnerUploaderFunc(func(buffer []*BackingUpMediaRequest, progressChannel chan *ProgressEvent) error {
			// nothing.
			return nil
		}),
		Listeners: []interface{}{receiver},
	}

	progressChannel, _ := run.start(context.TODO(), hintSize)
	tracker := NewTracker(progressChannel, options.Listener)

	err = run.waitToFinish()
	tracker.WaitToComplete()

	return receiver.collect(), err
}
