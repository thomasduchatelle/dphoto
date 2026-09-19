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

// AlbumCreated is filled in by ticket 01-03.
func (v *AlbumView) AlbumCreated(ctx context.Context, album catalog.Album) error {
	return nil
}

// AlbumRenamed keeps the album-list view in sync with a rename. When the folder name is
// unchanged (name-only edit) the existing viewer rows are updated in place through
// SetDisplayFields. When the folder changes, the old album's rows are wiped and a fresh
// row is written for the new album on every viewer — semantically equivalent to a delete
// followed by a create. In both cases, any transfer sources carried by the event are
// re-counted through SetCounts so display fields on those rows are not clobbered. The old
// album id is excluded from the recount when the folder changed: its rows have just been
// deleted and must not be resurrected without display fields.
func (v *AlbumView) AlbumRenamed(ctx context.Context, event catalog.AlbumRenamed) error {
	renamed := event.RenamedAlbum

	if event.ExistingAlbum.AlbumId == renamed.AlbumId {
		viewers, err := v.ListUsersWhoCanAccessAlbumPort.ListUsersWhoCanAccessAlbum(ctx, renamed.AlbumId)
		if err != nil {
			return err
		}
		if err = v.Repository.SetDisplayFields(ctx, renamed.AlbumId, viewers[renamed.AlbumId], renamed.Name, renamed.Start, renamed.End); err != nil {
			return err
		}
		return v.recountAlbums(ctx, filterAlbumIds(event.TransferredMedias.FromAlbums, event.ExistingAlbum.AlbumId))
	}

	if err := v.Repository.DeleteAllRowsForAlbum(ctx, event.ExistingAlbum.AlbumId); err != nil {
		return err
	}

	viewers, err := v.ListUsersWhoCanAccessAlbumPort.ListUsersWhoCanAccessAlbum(ctx, renamed.AlbumId)
	if err != nil {
		return err
	}
	err = v.Repository.PutSummaries(ctx, []AlbumSummaryForUsers{
		{
			AlbumSummary: AlbumSummary{
				AlbumId:    renamed.AlbumId,
				Name:       renamed.Name,
				Start:      renamed.Start,
				End:        renamed.End,
				MediaCount: len(event.TransferredMedias.Transfers[renamed.AlbumId]),
			},
			Users: viewers[renamed.AlbumId],
		},
	})
	if err != nil {
		return err
	}

	return v.recountAlbums(ctx, filterAlbumIds(event.TransferredMedias.FromAlbums, event.ExistingAlbum.AlbumId))
}

// AlbumDatesAmended keeps the album-list view in sync with an amend-dates change. The
// amended album's display fields (start/end) are updated in place on every viewer row via
// SetDisplayFields, and any transfer source/destination albums carried by the event are
// re-counted via SetCounts so their display fields survive.
func (v *AlbumView) AlbumDatesAmended(ctx context.Context, event catalog.AlbumDatesAmended) error {
	amended := event.DatesUpdate.UpdatedAlbum

	viewers, err := v.ListUsersWhoCanAccessAlbumPort.ListUsersWhoCanAccessAlbum(ctx, amended.AlbumId)
	if err != nil {
		return err
	}
	if err = v.Repository.SetDisplayFields(ctx, amended.AlbumId, viewers[amended.AlbumId], amended.Name, amended.Start, amended.End); err != nil {
		return err
	}

	var affected []catalog.AlbumId
	for albumId := range event.TransferredMedias.Transfers {
		if !slices.Contains(affected, albumId) {
			affected = append(affected, albumId)
		}
	}
	for _, albumId := range event.TransferredMedias.FromAlbums {
		if !slices.Contains(affected, albumId) {
			affected = append(affected, albumId)
		}
	}
	return v.recountAlbums(ctx, affected)
}

// recountAlbums re-reads the canonical media count of the given albums and writes each
// through SetCounts so the display fields on those rows are not touched.
func (v *AlbumView) recountAlbums(ctx context.Context, albumIds []catalog.AlbumId) error {
	if len(albumIds) == 0 {
		return nil
	}

	viewers, err := v.ListUsersWhoCanAccessAlbumPort.ListUsersWhoCanAccessAlbum(ctx, albumIds...)
	if err != nil {
		return err
	}
	counts, err := v.MediaCounterPort.CountMedia(ctx, albumIds...)
	if err != nil {
		return err
	}

	updates := make([]AlbumMediaCountForUsers, 0, len(albumIds))
	for _, albumId := range albumIds {
		updates = append(updates, AlbumMediaCountForUsers{
			AlbumId:    albumId,
			Users:      viewers[albumId],
			MediaCount: counts[albumId],
		})
	}
	return v.Repository.SetCounts(ctx, updates)
}

func filterAlbumIds(ids []catalog.AlbumId, exclude ...catalog.AlbumId) []catalog.AlbumId {
	var result []catalog.AlbumId
	for _, id := range ids {
		if !slices.Contains(exclude, id) {
			result = append(result, id)
		}
	}
	return result
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
