package backup

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/backup/chain"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

type BatchScanner struct {
	CataloguerFactory CataloguerFactory
	DetailsReaders    []DetailsReaderAdapter
}

func (s *BatchScanner) Scan(ctx context.Context, owner ownermodel.Owner, volume SourceVolume, optionSlice ...Options) ([]*ScannedFolder, error) {

	options := ReduceOptions(optionSlice...)

	launcher, reportBuilder, err := s.prepareVolumeScan(ctx, options, volume.String(), owner, volume)
	if err != nil {
		return nil, err
	}

	err = <-launcher.Process(ctx, volume)

	return reportBuilder.build(), err
}

//type scanMultithreadedOrchestration struct {
//	volumeProcess     *volumeAdapter
//	analyserProcess   foundMediaObserver
//	bufferProcess     AnalysedMediaObserver // TODO is also closable
//	cataloguerProcess analysedMediasBatchObserver
//}

func (s *BatchScanner) prepareVolumeScan(ctx context.Context, options Options, volumeName string, owner ownermodel.Owner, volume SourceVolume) (analyserLauncher, *scanReportBuilder, error) {
	tracker, _ := newTrackerV2(options)
	reportBuilder := newScanReportBuilder()
	scanLogger := newLogger(volumeName)

	cataloguer, err := s.CataloguerFactory.NewDryRunCataloguer(ctx, owner)
	if err != nil {
		return nil, nil, err
	}

	launcher := &chain.SingleLauncher[SourceVolume, FoundMedia]{
		Function: func(ctx context.Context, consumed SourceVolume) ([]FoundMedia, error) {
			medias, err := consumed.FindMedias(ctx)
			if err != nil {
				return nil, err
			}

			return medias, tracker.OnScanComplete(ctx, len(medias), sizeOfAllMedias(medias))
		},
		Next: &chain.MultithreadedLink[FoundMedia, *AnalysedMedia]{
			NumberOfRoutines: options.ConcurrencyParameters.NumberOfConcurrentAnalyserRoutines(),
			ConsumerBuilder: func(consumer chain.Consumer[*AnalysedMedia]) chain.Consumer[FoundMedia] {
				analyser := &analyserAdapter{
					analyser:     options.GetAnalyserDecorator().Decorate(newDefaultAnalyser(s.DetailsReaders...)),
					analysed:     []AnalysedMediaObserver{AnalysedMediaObserverFunc(consumer.Consume)},
					beforeFilter: []AnalysedMediaObserver{scanLogger},
					filteredOut:  []RejectedMediaObserver{scanLogger, tracker, reportBuilder},
					rejected:     []RejectedMediaObserver{scanLogger, tracker},
				}
				return chain.ConsumerFunc[FoundMedia](analyser.OnFoundMedia)
			},
			Next: &BufferLink[*AnalysedMedia]{
				Buffer: buffer[*AnalysedMedia]{
					content: make([]*AnalysedMedia, 0, defaultValue(options.BatchSize, 1)),
					// note - consumer is set during "Starts" call
				},
				Next: &chain.MultithreadedLink[[]*AnalysedMedia, []BackingUpMediaRequest]{
					NumberOfRoutines: options.ConcurrencyParameters.NumberOfConcurrentCataloguerRoutines(),
					ConsumerBuilder: func(consumer chain.Consumer[[]BackingUpMediaRequest]) chain.Consumer[[]*AnalysedMedia] {
						adapter := &cataloguerAdapter{
							cataloguer:  cataloguer,
							options:     options,
							preFilters:  []CatalogReferencerObserver{scanLogger},
							catalogued:  []CatalogReferencerObserver{tracker, CatalogReferencerObserverFunc(consumer.Consume)},
							filteredOut: []CataloguerFilterObserver{scanLogger, tracker},
						}
						return chain.ConsumerFunc[[]*AnalysedMedia](adapter.OnBatchOfAnalysedMedia)
					},
					Next: &chain.MultithreadedLink[[]BackingUpMediaRequest, []BackingUpMediaRequest]{
						NumberOfRoutines: 0,
						ConsumerBuilder: func(consumer chain.Consumer[[]BackingUpMediaRequest]) chain.Consumer[[]BackingUpMediaRequest] {
							return chain.ConsumerFunc[[]BackingUpMediaRequest](func(ctx context.Context, requests []BackingUpMediaRequest) error {
								err := reportBuilder.OnMediaCatalogued(ctx, requests)
								if err != nil {
									return err
								}

								return consumer.Consume(ctx, requests)
							})
						},
						Next: &chain.EnderChainLink[[]BackingUpMediaRequest]{
							Operator: func(ctx context.Context, requests []BackingUpMediaRequest) error {
								return nil
							},
						},
					},
				},
			},
		},
	}

	err = launcher.Starts(ctx, chain.NewErrorCollector())
	if err != nil {
		return nil, nil, err
	}

	return launcher, reportBuilder, nil
}

