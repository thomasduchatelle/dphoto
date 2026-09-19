package catalogviews

import (
	"context"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// AlbumViewAlbumSharedObserver adapts AlbumView.AlbumShared to the catalogacl.AlbumSharedObserver
// interface so the share use-case can notify the read model directly.
type AlbumViewAlbumSharedObserver struct {
	AlbumView *AlbumView
}

func (o *AlbumViewAlbumSharedObserver) AlbumShared(ctx context.Context, album catalog.Album, userEmail usermodel.UserId) error {
	return o.AlbumView.AlbumShared(ctx, album, userEmail)
}

// AlbumViewAlbumUnSharedObserver adapts AlbumView.AlbumUnshared to the
// catalogacl.AlbumUnSharedObserver interface.
type AlbumViewAlbumUnSharedObserver struct {
	AlbumView *AlbumView
}

func (o *AlbumViewAlbumUnSharedObserver) AlbumUnShared(ctx context.Context, albumId catalog.AlbumId, userEmail usermodel.UserId) error {
	return o.AlbumView.AlbumUnshared(ctx, albumId, userEmail)
}
