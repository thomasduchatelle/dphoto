package backup

import (
	"context"
)

type Cataloguer interface {
	Reference(ctx context.Context, medias []*AnalysedMedia, observer CatalogReferencerObserver) error
}

// CatalogReference is used to project where a media will fit in the catalog: its ID and its album.
type CatalogReference interface {
	// Exists returns true if the media exists in the catalog
	Exists() bool

	// AlbumCreated returns true if the album was created during the cataloger process
	AlbumCreated() bool

	// AlbumFolderName returns the name of the album where the media would be stored
	AlbumFolderName() string

	// UniqueIdentifier is identifying the media no matter its filename, its id in the catalog (if it's in it or not), its album, ... It's its signature.
	UniqueIdentifier() string

	// MediaId is the id of the media in the catalog and in the archive
	MediaId() string
}

type CatalogReferencerObserver interface {
	OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error
}

type CataloguerFilterObserver interface {
	OnFilteredOut(ctx context.Context, media AnalysedMedia, reference CatalogReference, cause error) error
}

type cataloguerAggregate struct {
	cataloguer          Cataloguer
	observerWithFilters applyFiltersOnCataloguer
}

func (s *cataloguerAggregate) OnBatchOfAnalysedMedia(ctx context.Context, batch []*AnalysedMedia) error {
	return s.cataloguer.Reference(ctx, batch, &s.observerWithFilters)
}

func postCataloguerFiltersList(options Options) []cataloguerFilter {
	filters := []cataloguerFilter{
		mustNotExists(),
		mustBeUniqueInVolume(),
	}

	if len(options.RestrictedAlbumFolderName) > 0 {
		var albumFolderNames []string
		for albumFolderName := range options.RestrictedAlbumFolderName {
			albumFolderNames = append(albumFolderNames, albumFolderName)
		}
		filters = append(filters, mustBeInAlbum(albumFolderNames...))
	}

	return filters
}

type applyFiltersOnCataloguer struct {
	CatalogReferencerObservers CatalogReferencerObservers
	CataloguerFilterObservers  CataloguerFilterObservers
	CataloguerFilters          []cataloguerFilter
}

func (a *applyFiltersOnCataloguer) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	validRequests := make([]BackingUpMediaRequest, 0, len(requests))
	for _, request := range requests {
		if filterErr := a.filteredOut(ctx, request); filterErr != nil {
			err := a.CataloguerFilterObservers.OnFilteredOut(ctx, *request.AnalysedMedia, request.CatalogReference, filterErr)
			if err != nil {
				return err
			}

		} else {
			validRequests = append(validRequests, request)
		}
	}

	if len(validRequests) > 0 {
		return a.CatalogReferencerObservers.OnMediaCatalogued(ctx, validRequests)
	}

	return nil
}

func (a *applyFiltersOnCataloguer) filteredOut(ctx context.Context, request BackingUpMediaRequest) error {
	for _, filter := range a.CataloguerFilters {
		if err := filter.FilterOut(ctx, *request.AnalysedMedia, request.CatalogReference); err != nil {
			return err
		}
	}

	return nil
}

type CatalogReferencerObserverFunc func(ctx context.Context, requests []BackingUpMediaRequest) error

func (c CatalogReferencerObserverFunc) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	return c(ctx, requests)
}

type CatalogReferencerObservers []CatalogReferencerObserver

func (c CatalogReferencerObservers) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	for _, observer := range c {
		if err := observer.OnMediaCatalogued(ctx, requests); err != nil {
			return err
		}
	}

	return nil
}

type CataloguerFilterObservers []CataloguerFilterObserver

func (c CataloguerFilterObservers) OnFilteredOut(ctx context.Context, media AnalysedMedia, reference CatalogReference, cause error) error {
	for _, observer := range c {
		if err := observer.OnFilteredOut(ctx, media, reference, cause); err != nil {
			return err
		}
	}

	return nil
}
