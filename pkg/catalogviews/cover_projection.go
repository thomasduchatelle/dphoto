package catalogviews

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

// CoverProjection denormalises the canonical cover set of an album into every viewer's
// cover row. It implements catalog.CoversChangedObserver so cover-maintenance operations
// in pkg/catalog fire into it without a direct dependency on the view.
//
// Writing an empty set deletes the cover row for every viewer; writing a non-empty set
// upserts it. Fan-out across owner + visitors is handled by the repository.
type CoverProjection struct {
	Repository AlbumSummaryRepository
}

func (c *CoverProjection) OnCoversChanged(ctx context.Context, albumId catalog.AlbumId, covers []catalog.Cover) error {
	if len(covers) == 0 {
		if err := c.Repository.DeleteCoversForAllViewers(ctx, albumId); err != nil {
			return errors.Wrapf(err, "CoverProjection: failed to delete covers for %s", albumId)
		}
		return nil
	}

	if err := c.Repository.PutCoversForAllViewers(ctx, albumId, covers); err != nil {
		return errors.Wrapf(err, "CoverProjection: failed to put covers for %s", albumId)
	}
	return nil
}
