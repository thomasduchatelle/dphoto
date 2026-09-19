package catalogviews

import (
	"context"
	"slices"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// NewAlbumView constructs the AlbumView read model.
//
// The five collaborators are kept as fields so that the event methods, filled in by
// subsequent tickets (01-03..01-07), can be implemented without further signature churn on the
// constructor:
//   - Repository is the read+write projection carrying every viewer row.
//   - GetAlbumSharingGridPort decorates owned rows with the current sharing grid on read.
//   - MediaCounterPort re-reads canonical counts for count-mutation events.
//   - FindAlbumsByIdsPort loads canonical album records for events that need display fields.
//   - ListUsersWhoCanAccessAlbumPort resolves the set of viewers (owner + visitors) that
//     must see an album, used by every event that writes viewer rows.
func NewAlbumView(
	repository AlbumSummaryRepository,
	getAlbumSharingGridPort GetAlbumSharingGridPort,
	mediaCounterPort MediaCounterPort,
	findAlbumsByIdsPort FindAlbumsByIdsPort,
	listUsersWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort,
) *AlbumView {
	return &AlbumView{
		Repository:                     repository,
		GetAlbumSharingGridPort:        getAlbumSharingGridPort,
		MediaCounterPort:               mediaCounterPort,
		FindAlbumsByIdsPort:            findAlbumsByIdsPort,
		ListUsersWhoCanAccessAlbumPort: listUsersWhoCanAccessAlbumPort,
	}
}

// AlbumView is the album-list read model: it serves ListAlbums from the projection and
// keeps the projection in sync by observing catalog domain events.
type AlbumView struct {
	Repository                     AlbumSummaryRepository
	GetAlbumSharingGridPort        GetAlbumSharingGridPort
	MediaCounterPort               MediaCounterPort
	FindAlbumsByIdsPort            FindAlbumsByIdsPort
	ListUsersWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort
}

// ListAlbums returns the albums visible by the user (owned + shared) served from the
// projection alone: one ListSummariesForUser query, plus (when the user has an owner) one
// GetAlbumSharingGrid query to decorate owned rows with their visitors.
func (v *AlbumView) ListAlbums(ctx context.Context, user usermodel.CurrentUser, filter ListAlbumsFilter) ([]*VisibleAlbum, error) {
	summaries, err := v.Repository.ListSummariesForUser(ctx, user.UserId)
	if err != nil {
		return nil, err
	}

	var sharingGrid map[catalog.AlbumId][]usermodel.UserId
	if user.Owner != nil {
		sharingGrid, err = v.GetAlbumSharingGridPort.GetAlbumSharingGrid(ctx, *user.Owner)
		if err != nil {
			return nil, err
		}
	}

	var albums []*VisibleAlbum
	for _, summary := range summaries {
		if filter.OnlyDirectlyOwned && !summary.Availability.AsOwner {
			continue
		}

		visible := &VisibleAlbum{
			Album: catalog.Album{
				AlbumId: summary.AlbumSummary.AlbumId,
				Name:    summary.AlbumSummary.Name,
				Start:   summary.AlbumSummary.Start,
				End:     summary.AlbumSummary.End,
			},
			MediaCount:         summary.AlbumSummary.MediaCount,
			OwnedByCurrentUser: summary.Availability.AsOwner,
		}
		if summary.Availability.AsOwner {
			visible.Visitors = sharingGrid[summary.AlbumSummary.AlbumId]
		}
		albums = append(albums, visible)
	}

	slices.SortFunc(albums, func(a, b *VisibleAlbum) int {
		if a.Start.Equal(b.Start) {
			return b.End.Compare(a.End)
		}
		return b.Start.Compare(a.Start)
	})

	return albums, nil
}

// AlbumCreated writes a full summary row for every viewer of the newly created album
// (owner + any visitors already provisioned) and, when medias were transferred into the
// album at creation time, re-counts each source album whose medias were moved out.
//
// The destination row's Count is derived from the transfer payload itself, so no
// MediaCounterPort query is needed for the new album. The source-album counts are
// re-read from the canonical MediaCounterPort and written via SetCounts so the display
// fields on those rows are untouched.
func (v *AlbumView) AlbumCreated(ctx context.Context, event catalog.AlbumCreated) error {
	album := event.CreatedAlbum

	viewersByAlbum, err := v.ListUsersWhoCanAccessAlbumPort.ListUsersWhoCanAccessAlbum(ctx, album.AlbumId)
	if err != nil {
		return err
	}

	err = v.Repository.PutSummaries(ctx, []AlbumSummaryForUsers{
		{
			AlbumSummary: AlbumSummary{
				AlbumId:    album.AlbumId,
				Name:       album.Name,
				Start:      album.Start,
				End:        album.End,
				MediaCount: len(event.TransferredMedias.Transfers[album.AlbumId]),
			},
			Users: viewersByAlbum[album.AlbumId],
		},
	})
	if err != nil {
		return err
	}

	if len(event.TransferredMedias.FromAlbums) == 0 {
		return nil
	}

	sourceViewers, err := v.ListUsersWhoCanAccessAlbumPort.ListUsersWhoCanAccessAlbum(ctx, event.TransferredMedias.FromAlbums...)
	if err != nil {
		return err
	}

	sourceCounts, err := v.MediaCounterPort.CountMedia(ctx, event.TransferredMedias.FromAlbums...)
	if err != nil {
		return err
	}

	updates := make([]AlbumMediaCountForUsers, 0, len(event.TransferredMedias.FromAlbums))
	for _, sourceId := range event.TransferredMedias.FromAlbums {
		updates = append(updates, AlbumMediaCountForUsers{
			AlbumId:    sourceId,
			Users:      sourceViewers[sourceId],
			MediaCount: sourceCounts[sourceId],
		})
	}
	return v.Repository.SetCounts(ctx, updates)
}

// AlbumRenamedInPlace is filled in by ticket 01-05.
func (v *AlbumView) AlbumRenamedInPlace(ctx context.Context, albumId catalog.AlbumId, newName string) error {
	return nil
}

// AlbumDatesAmended is filled in by ticket 01-05.
func (v *AlbumView) AlbumDatesAmended(ctx context.Context, update catalog.DatesUpdate) error {
	return nil
}

// AlbumDeleted is filled in by ticket 01-06.
func (v *AlbumView) AlbumDeleted(ctx context.Context, albumId catalog.AlbumId) error {
	return nil
}

// AlbumShared is filled in by ticket 01-07.
func (v *AlbumView) AlbumShared(ctx context.Context, album catalog.Album, userId usermodel.UserId) error {
	return nil
}

// AlbumUnshared is filled in by ticket 01-07.
func (v *AlbumView) AlbumUnshared(ctx context.Context, albumId catalog.AlbumId, userId usermodel.UserId) error {
	return nil
}

// MediasInserted is filled in by ticket 01-04.
func (v *AlbumView) MediasInserted(ctx context.Context, medias map[catalog.AlbumId][]catalog.MediaId) error {
	return nil
}

// MediasTransferred is filled in by ticket 01-04.
func (v *AlbumView) MediasTransferred(ctx context.Context, transfers catalog.TransferredMedias) error {
	return nil
}
