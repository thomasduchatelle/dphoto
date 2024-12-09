package backup

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/backup/chain"
	"slices"
)

type scanCompleteObserver interface {
	OnScanComplete(ctx context.Context, count, size int) error
}

type analyserLauncher interface {
	Process(ctx context.Context, volume SourceVolume) chan error
}

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
			Next: &chain.BufferLink[*AnalysedMedia]{
				BufferCapacity: options.BatchSize,
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

func sizeOfAllMedias(medias []FoundMedia) int {
	size := 0
	for _, media := range medias {
		size += media.Size()
	}
	return size
}
