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

type CataloguerFilterObserver interface {
	OnFilteredOut(ctx context.Context, media AnalysedMedia, reference CatalogReference, cause error) error
}

type CataloguerObserver interface {
	CatalogReferencerObserver
	CataloguerFilterObserver
}

type CatalogReferencer interface {
	Reference(ctx context.Context, medias []*AnalysedMedia, observer CatalogReferencerObserver) error
}

type Cataloguer interface {
	Catalog(ctx context.Context, medias []*AnalysedMedia, observer CataloguerObserver) error
}

type CataloguerWithFilters struct {
	Delegate          CatalogReferencer
	CataloguerFilters []CataloguerFilter
}

func (c *CataloguerWithFilters) Catalog(ctx context.Context, medias []*AnalysedMedia, observer CataloguerObserver) error {
	filters := &ApplyFiltersOnCataloguer{
		Delegate:          observer,
		Observer:          observer,
		CataloguerFilters: c.CataloguerFilters,
	}
	return c.Delegate.Reference(ctx, medias, filters)
}

func NewCataloguer(owner ownermodel.Owner, options Options) (Cataloguer, error) {
	var referencer CatalogReferencer
	var err error

	if options.DryRun {
		referencer, err = referencerFactory.NewDryRunReferencer(context.TODO(), owner)
	} else {
		referencer, err = referencerFactory.NewCreatorReferencer(context.TODO(), owner)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create a cataloguer for %s with options %+v", owner, options)
	}

	if len(options.RestrictedAlbumFolderName) > 0 {
		var albumFolderNames []string
		for albumFolderName := range options.RestrictedAlbumFolderName {
			albumFolderNames = append(albumFolderNames, albumFolderName)
		}

		return &CataloguerWithFilters{
			Delegate: referencer,
			CataloguerFilters: []CataloguerFilter{
				mustNotExists(),
				mustBeInAlbum(albumFolderNames...),
			},
		}, nil
	}

	return &CataloguerWithFilters{
		Delegate: referencer,
		CataloguerFilters: []CataloguerFilter{
			mustNotExists(),
		},
	}, nil
}
