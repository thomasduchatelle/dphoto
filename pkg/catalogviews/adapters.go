package catalogviews

import (
	"context"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

// AlbumViewCreateAlbumObserver adapts AlbumView.AlbumCreated to catalog.CreateAlbumObserver so
// the read-model can be plugged into catalog.CreateAlbum's observer chain.
type AlbumViewCreateAlbumObserver struct {
	AlbumView *AlbumView
}

func (a *AlbumViewCreateAlbumObserver) ObserveCreateAlbum(ctx context.Context, createdAlbum catalog.Album) error {
	return a.AlbumView.AlbumCreated(ctx, createdAlbum)
}
