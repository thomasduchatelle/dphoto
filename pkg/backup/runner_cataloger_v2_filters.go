package backup

import (
	"context"
	"github.com/pkg/errors"
	"slices"
)

var (
	CatalogerFilterMustBeInAlbumError        = errors.New("media must be in album")
	CatalogerFilterMustNotAlreadyExistsError = errors.New("media must not already exists")
)

type CataloguerFilter interface {
	// FilterOut returns an error if the media must be filtered out
	FilterOut(ctx context.Context, media AnalysedMedia, reference CatalogReference) error
}

func mustBeInAlbum(albumFolderNames ...string) CataloguerFilter {
	return &mustBeInAlbumCatalogerFilter{albumFolderNames: albumFolderNames}
}

func mustNotExists() CataloguerFilter {
	return new(mustNotAlreadyExistsCatalogerFilter)
}

type mustBeInAlbumCatalogerFilter struct {
	albumFolderNames []string
}

func (m mustBeInAlbumCatalogerFilter) FilterOut(ctx context.Context, media AnalysedMedia, reference CatalogReference) error {
	if slices.Contains(m.albumFolderNames, reference.AlbumFolderName()) {
		return nil
	}

	return CatalogerFilterMustBeInAlbumError
}

type mustNotAlreadyExistsCatalogerFilter struct{}

func (m mustNotAlreadyExistsCatalogerFilter) FilterOut(ctx context.Context, media AnalysedMedia, reference CatalogReference) error {
	if reference.Exists() {
		return CatalogerFilterMustNotAlreadyExistsError
	}
	return nil
}

type ApplyFiltersOnCataloguer struct {
	Delegate          CatalogReferencerObserver
	Observer          CataloguerFilterObserver
	CataloguerFilters []CataloguerFilter
}

func (a *ApplyFiltersOnCataloguer) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	validRequests := make([]BackingUpMediaRequest, 0, len(requests))
	for _, request := range requests {
		filtered, err := a.filteredOut(ctx, request)
		if err != nil {
			return err
		}

		if !filtered {
			validRequests = append(validRequests, request)
		}
	}

	if len(validRequests) > 0 {
		return a.Delegate.OnMediaCatalogued(ctx, validRequests)
	}

	return nil
}

func (a *ApplyFiltersOnCataloguer) filteredOut(ctx context.Context, request BackingUpMediaRequest) (bool, error) {
	for _, filter := range a.CataloguerFilters {
		if err := filter.FilterOut(ctx, *request.AnalysedMedia, request.CatalogReference); err != nil {
			return true, a.Observer.OnFilteredOut(ctx, *request.AnalysedMedia, request.CatalogReference, err)
		}
	}

	return false, nil
}
