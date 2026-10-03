package catalog

import (
	"context"

	"github.com/pkg/errors"
)

// MediasInsertedCoverObserver is an InsertMediasObserver that reconciles the cover
// set of every album that received new medias. Because the just-inserted medias are
// among the eligible candidates returned by the album query, the invariant is applied
// via Reconcile with an empty `added` set: the fallback query naturally picks them up.
// This keeps the observer independent from the MediasInserted event payload.
type MediasInsertedCoverObserver struct {
	ReconcileCoversPort ReconcileCoversPort
}

func (m *MediasInsertedCoverObserver) OnMediasInserted(ctx context.Context, medias map[AlbumId][]MediaId) error {
	for albumId := range medias {
		if err := m.ReconcileCoversPort.Reconcile(ctx, albumId, nil, nil); err != nil {
			return errors.Wrapf(err, "MediasInsertedCoverObserver: failed to reconcile covers for %s", albumId)
		}
	}
	return nil
}
