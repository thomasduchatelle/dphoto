package catalogviews

import (
	"context"
	"slices"

	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// NewAlbumView constructs the AlbumView read model.
//
// Collaborators kept as fields so subsequent event methods can be filled in without churn:
//   - Repository is the read+write projection carrying every viewer row.
//   - GetAlbumSharingGridPort decorates owned rows with the current sharing grid on read.
//   - MediaCounterPort re-reads canonical counts for count-mutation events.
//   - FindAlbumsByIdsPort loads canonical album records; used by AlbumRenamedInPlace to
//     preserve Start/End (the observer signature only carries the new name).
//   - ListUserWhoCanAccessAlbumPort fans display-field updates out to every viewer row.
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

// AlbumRenamedInPlace refreshes AlbumName on every viewer row for the given album while
// leaving Start/End at their current canonical values. The rename observer signature only
// carries the new name, so the current album is loaded via FindAlbumsByIdsPort to source
// Start/End (SetDisplayFields SETs the three attributes together on the DynamoDB adapter).
func (v *AlbumView) AlbumRenamedInPlace(ctx context.Context, albumId catalog.AlbumId, newName string) error {
	albums, err := v.FindAlbumsByIdsPort.FindAlbumsById(ctx, []catalog.AlbumId{albumId})
	if err != nil {
		return err
	}
	if len(albums) == 0 {
		return errors.Wrapf(catalog.AlbumNotFoundErr, "AlbumRenamedInPlace(%s)", albumId)
	}

	users, err := v.ListUserWhoCanAccessAlbumPort.ListUsersWhoCanAccessAlbum(ctx, albumId)
	if err != nil {
		return err
	}

	return v.Repository.SetDisplayFields(ctx, albumId, users[albumId], newName, albums[0].Start, albums[0].End)
}

// AlbumDatesAmended writes the new Start/End (and preserves the current Name, carried by
// DatesUpdate.UpdatedAlbum) onto every viewer row for the amended album.
func (v *AlbumView) AlbumDatesAmended(ctx context.Context, update catalog.DatesUpdate) error {
	albumId := update.UpdatedAlbum.AlbumId
	users, err := v.ListUserWhoCanAccessAlbumPort.ListUsersWhoCanAccessAlbum(ctx, albumId)
	if err != nil {
		return err
	}

	return v.Repository.SetDisplayFields(ctx, albumId, users[albumId], update.UpdatedAlbum.Name, update.UpdatedAlbum.Start, update.UpdatedAlbum.End)
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
