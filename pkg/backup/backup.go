// Package backup is providing commands to inspect a file system (hard-drive, USB, Android, S3) and backup medias to a remote DPhoto storage.
package backup

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/backup/chain"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"slices"
)

type SourceVolume interface {
	String() string
	FindMedias(ctx context.Context) ([]FoundMedia, error)
}

type BatchBackup struct {
	CataloguerFactory CataloguerFactory
	DetailsReaders    []DetailsReaderAdapter
	InsertMediaPort   InsertMediaPort
	ArchivePort       BArchiveAdapter
}

// Backup is analysing each media and is backing it up if not already in the catalog.
func (b *BatchBackup) Backup(ctx context.Context, owner ownermodel.Owner, volume SourceVolume, optionsSlice ...Options) (CompletionReport, error) {
	launcher, tracker, err := b.prepareVolumeBackup(ctx, ReduceOptions(optionsSlice...), volume.String(), owner)
	if err != nil {
		return nil, err
	}

	err, _ = <-launcher.Process(ctx, volume)

	return tracker, err
}

func (b *BatchBackup) prepareVolumeBackup(ctx context.Context, options Options, volumeName string, owner ownermodel.Owner) (analyserLauncher, *trackerV2, error) {
	tracker, report := newTrackerV2(options) // TODO is using the tracker to collect the report the best way to do it ?
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
			PostAnalyserFilterRejects: []RejectedMediaObserver{scanLogger, tracker /*, reportBuilder*/},
			PostAnalyserRejects:       []RejectedMediaObserver{scanLogger, tracker},
			PreCataloguerFilter:       []CatalogReferencerObserver{scanLogger},
			PostCatalogFiltersOut:     []CataloguerFilterObserver{scanLogger, tracker},
			Wrappers:                  []chain.CloserFunc{tracker.NoMoreEvents},
			// PostCatalogFiltersIn is not supported.
		},
		Uploader: &uploader{
			Owner:            owner,
			InsertMediaPort:  b.InsertMediaPort,
			ArchivePort:      b.ArchivePort,
			UploaderObserver: tracker,
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

	//uploaderLauncher := &multithreadedUploaderLauncher{
	//	uploader: &uploader{
	//		Owner:            owner,
	//		InsertMediaPort:  b.InsertMediaPort,
	//		ArchivePort:      b.ArchivePort,
	//		UploaderObserver: tracker,
	//	},
	//	channel:                    uploaderChannelInstance,
	//	done:                       make(chan interface{}),
	//	concurrentUploaderRoutines: options.ConcurrencyParameters.NumberOfConcurrentUploaderRoutines(),
	//}
	launcher := scanAndBackupCommonLauncher(&config.scanConfiguration, options,
		&reBufferLink[BackingUpMediaRequest]{ // TODO rebuffer is not tested...
			bufferLink: bufferLink[BackingUpMediaRequest]{
				Buffer: buffer[BackingUpMediaRequest]{
					content: make([]BackingUpMediaRequest, 0, defaultValue(options.BatchSize, 1)),
				},
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

		//&chain.MultithreadedLink[[]BackingUpMediaRequest, []BackingUpMediaRequest]{
		//	NumberOfRoutines: 1,
		//	ConsumerBuilder:  chain.PassThrough[[]BackingUpMediaRequest](),
		//	Next: &chain.CloseWrapperLink[[]BackingUpMediaRequest]{
		//		CloserFuncs: slices.Concat(config.Wrappers, []chain.CloserFunc{chain.CloserFunc(cancelFunc)}),
		//		Next:        chain.EndOfTheChain[[]BackingUpMediaRequest](finalizer(config.PostCatalogFiltersIn)...),
		//	},
		//},
	)

	err := launcher.Starts(ctx, chain.NewErrorCollector(func(err error) {
		cancelFunc()
	}))
	return launcher, err
}

type reBufferLink[Consumed any] struct {
	bufferLink[Consumed]
}

func (l *reBufferLink[Consumed]) Consume(ctx context.Context, buf []Consumed) error {
	for _, item := range buf {
		err := l.bufferLink.Consume(ctx, item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BatchBackup) _prepareVolumeBackup(ctx context.Context, options Options, volumeName string, owner ownermodel.Owner) (analyserLauncher, *trackerV2, error) {
	tracker, report := newTrackerV2(options) // TODO is using the tracker to collect the report the best way to do it ?
	//reportBuilder := newScanReportBuilder()
	scanLogger := newLogger(volumeName)

	//uploaderChannelInstance := make(uploaderChannel)

	monitoring := &scanListeners{
		scanCompleteObserver:      tracker,
		PostAnalyserSuccess:       []AnalysedMediaObserver{scanLogger},
		PostAnalyserRejects:       []RejectedMediaObserver{scanLogger, tracker},
		PostAnalyserFilterRejects: []RejectedMediaObserver{scanLogger, tracker /*, reportBuilder*/},
		PreCataloguerFilter:       []CatalogReferencerObserver{scanLogger},
		PostCatalogFiltersIn: []CatalogReferencerObserver{
			tracker, /*, reportBuilder*/
			//uploaderChannelInstance,
		},
		PostCatalogFiltersOut: []CataloguerFilterObserver{scanLogger, tracker},
	}
	if options.SkipRejects {
		monitoring.PostAnalyserRejects = append(monitoring.PostAnalyserRejects /*, reportBuilder*/)
	}

	uploaderLauncher := &multithreadedUploaderLauncher{
		uploader: &uploader{
			Owner:            owner,
			InsertMediaPort:  b.InsertMediaPort,
			ArchivePort:      b.ArchivePort,
			UploaderObserver: tracker,
		},
		//channel:                    uploaderChannelInstance,
		done:                       make(chan interface{}),
		concurrentUploaderRoutines: options.ConcurrencyParameters.NumberOfConcurrentUploaderRoutines(),
	}

	controller := newMultiThreadedController(options.ConcurrencyParameters, monitoring)
	controller.registerWrappers(uploaderLauncher)
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
	return uploaderLauncher.wrapLauncher(launcher), report, err
}

type multithreadedUploaderLauncher struct {
	launcher                   analyserLauncher
	uploader                   CatalogReferencerObserver
	channel                    chan []BackingUpMediaRequest
	done                       chan interface{}
	concurrentUploaderRoutines int
}

func (m *multithreadedUploaderLauncher) Close() error {
	close(m.channel)
	<-m.done
	return nil
}

func (m *multithreadedUploaderLauncher) Process(ctx context.Context, volume SourceVolume) chan error {
	startsInParallel(ctx, m.concurrentUploaderRoutines, func(ctx context.Context) {
		for requests := range m.channel {
			err := m.uploader.OnMediaCatalogued(ctx, requests)
			if err != nil {
				// TODO handle the error to abort the process
			}
		}
	}, func() {
		close(m.done)
	})

	return m.launcher.Process(ctx, volume)
}

func (m *multithreadedUploaderLauncher) wrapLauncher(launcher scanningLauncher) analyserLauncher {
	m.launcher = launcher
	return m
}
