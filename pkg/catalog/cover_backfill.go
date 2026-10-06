package catalog

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

type FindAlbumByOwnerPort interface {
	FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*Album, error)
}

type CoverBackfillObserver interface {
	OnCoverBackfilled(ctx context.Context, coversByAlbumId map[AlbumId][]Cover) error
}

// BackfillCovers is a utility to reconcile the covers of each album of an owner.
type BackfillCovers struct {
	FindAlbumByOwnerPort FindAlbumByOwnerPort
	CoverService         CoverServicePort
	Observers            []CoverBackfillObserver
}

func (b *BackfillCovers) BackfillForOwner(ctx context.Context, owner ownermodel.Owner) (map[AlbumId][]Cover, error) {
	albums, err := b.FindAlbumByOwnerPort.FindAlbumsByOwner(ctx, owner)
	if err != nil {
		return nil, errors.Wrapf(err, "BackfillCovers(%s) failed to list albums", owner)
	}

	albumIds := make([]AlbumId, len(albums), len(albums))
	for i, album := range albums {
		albumIds[i] = album.AlbumId
	}

	changed, err := b.CoverService.StableRandomise(ctx, albumIds...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to randomise albums of owner %s", owner)
	}

	for _, observer := range b.Observers {
		if err := observer.OnCoverBackfilled(ctx, changed); err != nil {
			return nil, errors.Wrapf(err, "failed to notify CoverBackfillObserver %T", observer)
		}
	}

	return changed, nil
}
