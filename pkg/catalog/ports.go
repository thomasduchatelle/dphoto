package catalog

import (
	"context"
	"time"

	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// TimelineRepository is the single port through which the album use cases both build the
// TimelineAggregate for an owner and persist mutations to individual albums. LoadTimeline
// returns a ready-to-use aggregate: it fails if the persisted albums cannot form a valid
// timeline (e.g. duplicated AlbumId).
type TimelineRepository interface {
	LoadTimeline(ctx context.Context, owner ownermodel.Owner) (*TimelineAggregate, error)
	InsertAlbum(ctx context.Context, album Album) error
	DeleteAlbum(ctx context.Context, albumId AlbumId) error
	UpdateAlbumName(ctx context.Context, albumId AlbumId, newName string) error
	AmendDates(ctx context.Context, albumId AlbumId, start, end time.Time) error
}
