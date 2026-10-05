package catalog

import (
	"context"

	"github.com/pkg/errors"
)

type AlbumCoversRandomised struct {
	AlbumId AlbumId
	Covers  []Cover
}

type AlbumCoversRandomisedObserver interface {
	OnAlbumCoversRandomised(ctx context.Context, event AlbumCoversRandomised) error
}

type AlbumCoversRandomisedObserverFunc func(ctx context.Context, event AlbumCoversRandomised) error

func (f AlbumCoversRandomisedObserverFunc) OnAlbumCoversRandomised(ctx context.Context, event AlbumCoversRandomised) error {
	return f(ctx, event)
}

func NewRandomizeAlbumCovers(
	coverService CoverServicePort,
	findCoversByAlbum FindCoversByAlbumPort,
	observers ...AlbumCoversRandomisedObserver,
) *RandomizeAlbumCovers {
	return &RandomizeAlbumCovers{
		CoverService:          coverService,
		FindCoversByAlbumPort: findCoversByAlbum,
		Observers:             observers,
	}
}

type FindCoversByAlbumPort interface {
	FindCoversByAlbum(ctx context.Context, albumId AlbumId) ([]Cover, error)
}

// RandomizeAlbumCovers is the owner-triggered use case to re-randomise the covers of a
// single album. It delegates the cover maintenance to the CoverService (non-stable
// variant: RANDOM covers are dropped and re-picked, CHERRY_PICKED covers are preserved,
// empty slots are filled up to MaxCoversPerAlbum). It fires AlbumCoversRandomised when
// the cover set actually changed, so view projections mirror the owner's action without
// doing a redundant write when the redraw happened to produce the same set.
type RandomizeAlbumCovers struct {
	CoverService          CoverServicePort
	FindCoversByAlbumPort FindCoversByAlbumPort
	Observers             []AlbumCoversRandomisedObserver
}

func (r *RandomizeAlbumCovers) Randomize(ctx context.Context, albumId AlbumId) ([]Cover, error) {
	changed, err := r.CoverService.Randomise(ctx, false, albumId)
	if err != nil {
		return nil, errors.Wrapf(err, "RandomizeAlbumCovers failed to randomise %s", albumId)
	}

	if covers, ok := changed[albumId]; ok {
		event := AlbumCoversRandomised{AlbumId: albumId, Covers: covers}
		for _, observer := range r.Observers {
			if err := observer.OnAlbumCoversRandomised(ctx, event); err != nil {
				return nil, errors.Wrapf(err, "failed to notify AlbumCoversRandomisedObserver %T", observer)
			}
		}
		return covers, nil
	}

	covers, err := r.FindCoversByAlbumPort.FindCoversByAlbum(ctx, albumId)
	if err != nil {
		return nil, errors.Wrapf(err, "RandomizeAlbumCovers failed to read current covers of %s", albumId)
	}
	return covers, nil
}
