// Package backup is providing commands to inspect a file system (hard-drive, USB, Android, S3) and backup medias to a remote DPhoto storage.
package backup

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"regexp"
	"strings"
	"time"
)

type SourceVolume interface {
	String() string
	FindMedias() ([]FoundMedia, error)
}

type BatchBackup struct {
	CataloguerFactory CataloguerFactory
	DetailsReaders    []DetailsReaderAdapter
	InsertMediaPort   InsertMediaPort
	ArchivePort       BArchiveAdapter
}

// Backup is analysing each media and is backing it up if not already in the catalog.
func (b *BatchBackup) Backup(ctx context.Context, owner ownermodel.Owner, volume SourceVolume, optionsSlice ...Options) (CompletionReport, error) {
	launcher, tracker, err := b.prepareVolumeScan(ctx, ReduceOptions(optionsSlice...), volume.String(), owner)
	if err != nil {
		return nil, err
	}

	err, _ = <-launcher.process(ctx, volume)

	return tracker, err
}

func (b *BatchBackup) Backup2(ctx context.Context, owner ownermodel.Owner, volume SourceVolume, optionsSlice ...Options) (CompletionReport, error) {
	unsafeChar := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	backupId := fmt.Sprintf("%s_%s", strings.Trim(unsafeChar.ReplaceAllString(volume.String(), "_"), "_"), time.Now().Format("20060102_150405"))
	mdc := log.WithFields(log.Fields{
		"BackupId": backupId,
		"Volume":   volume.String(),
	})

	options := ReduceOptions(optionsSlice...)

	referencer, err := b.newCataloguer(ctx, owner, options.DryRun)
	if err != nil {
		return nil, err
	}

	publisher, hintSize, err := newPublisher(volume)

	run := runner{
		MDC:               mdc,
		Options:           options,
		Publisher:         publisher,
		Analyser:          options.GetAnalyserDecorator().Decorate(newDefaultAnalyser(b.DetailsReaders...)),
		CatalogReferencer: referencer,
		UniqueFilter:      newUniqueFilter(),
		Uploader:          &uploader{Owner: owner, InsertMediaPort: b.InsertMediaPort, ArchivePort: b.ArchivePort},
	}

	progressChannel, _ := run.start(ctx, hintSize)
	backupReport := NewTracker(progressChannel, options.Listener)

	err = run.waitToFinish()
	backupReport.WaitToComplete()

	if err == nil {
		mdc.Infof("Backup completed, %d medias backed up.", backupReport.MediaCount())
	} else {
		mdc.WithError(err).Errorf("Backup faifed with err: %s", err.Error())
	}
	return backupReport, err
}

func (b *BatchBackup) newCataloguer(ctx context.Context, owner ownermodel.Owner, dryRun bool) (Cataloguer, error) {
	var referencer Cataloguer
	var err error

	if dryRun {
		referencer, err = b.CataloguerFactory.NewDryRunCataloguer(ctx, owner)
	} else {
		referencer, err = b.CataloguerFactory.NewAlbumCreatorCataloguer(ctx, owner)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create a cataloguer for %s with dryRun=%t", owner, dryRun)
	}

	return referencer, nil
}

func (b *BatchBackup) prepareVolumeScan(ctx context.Context, options Options, volumeName string, owner ownermodel.Owner) (analyserLauncher, *trackerV2, error) {
	tracker, report := newTrackerV2(options) // TODO is using the tracker to collect the report the best way to do it ?
	//reportBuilder := newScanReportBuilder()
	scanLogger := newLogger(volumeName)

	monitoring := &scanListeners{
		scanCompleteObserver:      tracker,
		PostAnalyserSuccess:       []AnalysedMediaObserver{scanLogger},
		PostAnalyserRejects:       []RejectedMediaObserver{scanLogger, tracker},
		PostAnalyserFilterRejects: []RejectedMediaObserver{scanLogger, tracker /*, reportBuilder*/},
		PreCataloguerFilter:       []CatalogReferencerObserver{scanLogger},
		PostCatalogFiltersIn: []CatalogReferencerObserver{
			tracker, /*, reportBuilder*/
			&uploader{
				Owner:            owner,
				InsertMediaPort:  b.InsertMediaPort,
				ArchivePort:      b.ArchivePort,
				UploaderObserver: tracker,
			},
		},
		PostCatalogFiltersOut: []CataloguerFilterObserver{scanLogger, tracker},
	}
	if options.SkipRejects {
		monitoring.PostAnalyserRejects = append(monitoring.PostAnalyserRejects /*, reportBuilder*/)
	}

	controller := newMultiThreadedController(options.ConcurrencyParameters, monitoring)
	controller.registerWrappers(tracker)

	cataloguer, err := b.newCataloguer(ctx, owner, options.DryRun)
	if err != nil {
		return nil, nil, err
	}

	launcher, err := newScanningChain(ctx, controller, scanningOptions{
		Options:    options,
		cataloguer: cataloguer,
		analyser:   options.GetAnalyserDecorator().Decorate(newDefaultAnalyser(b.DetailsReaders...)),
	})
	return launcher, report, err
}
