package catalogviews

import (
	"context"
	"slices"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// NewAlbumView constructs the AlbumView read model.
//
// The collaborators are kept as fields so that the event methods, filled in by subsequent
// tickets (01-03..01-07), can be implemented without further signature churn on the
// constructor:
//   - Repository is the read+write projection carrying every viewer row.
//   - GetAlbumSharingGridPort decorates owned rows with the current sharing grid on read.
//   - MediaCounterPort re-reads canonical counts for count-mutation events.
//   - FindAlbumsByIdsPort loads canonical album records for events that need display fields.
//   - ListUserWhoCanAccessAlbumPort resolves the set of viewer rows to touch on count-only
//     updates driven by domain events (e.g. transfer-destination recount on AlbumDeleted).
func NewAlbumView(
	repository AlbumSummaryRepository,
	getAlbumSharingGridPort GetAlbumSharingGridPort,
	mediaCounterPort MediaCounterPort,
	findAlbumsByIdsPort FindAlbumsByIdsPort,
	listUserWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort,
) *AlbumView {
	return &AlbumView{
		Repository:                    repository,
		GetAlbumSharingGridPort:       getAlbumSharingGridPort,
		MediaCounterPort:              mediaCounterPort,
		FindAlbumsByIdsPort:           findAlbumsByIdsPort,
		ListUserWhoCanAccessAlbumPort: listUserWhoCanAccessAlbumPort,
	}
}

// AlbumView is the album-list read model: it serves ListAlbums from the projection and
// keeps the projection in sync by observing catalog domain events.
type AlbumView struct {
	Repository                    AlbumSummaryRepository
	GetAlbumSharingGridPort       GetAlbumSharingGridPort
	MediaCounterPort              MediaCounterPort
	FindAlbumsByIdsPort           FindAlbumsByIdsPort
	ListUserWhoCanAccessAlbumPort ListUserWhoCanAccessAlbumPort
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

// AlbumRenamedInPlace is filled in by ticket 01-05.
func (v *AlbumView) AlbumRenamedInPlace(ctx context.Context, albumId catalog.AlbumId, newName string) error {
	return nil
}

// AlbumDatesAmended is filled in by ticket 01-05.
func (v *AlbumView) AlbumDatesAmended(ctx context.Context, update catalog.DatesUpdate) error {
	return nil
}

// AlbumDeleted wipes every viewer row for the deleted album and, when the deletion moved
// medias into surrounding albums, refreshes the count of those destination rows from the
// canonical MediaCounterPort. Display fields on the destination rows are left untouched
// (SET Count = :c only).
func (v *AlbumView) AlbumDeleted(ctx context.Context, event catalog.AlbumDeleted) error {
	if err := v.Repository.DeleteAllRowsForAlbum(ctx, event.DeletedAlbumId); err != nil {
		return err
	}

	if event.TransferredMedias.IsEmpty() {
		return nil
	}

	var destinationIds []catalog.AlbumId
	for albumId := range event.TransferredMedias.Transfers {
		destinationIds = append(destinationIds, albumId)
	}

	availabilities, err := v.ListUserWhoCanAccessAlbumPort.ListUsersWhoCanAccessAlbum(ctx, destinationIds...)
	if err != nil {
		return err
	}

	counts, err := v.MediaCounterPort.CountMedia(ctx, destinationIds...)
	if err != nil {
		return err
	}

	updates := make([]AlbumMediaCountForUsers, 0, len(destinationIds))
	for _, albumId := range destinationIds {
		updates = append(updates, AlbumMediaCountForUsers{
			AlbumId:    albumId,
			Users:      availabilities[albumId],
			MediaCount: counts[albumId],
		})
	}

	return v.Repository.SetCounts(ctx, updates)
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
