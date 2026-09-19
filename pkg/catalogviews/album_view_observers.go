package catalogviews

import (
	"context"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

// AlbumViewDeleteAlbumObserver adapts the catalog.AlbumDeletedObserver interface to the
// AlbumView read model: it forwards the domain event to AlbumView.AlbumDeleted so the
// projection wipes the deleted album's rows and refreshes the counts of any transfer
// destination albums.
type AlbumViewDeleteAlbumObserver struct {
	AlbumView *AlbumView
}

func (a *AlbumViewDeleteAlbumObserver) OnAlbumDeleted(ctx context.Context, event catalog.AlbumDeleted) error {
	return a.AlbumView.AlbumDeleted(ctx, event)
}
