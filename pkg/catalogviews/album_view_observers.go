package catalogviews

import (
	"context"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

// AlbumViewRenameObserver bridges the catalog domain's AlbumRenamedObserver contract to
// AlbumView.AlbumRenamed so the album-list read model stays in sync with rename events.
type AlbumViewRenameObserver struct {
	AlbumView *AlbumView
}

func (o *AlbumViewRenameObserver) OnAlbumRenamed(ctx context.Context, event catalog.AlbumRenamed) error {
	return o.AlbumView.AlbumRenamed(ctx, event)
}

// AlbumViewAmendDatesObserver bridges the catalog domain's AlbumDatesAmendedObserver
// contract to AlbumView.AlbumDatesAmended so the album-list read model stays in sync with
// date-amendment events.
type AlbumViewAmendDatesObserver struct {
	AlbumView *AlbumView
}

func (o *AlbumViewAmendDatesObserver) OnAlbumDatesAmended(ctx context.Context, event catalog.AlbumDatesAmended) error {
	return o.AlbumView.AlbumDatesAmended(ctx, event)
}
