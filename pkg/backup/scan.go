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

	err, _ = <-launcher.process(ctx, volume)

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

	launcher := chain.SliceLauncher[FoundMedia]{
		Producer: volume.FindMedias,
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
			Next: &chain.MultithreadedLink[*AnalysedMedia, []*AnalysedMedia]{
				NumberOfRoutines: 1,
				ConsumerBuilder: func(c chain.Consumer[[]*AnalysedMedia]) chain.Consumer[*AnalysedMedia] {
					// TODO buffer needs to be closed
					return chain.ConsumerFunc[*AnalysedMedia](newAnalysedMediaBufferAdapter(options, analysedMediasBatchObserverFunc(c.Consume)).OnAnalysedMedia)
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

	//controller := scanMultithreadedOrchestration{
	//	volumeProcess:     volumeAdapter{
	//		observer: nil,
	//	},
	//	analyserProcess:   nil,
	//	bufferProcess:     nil,
	//	cataloguerProcess: nil,
	//}

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

	launcher, err := newScanningChain(ctx, controller, scanningOptions{
		Options:    options,
		cataloguer: cataloguer,
		analyser:   options.GetAnalyserDecorator().Decorate(newDefaultAnalyser(s.DetailsReaders...)),
	})
	return launcher, reportBuilder, err
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
