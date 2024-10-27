package backup

import (
	"context"
	"github.com/pkg/errors"
	"slices"
)

var (
	ErrAnalyserNoDateTime                  = errors.New("media must have a date time included in the metadata")
	ErrCatalogerFilterMustBeInAlbum        = errors.New("media must be in album")
	ErrCatalogerFilterMustNotAlreadyExists = errors.New("media must not already exists")
	ErrMediaMustNotBeDuplicated            = errors.New("media is present twice in the volume")
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

func mustBeUniqueInVolume() CataloguerFilter {
	return &uniqueFilter{
		uniqueIndexes: make(map[string]interface{}),
	}
}

type applyFiltersOnCataloguer struct {
	CatalogReferencerObservers []CatalogReferencerObserver
	CataloguerFilterObservers  []CataloguerFilterObserver
	CataloguerFilters          []CataloguerFilter
}

func (a *applyFiltersOnCataloguer) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
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
		return a.fireOnMediaCatalogued(ctx, validRequests)
	}

	return nil
}

func (a *applyFiltersOnCataloguer) filteredOut(ctx context.Context, request BackingUpMediaRequest) (bool, error) {
	for _, filter := range a.CataloguerFilters {
		if err := filter.FilterOut(ctx, *request.AnalysedMedia, request.CatalogReference); err != nil {
			return true, a.fireOnFilteredOut(ctx, request, err)
		}
	}

	return false, nil
}

func (a *applyFiltersOnCataloguer) fireOnMediaCatalogued(ctx context.Context, validRequests []BackingUpMediaRequest) error {
	for _, observer := range a.CatalogReferencerObservers {
		if err := observer.OnMediaCatalogued(ctx, validRequests); err != nil {
			return err
		}
	}

	return nil
}

func (a *applyFiltersOnCataloguer) fireOnFilteredOut(ctx context.Context, request BackingUpMediaRequest, cause error) error {
	for _, observer := range a.CataloguerFilterObservers {
		if err := observer.OnFilteredOut(ctx, *request.AnalysedMedia, request.CatalogReference, cause); err != nil {
			return err
		}
	}

	return nil
}

type mustBeInAlbumCatalogerFilter struct {
	albumFolderNames []string
}

func (m mustBeInAlbumCatalogerFilter) FilterOut(ctx context.Context, media AnalysedMedia, reference CatalogReference) error {
	if slices.Contains(m.albumFolderNames, reference.AlbumFolderName()) {
		return nil
	}

	return ErrCatalogerFilterMustBeInAlbum
}

type mustNotAlreadyExistsCatalogerFilter struct{}

func (m mustNotAlreadyExistsCatalogerFilter) FilterOut(ctx context.Context, media AnalysedMedia, reference CatalogReference) error {
	if reference.Exists() {
		return ErrCatalogerFilterMustNotAlreadyExists
	}
	return nil
}

type uniqueFilter struct {
	uniqueIndexes map[string]interface{}
}

func (u *uniqueFilter) FilterOut(ctx context.Context, media AnalysedMedia, reference CatalogReference) error {
	uniqueId := reference.UniqueIdentifier()
	if _, filterOut := u.uniqueIndexes[uniqueId]; filterOut {
		return ErrMediaMustNotBeDuplicated
	}

	u.uniqueIndexes[uniqueId] = nil
	return nil
}
