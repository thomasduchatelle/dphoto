package backup

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/backup/chain"
)

type scanCompleteObserver interface {
	OnScanComplete(ctx context.Context, count, size int) error
}

type analyserLauncher interface {
	Process(ctx context.Context, volume SourceVolume) chan error
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
				analyser := &analyserAggregate{
					analyser:               config.Analyser,
					analysedMediaObservers: []AnalysedMediaObserver{AnalysedMediaObserverFunc(consumer.Consume)},
					rejectedMediaObservers: config.PostAnalyserRejects,
				}
				return chain.ConsumerFunc[FoundMedia](analyser.OnFoundMedia)
			},
			Next: &chain.BufferLink[*AnalysedMedia]{
				BufferCapacity: options.BatchSize,
				Next: &chain.MultithreadedLink[[]*AnalysedMedia, []BackingUpMediaRequest]{
					NumberOfRoutines: options.ConcurrencyParameters.NumberOfConcurrentCataloguerRoutines(),
					ConsumerBuilder: func(consumer chain.Consumer[[]BackingUpMediaRequest]) chain.Consumer[[]*AnalysedMedia] {
						adapter := &cataloguerAggregate{
							cataloguer: config.Cataloguer,
							observerWithFilters: applyFiltersOnCataloguer{
								CatalogReferencerObservers: []CatalogReferencerObserver{CatalogReferencerObserverFunc(consumer.Consume)},
								CataloguerFilterObservers:  config.PostCataloguerFiltersOut,
								CataloguerFilters:          postCataloguerFiltersList(options),
							},
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
