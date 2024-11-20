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
	tracker, _ := newTrackerV2(options)
	report := newBackupReportBuilder()
	scanLogger := newLogger(volumeName)

	cataloguer, err := b.newCataloguer(ctx, owner)
	if err != nil {
		return nil, nil, err
	}

	config := &backupConfiguration{
		scanConfiguration: scanConfiguration{
			Analyser:                 options.GetAnalyserDecorator().Decorate(newDefaultAnalyser(b.DetailsReaders...), tracker),
			Cataloguer:               cataloguer,
			ScanCompleteObserver:     tracker,
			PostAnalyserRejects:      []RejectedMediaObserver{scanLogger, tracker, report},
			PostCatalogFiltersIn:     []CatalogReferencerObserver{scanLogger, tracker},
			PostCataloguerFiltersOut: []CataloguerFilterObserver{scanLogger, tracker, report},
			Wrappers:                 []chain.CloserFunc{tracker.NoMoreEvents},
		},
		Uploader: &uploader{
			Owner:             owner,
			InsertMediaPort:   b.InsertMediaPort,
			ArchivePort:       b.ArchivePort,
			UploaderObservers: []uploaderObserver{tracker, report},
		},
	}
	if !options.SkipRejects && options.RejectDir == "" {
		config.PostAnalyserRejects = append(config.PostAnalyserRejects, new(analyserFailsFastObserver))
	}
	if options.RejectDir != "" {
		observer, err := newCopyRejectsObserver(options.RejectDir)
		if err != nil {
			return nil, nil, err
		}
		config.PostAnalyserRejects = append(config.PostAnalyserRejects, observer)
	}

	launcher, err := multithreadedBackupRuntime(ctx, options, config)
	return launcher, report, err
}

func (b *BatchBackup) newCataloguer(ctx context.Context, owner ownermodel.Owner) (Cataloguer, error) {
	referencer, err := b.CataloguerFactory.NewOwnerScopedCataloguer(ctx, owner)
	return referencer, errors.Wrapf(err, "failed to create a cataloguer for %s", owner)
}

type backupConfiguration struct {
	scanConfiguration

	Uploader *uploader
}

func multithreadedBackupRuntime(ctxNonCancelable context.Context, options Options, config *backupConfiguration) (analyserLauncher, error) {
	ctx, cancelFunc := context.WithCancel(ctxNonCancelable)

	launcher := scanAndBackupCommonLauncher(&config.scanConfiguration, options,
		&chain.ReBufferLink[BackingUpMediaRequest]{
			BufferLink: chain.BufferLink[BackingUpMediaRequest]{
				BufferCapacity: options.BatchSize,
				Next: &chain.MultithreadedLink[[]BackingUpMediaRequest, []BackingUpMediaRequest]{
					NumberOfRoutines: options.ConcurrencyParameters.NumberOfConcurrentUploaderRoutines(),
					ConsumerBuilder: func(consumer chain.Consumer[[]BackingUpMediaRequest]) chain.Consumer[[]BackingUpMediaRequest] {
						observers := CatalogReferencerObservers(slices.Concat(
							config.PostCatalogFiltersIn,
							CatalogReferencerObservers{config.Uploader, CatalogReferencerObserverFunc(consumer.Consume)},
						))

						return chain.ConsumerFunc[[]BackingUpMediaRequest](func(ctx context.Context, consumed []BackingUpMediaRequest) error {
							return observers.OnMediaCatalogued(ctx, consumed)
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
