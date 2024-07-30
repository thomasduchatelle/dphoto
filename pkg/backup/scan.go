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
func Scan(owner string, volume SourceVolume, optionSlice ...Options) ([]*ScannedFolder, []FoundMedia, error) {
	unsafeChar := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	scanId := fmt.Sprintf("%s_%s", strings.Trim(unsafeChar.ReplaceAllString(volume.String(), "_"), "_"), time.Now().Format("20060102_150405"))
	mdc := log.WithFields(log.Fields{
		"ScanId": scanId,
		"Volume": volume.String(),
	})

	options := readOptions(optionSlice)
	options.DryRun = true

	publisher, hintSize, err := newPublisher(volume)
	if err != nil {
		return nil, nil, err
	}

	cataloger, err := NewCataloger(ownermodel.Owner(owner), options)
	if err != nil {
		return nil, nil, err
	}

	receiver := newScanReceiver(mdc, volume)
	run := runner{
		MDC:                  mdc,
		Publisher:            publisher,
		Analyser:             options.GetAnalyserDecorator().Decorate(newBackupAnalyseMedia()),
		Cataloger:            cataloger,
		UniqueFilter:         newUniqueFilter(),
		Uploader:             RunnerUploaderFunc(receiver.receive),
		ConcurrentAnalyser:   options.ConcurrentAnalyser,
		ConcurrentCataloguer: options.ConcurrentCataloguer,
		ConcurrentUploader:   1,
		BatchSize:            options.BatchSize,
		SkipRejects:          options.SkipRejects,
	}

	progressChannel, _ := run.start(context.TODO(), hintSize)
	tracker := NewTracker(progressChannel, options.Listener)

	err = run.waitToFinish()
	tracker.WaitToComplete()

	return receiver.collect(), receiver.rejects, err
}
