package catalogviews

import (
	"context"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

// AlbumViewMediasInsertedObserver adapts AlbumView.MediasInserted to
// catalog.InsertMediasObserver so the InsertMedias use case can drive the read model.
type AlbumViewMediasInsertedObserver struct {
	AlbumView *AlbumView
}

func (o *AlbumViewMediasInsertedObserver) OnMediasInserted(ctx context.Context, medias map[catalog.AlbumId][]catalog.MediaId) error {
	return o.AlbumView.MediasInserted(ctx, medias)
}
