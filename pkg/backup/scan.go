package backup

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

var (
	datePrefix = regexp.MustCompile("^[0-9]{4}-[01Q][0-9][-_]")
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

	publisher, hintSize, err := newPublisher(volume)

	receiver := newScanReceiver(mdc, volume)

	run := runner{
		MDC:                  mdc,
		Publisher:            publisher,
		Analyser:             newBackupAnalyseMedia(),
		Cataloger:            newScannerCataloger(owner),
		UniqueFilter:         newUniqueFilter(),
		Uploader:             receiver.receive,
		ConcurrentAnalyser:   ConcurrentAnalyser,
		ConcurrentCataloguer: ConcurrentCataloguer,
		ConcurrentUploader:   1,
		BatchSize:            BatchSize,
	}

	progressChannel, _ := run.start(context.TODO(), hintSize)
	tracker := NewTracker(progressChannel, options.Listener)

	err = run.waitToFinish()
	tracker.WaitToComplete()

	return receiver.collect(), receiver.rejects, err
}
