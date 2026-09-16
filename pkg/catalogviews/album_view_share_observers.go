package catalogviews

import (
	"context"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// AlbumViewAlbumSharedObserver adapts AlbumView to the catalogacl.AlbumSharedObserver contract
// so the album-list projection stays in sync when an album is shared with a new visitor.
type AlbumViewAlbumSharedObserver struct {
	AlbumView *AlbumView
}

func (a *AlbumViewAlbumSharedObserver) AlbumShared(ctx context.Context, album catalog.Album, userEmail usermodel.UserId) error {
	return a.AlbumView.AlbumShared(ctx, album, userEmail)
}

// AlbumViewAlbumUnSharedObserver adapts AlbumView to the catalogacl.AlbumUnSharedObserver contract
// so the visitor's projection row is removed when an album stops being shared.
type AlbumViewAlbumUnSharedObserver struct {
	AlbumView *AlbumView
}

func (a *AlbumViewAlbumUnSharedObserver) AlbumUnShared(ctx context.Context, albumId catalog.AlbumId, userEmail usermodel.UserId) error {
	return a.AlbumView.AlbumUnshared(ctx, albumId, userEmail)
}