func (s *BatchScanner) _ref_prepareVolumeScan(ctx context.Context, options Options, volumeName string, owner ownermodel.Owner) (analyserLauncher, *scanReportBuilder, error) {
	tracker, _ := newTrackerV2(options)
	reportBuilder := newScanReportBuilder()
	scanLogger := newLogger(volumeName)

	monitoring := &scanListeners{
		scanCompleteObserver:      tracker,
		PostAnalyserSuccess:       []AnalysedMediaObserver{scanLogger},
		PostAnalyserRejects:       []RejectedMediaObserver{scanLogger, tracker},
		PostAnalyserFilterRejects: []RejectedMediaObserver{scanLogger, tracker, reportBuilder},
		PreCataloguerFilter:       []CatalogReferencerObserver{scanLogger},
		PostCatalogFiltersIn:      []CatalogReferencerObserver{tracker, reportBuilder},
		PostCatalogFiltersOut:     []CataloguerFilterObserver{scanLogger, tracker},
	}
	if options.SkipRejects {
		monitoring.PostAnalyserRejects = append(monitoring.PostAnalyserRejects, reportBuilder)
	}

	controller := newMultiThreadedController(options.ConcurrencyParameters, monitoring)
	controller.registerWrappers(tracker)

	cataloguer, err := s.CataloguerFactory.NewDryRunCataloguer(ctx, owner)
	if err != nil {
		return nil, nil, err
	}

	launcher, err := newScanningChain(ctx, controller, scanningOptions{
		Options:    options,
		cataloguer: cataloguer,
		analyser:   options.GetAnalyserDecorator().Decorate(newDefaultAnalyser(s.DetailsReaders...)),
	})
	return launcher, reportBuilder, err
}

type BufferLink[Consumed any] struct {
	Next    chain.Link[[]Consumed]
	channel chan Consumed
	Buffer  buffer[Consumed]
}

func (l *BufferLink[Consumed]) Consume(ctx context.Context, consumed Consumed) error {
	l.channel <- consumed
	return nil
}

func (l *BufferLink[Consumed]) Starts(ctx context.Context, collector chain.ChainableErrorCollector) error {
	l.channel = make(chan Consumed, 255)
	l.Buffer.consumer = l.Next.Consume

	go func() {
		defer l.Next.NotifyUpstreamCompleted()

		for {
			select {
			case consumed, more := <-l.channel:
				if more {
					err := l.Buffer.Append(ctx, consumed)
					if err != nil {
						collector.OnError(err)
					}
				} else {
					err := l.Buffer.Flush(ctx)
					if err != nil {
						collector.OnError(err)
					}
					return
				}
			}
		}
	}()

	return l.Next.Starts(ctx, collector)
}

func (l *BufferLink[Consumed]) WaitForCompletion() chan error {
	return l.Next.WaitForCompletion()
}

func (l *BufferLink[Consumed]) NotifyUpstreamCompleted() {
	close(l.channel)
}
