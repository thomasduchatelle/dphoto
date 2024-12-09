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
