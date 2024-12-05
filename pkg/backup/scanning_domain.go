package backup

import (
	"context"
	"slices"
)

type CatalogReferencerObservers []CatalogReferencerObserver

func (c CatalogReferencerObservers) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	for _, observer := range c {
		if err := observer.OnMediaCatalogued(ctx, requests); err != nil {
			return err
		}
	}

	return nil
}

type analyserAdapter struct {
	analyser     Analyser
	analysed     AnalysedMediaObservers
	beforeFilter AnalysedMediaObservers
	filteredOut  RejectedMediaObservers
	rejected     RejectedMediaObservers
}

func (a *analyserAdapter) OnFoundMedia(ctx context.Context, media FoundMedia) error {
	return a.analyser.Analyse(
		ctx,
		media,
		slices.Concat(a.beforeFilter, []AnalysedMediaObserver{&analyserNoDateTimeFilter{
			analysedMediaObserver: a.analysed,
			rejectedMediaObserver: a.filteredOut,
		}}),
		&a.rejected,
	)
}

type cataloguerAdapter struct {
	cataloguer  Cataloguer
	options     Options
	preFilters  CatalogReferencerObservers
	catalogued  CatalogReferencerObservers
	filteredOut []CataloguerFilterObserver
}

func (s *cataloguerAdapter) OnBatchOfAnalysedMedia(ctx context.Context, batch []*AnalysedMedia) error {
	return s.cataloguer.Reference(ctx, batch, slices.Concat(
		s.preFilters,
		[]CatalogReferencerObserver{
			&applyFiltersOnCataloguer{
				CatalogReferencerObservers: s.catalogued,
				CataloguerFilterObservers:  s.filteredOut,
				CataloguerFilters:          postCatalogFiltersList(s.options),
			},
		},
	))
}

type analyserNoDateTimeFilter struct {
	analysedMediaObserver AnalysedMediaObserver
	rejectedMediaObserver RejectedMediaObserver
}

func (a *analyserNoDateTimeFilter) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	if media.Details.DateTime.IsZero() {
		return a.rejectedMediaObserver.OnRejectedMedia(ctx, media.FoundMedia, ErrAnalyserNoDateTime)
	}

	return a.analysedMediaObserver.OnAnalysedMedia(ctx, media)
}

// TODO should be removed ?

func newAnalysedMediaBufferAdapter(options Options, observer analysedMediasBatchObserver) *analysedMediaBufferAdapter {
	return &analysedMediaBufferAdapter{
		buffer: buffer[*AnalysedMedia]{
			consumer: observer.OnBatchOfAnalysedMedia,
			content:  make([]*AnalysedMedia, 0, defaultValue(options.BatchSize, 1)),
		},
	}
}

type analysedMediaBufferAdapter struct {
	buffer buffer[*AnalysedMedia]
}

func (a *analysedMediaBufferAdapter) OnAnalysedMedia(ctx context.Context, media *AnalysedMedia) error {
	return a.buffer.Append(ctx, media)
}

func (a *analysedMediaBufferAdapter) Close(ctx context.Context) error {
	return a.buffer.Flush(ctx)
}

type foundMediaObserver interface {
	OnFoundMedia(ctx context.Context, media FoundMedia) error
}

type volumeAdapter struct {
	observer foundMediaObserver
}

func (v *volumeAdapter) OnVolumeReceived(ctx context.Context, volume SourceVolume) error {
	medias, err := volume.FindMedias(ctx)
	if err != nil {
		return err
	}

	for _, media := range medias {
		if err := v.observer.OnFoundMedia(ctx, media); err != nil {
			return err
		}
	}

	return nil
}
