package catalogviews

import (
	"context"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

// AlbumViewCreateAlbumObserver adapts AlbumView.AlbumCreated to the catalog.AlbumCreatedObserver
// contract so the read model can subscribe to CreateAlbum's event fan-out.
type AlbumViewCreateAlbumObserver struct {
	AlbumView *AlbumView
}

func (a *AlbumViewCreateAlbumObserver) OnAlbumCreated(ctx context.Context, event catalog.AlbumCreated) error {
	return a.AlbumView.AlbumCreated(ctx, event)
}
