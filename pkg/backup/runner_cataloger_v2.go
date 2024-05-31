package backup

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"slices"
)

// CatalogReference is used to project where a media will fit in the catalog: its ID and its album.
type CatalogReference interface {
	// Exists returns true if the media exists in the catalog
	Exists() bool

	// AlbumCreated returns true if the album was created during the cataloger process
	AlbumCreated() bool

	// AlbumFolderName returns the name of the album where the media would be stored
	AlbumFolderName() string

	// MediaId is used for backward compatibility with the previous version of the cataloger
	MediaId() string
}

type CatalogReferencer interface {
	Reference(ctx context.Context, medias []*AnalysedMedia) (map[*AnalysedMedia]CatalogReference, error)
}

type CatalogerFilter interface {
	// FilterOut returns true if the media should be filtered out, with the cause of behind excluded
	FilterOut(media AnalysedMedia, reference CatalogReference) (ProgressEventType, bool)
}

type CatalogerFilterFunc func(media AnalysedMedia, reference CatalogReference) (ProgressEventType, bool)

func (f CatalogerFilterFunc) FilterOut(media AnalysedMedia, reference CatalogReference) (ProgressEventType, bool) {
	return f(media, reference)
}

// Cataloger returns a BackingUpMediaRequest with an album all the time: it will have been created if necessary.
type Cataloger struct {
	CatalogReferencer CatalogReferencer
	CatalogerFilters  []CatalogerFilter
}

func (c *Cataloger) Catalog(ctx context.Context, medias []*AnalysedMedia, progressChannel chan *ProgressEvent) ([]*BackingUpMediaRequest, error) {
	references, err := c.CatalogReferencer.Reference(ctx, medias)
	if err != nil {
		return nil, err
	}

	var requests []*BackingUpMediaRequest
	counts := make(map[ProgressEventType]MediaCounter)

	for analysed, reference := range references {
		if cause, filteredOut := c.firstMatchingFilter(c.CatalogerFilters, *analysed, reference); filteredOut {
			count, _ := counts[cause]
			counts[cause] = count.Add(1, analysed.FoundMedia.Size())
			continue
		}

		if reference.AlbumCreated() {
			progressChannel <- &ProgressEvent{Type: ProgressEventAlbumCreated, Count: 1, Album: reference.AlbumFolderName()}
		}

		requests = append(requests, &BackingUpMediaRequest{
			AnalysedMedia: analysed,
			Id:            reference.MediaId(),
			FolderName:    reference.AlbumFolderName(),
		})
		count, _ := counts[ProgressEventCatalogued]
		counts[ProgressEventCatalogued] = count.Add(1, analysed.FoundMedia.Size())
	}

	for event, count := range counts {
		progressChannel <- &ProgressEvent{Type: event, Count: count.Count, Size: count.Size}
	}

	return requests, nil
}

func (c *Cataloger) firstMatchingFilter(filters []CatalogerFilter, media AnalysedMedia, reference CatalogReference) (ProgressEventType, bool) {
	for _, filter := range filters {
		if cause, filteredOut := filter.FilterOut(media, reference); filteredOut {
			return cause, true
		}
	}

	return ProgressEventCatalogued, false
}

func NewCataloger(owner ownermodel.Owner, options Options) (RunnerCataloger, error) {
	var referencer CatalogReferencer
	var err error

	if options.DryRun {
		referencer, err = referencerFactory.NewDryRunReferencer(context.TODO(), owner)
	} else {
		referencer, err = referencerFactory.NewCreatorReferencer(context.TODO(), owner)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create a cataloger for %s with options %+v", owner, options)
	}

	if len(options.RestrictedAlbumFolderName) > 0 {
		var albumFolderNames []string
		for albumFolderName := range options.RestrictedAlbumFolderName {
			albumFolderNames = append(albumFolderNames, albumFolderName)
		}

		return &Cataloger{
			CatalogReferencer: referencer,
			CatalogerFilters: []CatalogerFilter{
				mustNotExists(),
				mustBeInAlbum(albumFolderNames...),
			},
		}, nil
	}

	return &Cataloger{
		CatalogReferencer: referencer,
		CatalogerFilters: []CatalogerFilter{
			mustNotExists(),
		},
	}, nil
}

func mustBeInAlbum(albumFolderNames ...string) CatalogerFilterFunc {
	return func(media AnalysedMedia, reference CatalogReference) (ProgressEventType, bool) {
		return ProgressEventWrongAlbum, !slices.Contains(albumFolderNames, reference.AlbumFolderName())
	}
}

func mustNotExists() CatalogerFilterFunc {
	return func(media AnalysedMedia, reference CatalogReference) (ProgressEventType, bool) {
		return ProgressEventAlreadyExists, reference.Exists()
	}
}
