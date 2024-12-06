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

	launcher, reportBuilder, err := s.prepareVolumeScan(ctx, options, volume.String(), owner)
	if err != nil {
		return nil, err
	}

	err = <-launcher.Process(ctx, volume)

	return reportBuilder.build(), err
}

// scanListeners list the listeners that will be notified during the scan process.
type scanConfiguration struct {
	Analyser                  Analyser
	Cataloguer                Cataloguer
	ScanCompleteObserver      scanCompleteObserver
	PostAnalyserSuccess       []AnalysedMediaObserver
	PostAnalyserFilterRejects []RejectedMediaObserver
	PostAnalyserRejects       []RejectedMediaObserver
	PreCataloguerFilter       []CatalogReferencerObserver
	PostCatalogFiltersIn      []CatalogReferencerObserver
	PostCatalogFiltersOut     []CataloguerFilterObserver
	closer                    func() error
}

func (s *BatchScanner) prepareVolumeScan(ctx context.Context, options Options, volumeName string, owner ownermodel.Owner) (analyserLauncher, *scanReportBuilder, error) {
	tracker, _ := newTrackerV2(options)
	reportBuilder := newScanReportBuilder()
	scanLogger := newLogger(volumeName)

	cataloguer, err := s.CataloguerFactory.NewDryRunCataloguer(ctx, owner)
	if err != nil {
		return nil, nil, err
	}

	config := &scanConfiguration{
		ScanCompleteObserver:      tracker,
		Analyser:                  options.GetAnalyserDecorator().Decorate(newDefaultAnalyser(s.DetailsReaders...)),
		Cataloguer:                cataloguer,
		PostAnalyserSuccess:       []AnalysedMediaObserver{scanLogger},
		PostAnalyserRejects:       []RejectedMediaObserver{scanLogger, tracker},
		PostAnalyserFilterRejects: []RejectedMediaObserver{scanLogger, tracker, reportBuilder},
		PreCataloguerFilter:       []CatalogReferencerObserver{scanLogger},
		PostCatalogFiltersIn:      []CatalogReferencerObserver{tracker, reportBuilder},
		PostCatalogFiltersOut:     []CataloguerFilterObserver{scanLogger, tracker},
		closer:                    tracker.Close,
	}
	if options.SkipRejects {
		config.PostAnalyserRejects = append(config.PostAnalyserRejects, reportBuilder)
	} else {
		config.PostAnalyserRejects = append(config.PostAnalyserRejects, new(analyserFailsFastObserver))
	}

	launcher, err := multithreadedScanRuntime(ctx, options, config)
	return launcher, reportBuilder, err
}

func multithreadedScanRuntime(ctx context.Context, options Options, config *scanConfiguration) (analyserLauncher, error) {
	launcher := &chain.SingleLauncher[SourceVolume, FoundMedia]{
		Function: func(ctx context.Context, volume SourceVolume) ([]FoundMedia, error) {
			medias, err := volume.FindMedias(ctx)
			if err != nil {
				return nil, err
			}

			return medias, config.ScanCompleteObserver.OnScanComplete(ctx, len(medias), sizeOfAllMedias(medias))
		},
		Next: &chain.MultithreadedLink[FoundMedia, *AnalysedMedia]{
			NumberOfRoutines: options.ConcurrencyParameters.NumberOfConcurrentAnalyserRoutines(),
			ConsumerBuilder: func(consumer chain.Consumer[*AnalysedMedia]) chain.Consumer[FoundMedia] {
				analyser := &analyserAdapter{
					analyser:     config.Analyser,
					analysed:     []AnalysedMediaObserver{AnalysedMediaObserverFunc(consumer.Consume)},
					beforeFilter: config.PostAnalyserSuccess,
					filteredOut:  config.PostAnalyserFilterRejects,
					rejected:     config.PostAnalyserRejects,
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
							cataloguer:  config.Cataloguer,
							options:     options,
							preFilters:  config.PreCataloguerFilter,
							catalogued:  []CatalogReferencerObserver{CatalogReferencerObserverFunc(consumer.Consume)},
							filteredOut: config.PostCatalogFiltersOut,
						}
						return chain.ConsumerFunc[[]*AnalysedMedia](adapter.OnBatchOfAnalysedMedia)
					},
					Next: &chain.MultithreadedLink[[]BackingUpMediaRequest, []BackingUpMediaRequest]{
						NumberOfRoutines: 1,
						ConsumerBuilder:  chain.PassThrough[[]BackingUpMediaRequest](),
						Next: &chain.CloseWrapperLink[[]BackingUpMediaRequest]{
							CloserFunc: config.closer,
							Next:       chain.EndOfTheChain[[]BackingUpMediaRequest](finalizer(config.PostCatalogFiltersIn)...),
						},
					},
				},
			},
		},
	}

	err := launcher.Starts(ctx, chain.NewErrorCollector())
	return launcher, err
}

func finalizer(in []CatalogReferencerObserver) []chain.ConsumerFunc[[]BackingUpMediaRequest] {
	functions := make([]chain.ConsumerFunc[[]BackingUpMediaRequest], len(in))
	for i, f := range in {
		functions[i] = f.OnMediaCatalogued
	}

	return functions
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
