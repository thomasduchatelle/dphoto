// Package backup is providing commands to inspect a file system (hard-drive, USB, Android, S3) and backup medias to a remote DPhoto storage.
package backup

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/backup/chain"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"slices"
)

type BatchBackup struct {
	CataloguerFactory CataloguerFactory
	DetailsReaders    []DetailsReader
	InsertMediaPort   InsertMediaPort
	ArchivePort       ArchiveMediaPort
}

// Backup is analysing each media and is backing it up if not already in the catalog.
func (b *BatchBackup) Backup(ctx context.Context, owner ownermodel.Owner, volume SourceVolume, optionsSlice ...Options) (Report, error) {
	launcher, report, err := b.prepareVolumeBackup(ctx, ReduceOptions(optionsSlice...), volume.String(), owner)
	if err != nil {
		return nil, err
	}

	err, _ = <-launcher.Process(ctx, volume)

	return report, err
}

func (b *BatchBackup) prepareVolumeBackup(ctx context.Context, options Options, volumeName string, owner ownermodel.Owner) (analyserLauncher, Report, error) {
	tracker, _ := newTrackerV2(options) // TODO is using the tracker to collect the report the best way to do it ?
	report := newBackupReportBuilder()
	scanLogger := newLogger(volumeName)

	cataloguer, err := b.newCataloguer(ctx, owner, options.DryRun)
	if err != nil {
		return nil, nil, err
	}

	config := &backupConfiguration{
		scanConfiguration: scanConfiguration{
			Analyser:                  options.GetAnalyserDecorator().Decorate(newDefaultAnalyser(b.DetailsReaders...)),
			Cataloguer:                cataloguer,
			ScanCompleteObserver:      tracker,
			PostAnalyserSuccess:       []AnalysedMediaObserver{scanLogger},
			PostAnalyserFilterRejects: []RejectedMediaObserver{scanLogger, tracker, report},
			PostAnalyserRejects:       []RejectedMediaObserver{scanLogger, tracker, report},
			PreCataloguerFilter:       []CatalogReferencerObserver{scanLogger},
			PostCatalogFiltersOut:     []CataloguerFilterObserver{scanLogger, tracker, report},
			Wrappers:                  []chain.CloserFunc{tracker.NoMoreEvents},
			// PostCatalogFiltersIn is not supported.
		},
		Uploader: &uploader{
			Owner:             owner,
			InsertMediaPort:   b.InsertMediaPort,
			ArchivePort:       b.ArchivePort,
			UploaderObservers: []uploaderObserver{tracker, report},
		},
	}
	if options.SkipRejects {
		config.PostAnalyserRejects = append(config.PostAnalyserRejects /*, reportBuilder*/)
	} else {
		config.PostAnalyserRejects = append(config.PostAnalyserRejects, new(analyserFailsFastObserver))
	}

	launcher, err := multithreadedBackupRuntime(ctx, options, config)
	return launcher, report, err
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

type backupConfiguration struct {
	scanConfiguration

	Uploader *uploader
}

func multithreadedBackupRuntime(ctxNonCancelable context.Context, options Options, config *backupConfiguration) (analyserLauncher, error) {
	ctx, cancelFunc := context.WithCancel(ctxNonCancelable)

	launcher := scanAndBackupCommonLauncher(&config.scanConfiguration, options,
		&chain.ReBufferLink[BackingUpMediaRequest]{ // TODO ReBufferLink is added a behaviour which is not tested (and if removed it would still work)
			BufferLink: chain.BufferLink[BackingUpMediaRequest]{
				BufferCapacity: options.BatchSize,
				Next: &chain.MultithreadedLink[[]BackingUpMediaRequest, []BackingUpMediaRequest]{
					NumberOfRoutines: options.ConcurrencyParameters.NumberOfConcurrentUploaderRoutines(),
					ConsumerBuilder: func(consumer chain.Consumer[[]BackingUpMediaRequest]) chain.Consumer[[]BackingUpMediaRequest] {
						up := config.Uploader
						return chain.ConsumerFunc[[]BackingUpMediaRequest](func(ctx context.Context, consumed []BackingUpMediaRequest) error {
							err := up.OnMediaCatalogued(ctx, consumed)
							if err != nil {
								return err
							}
							return consumer.Consume(ctx, consumed)
						})
					},
					Next: &chain.CloseWrapperLink[[]BackingUpMediaRequest]{
						CloserFuncs: slices.Concat(config.Wrappers, []chain.CloserFunc{chain.CloserFunc(cancelFunc)}),
						Next:        chain.EndOfTheChain[[]BackingUpMediaRequest](),
					},
				},
			},
		},
	)

	err := launcher.Starts(ctx, chain.NewErrorCollector(func(err error) {
		cancelFunc()
	}))
	return launcher, err
}
