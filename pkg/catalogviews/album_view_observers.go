package catalogviews

import (
	"context"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

// AlbumViewRenameObserver implements catalog.RenameAlbumObserver by forwarding in-place rename
// events to AlbumView.AlbumRenamedInPlace so the projection stays in sync with the canonical
// album name.
type AlbumViewRenameObserver struct {
	AlbumView *AlbumView
}

func (a *AlbumViewRenameObserver) OnAlbumRenamed(ctx context.Context, albumId catalog.AlbumId, newName string) error {
	return a.AlbumView.AlbumRenamedInPlace(ctx, albumId, newName)
}

// AlbumViewAmendDatesObserver implements catalog.AlbumDatesAmendedObserver by forwarding
// dates-amended events to AlbumView.AlbumDatesAmended so the projection stays in sync with
// the canonical album dates.
type AlbumViewAmendDatesObserver struct {
	AlbumView *AlbumView
}

func (a *AlbumViewAmendDatesObserver) OnAlbumDatesAmended(ctx context.Context, amendedAlbum catalog.DatesUpdate) error {
	return a.AlbumView.AlbumDatesAmended(ctx, amendedAlbum)
}
