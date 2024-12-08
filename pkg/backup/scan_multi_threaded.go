package backup

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/backup/chain"
	"slices"
)

func multithreadedScanRuntime(ctxNonCancelable context.Context, options Options, config *scanConfiguration) (analyserLauncher, error) {
	ctx, cancelFunc := context.WithCancel(ctxNonCancelable)

	launcher := scanAndBackupCommonLauncher(config, options, &chain.MultithreadedLink[[]BackingUpMediaRequest, []BackingUpMediaRequest]{
		NumberOfRoutines: 1,
		ConsumerBuilder:  chain.PassThrough[[]BackingUpMediaRequest](),
		Next: &chain.CloseWrapperLink[[]BackingUpMediaRequest]{
			CloserFuncs: slices.Concat(config.Wrappers, []chain.CloserFunc{chain.CloserFunc(cancelFunc)}),
			Next:        chain.EndOfTheChain[[]BackingUpMediaRequest](finalizer(config.PostCatalogFiltersIn)...),
		},
	})

	err := launcher.Starts(ctx, chain.NewErrorCollector(func(err error) {
		cancelFunc()
	}))
	return launcher, err
}

func scanAndBackupCommonLauncher(config *scanConfiguration, options Options, next chain.Link[[]BackingUpMediaRequest]) *chain.SingleLauncher[SourceVolume, FoundMedia] {
	return &chain.SingleLauncher[SourceVolume, FoundMedia]{
		Function: func(ctx context.Context, volume SourceVolume) ([]FoundMedia, error) {
			medias, err := volume.FindMedias(ctx)
			if err != nil || config.ScanCompleteObserver == nil {
				return medias, err
			}

			return medias, config.ScanCompleteObserver.OnScanComplete(ctx, len(medias), sizeOfAllMedias(medias))
		},
		Next: &chain.MultithreadedLink[FoundMedia, *AnalysedMedia]{
			NumberOfRoutines: options.ConcurrencyParameters.NumberOfConcurrentAnalyserRoutines(),
			Cancellable:      true,
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
			Next: &bufferLink[*AnalysedMedia]{
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
					Next: next,
				},
			},
		},
	}
}

func finalizer(in []CatalogReferencerObserver) []chain.ConsumerFunc[[]BackingUpMediaRequest] {
	functions := make([]chain.ConsumerFunc[[]BackingUpMediaRequest], len(in))
	for i, f := range in {
		functions[i] = f.OnMediaCatalogued
	}

	return functions
}

type bufferLink[Consumed any] struct {
	Next    chain.Link[[]Consumed]
	channel chan Consumed
	Buffer  buffer[Consumed]
}

func (l *bufferLink[Consumed]) Consume(ctx context.Context, consumed Consumed) error {
	l.channel <- consumed
	return nil
}

func (l *bufferLink[Consumed]) Starts(ctx context.Context, collector chain.ChainableErrorCollector) error {
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

func (l *bufferLink[Consumed]) WaitForCompletion() chan error {
	return l.Next.WaitForCompletion()
}

func (l *bufferLink[Consumed]) NotifyUpstreamCompleted() {
	close(l.channel)
}
