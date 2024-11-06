package backup

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

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

type CatalogReferencerObserverFunc func(ctx context.Context, requests []BackingUpMediaRequest) error

func (c CatalogReferencerObserverFunc) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	return c(ctx, requests)
}

type CataloguerFilterObserver interface {
	OnFilteredOut(ctx context.Context, media AnalysedMedia, reference CatalogReference, cause error) error
}

type CataloguerObserver interface {
	CatalogReferencerObserver
	CataloguerFilterObserver
}

type Cataloguer interface {
	Reference(ctx context.Context, medias []*AnalysedMedia, observer CatalogReferencerObserver) error
}

type CataloguerWithFilters struct {
	Delegate                  Cataloguer
	CataloguerFilters         []CataloguerFilter
	CatalogReferencerObserver CatalogReferencerObserver
	CataloguerFilterObserver  CataloguerFilterObserver
}

func (c *CataloguerWithFilters) Catalog(ctx context.Context, medias []*AnalysedMedia) error {
	if c.CatalogReferencerObserver == nil || c.CataloguerFilterObserver == nil {
		return errors.New("cataloguer must have a CatalogReferencerObserver and a CataloguerFilterObserver")
	}
	filters := &applyFiltersOnCataloguer{
		CatalogReferencerObservers: []CatalogReferencerObserver{c.CatalogReferencerObserver},
		CataloguerFilterObservers:  []CataloguerFilterObserver{c.CataloguerFilterObserver},
		CataloguerFilters:          c.CataloguerFilters,
	}
	return c.Delegate.Reference(ctx, medias, filters)
}

func NewReferencer(owner ownermodel.Owner, dryRun bool) (Cataloguer, error) {
	var referencer Cataloguer
	var err error

	if dryRun {
		referencer, err = referencerFactory.NewDryRunCataloguer(context.TODO(), owner)
	} else {
		referencer, err = referencerFactory.NewAlbumCreatorCataloguer(context.TODO(), owner)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create a cataloguer for %s with dryRun=%t", owner, dryRun)
	}
	return referencer, nil
}
