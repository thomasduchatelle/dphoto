package catalogviews

import (
	"context"
	"slices"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

func NewAlbumView(
	repository AlbumSummaryRepository,
	getAlbumSharingGridPort GetAlbumSharingGridPort,
	mediaCounterPort MediaCounterPort,
	findAlbumsByIdsPort FindAlbumsByIdsPort,
	ownerUserIdPort OwnerUserIdPort,
	findCoversByAlbumPort FindCoversByAlbumPort,
) *AlbumView {
	return &AlbumView{
		Repository:              repository,
		GetAlbumSharingGridPort: getAlbumSharingGridPort,
		MediaCounterPort:        mediaCounterPort,
		FindAlbumsByIdsPort:     findAlbumsByIdsPort,
		OwnerUserIdPort:         ownerUserIdPort,
		FindCoversByAlbumPort:   findCoversByAlbumPort,
	}
}

type AlbumView struct {
	Repository              AlbumSummaryRepository
	GetAlbumSharingGridPort GetAlbumSharingGridPort
	MediaCounterPort        MediaCounterPort
	FindAlbumsByIdsPort     FindAlbumsByIdsPort
	OwnerUserIdPort         OwnerUserIdPort
	FindCoversByAlbumPort   FindCoversByAlbumPort
}

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
			Covers:             summary.AlbumSummary.Covers,
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

func (v *AlbumView) OnAlbumCreated(ctx context.Context, event catalog.AlbumCreated) error {
	album := event.CreatedAlbum

	ownerUserId, err := v.OwnerUserIdPort.GetOwnerUserId(ctx, album.AlbumId.Owner)
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
			Users: []Availability{OwnerAvailability(ownerUserId)},
		},
	})
	if err != nil {
		return err
	}

	return v.recountAlbums(ctx, event.TransferredMedias.FromAlbums)
}

func (v *AlbumView) OnAlbumRenamed(ctx context.Context, event catalog.AlbumRenamed) error {
	renamed := event.RenamedAlbum

	if event.ExistingAlbum.AlbumId.IsEqual(renamed.AlbumId) {
		return v.Repository.SetDisplayFieldsForAllViewers(ctx, renamed.AlbumId, renamed.Name, renamed.Start, renamed.End)
	}

	return v.Repository.RenameAlbum(ctx, event.ExistingAlbum.AlbumId, renamed.AlbumId, renamed.Name)
}

func (v *AlbumView) OnAlbumDatesAmended(ctx context.Context, event catalog.AlbumDatesAmended) error {
	amended := event.DatesUpdate.UpdatedAlbum

	if err := v.Repository.SetDisplayFieldsForAllViewers(ctx, amended.AlbumId, amended.Name, amended.Start, amended.End); err != nil {
		return err
	}

	var affected []catalog.AlbumId
	for albumId := range event.TransferredMedias.Transfers {
		if !albumId.IsEqual(amended.AlbumId) && !slices.ContainsFunc(affected, albumId.IsEqual) {
			affected = append(affected, albumId)
		}
	}
	for _, albumId := range event.TransferredMedias.FromAlbums {
		if !albumId.IsEqual(amended.AlbumId) && !slices.ContainsFunc(affected, albumId.IsEqual) {
			affected = append(affected, albumId)
		}
	}
	return v.recountAlbums(ctx, affected)
}

func (v *AlbumView) OnAlbumDeleted(ctx context.Context, event catalog.AlbumDeleted) error {
	if err := v.Repository.DeleteAllRowsForAlbum(ctx, event.DeletedAlbumId); err != nil {
		return err
	}

	if event.TransferredMedias.IsEmpty() {
		return nil
	}

	destinationIds := make([]catalog.AlbumId, 0, len(event.TransferredMedias.Transfers))
	for albumId := range event.TransferredMedias.Transfers {
		destinationIds = append(destinationIds, albumId)
	}

	return v.recountAlbums(ctx, destinationIds)
}

func (v *AlbumView) AlbumShared(ctx context.Context, album catalog.Album, userId usermodel.UserId) error {
	counts, err := v.MediaCounterPort.CountMedia(ctx, album.AlbumId)
	if err != nil {
		return err
	}

	covers, err := v.FindCoversByAlbumPort.FindCoversByAlbum(ctx, album.AlbumId)
	if err != nil {
		return err
	}

	return v.Repository.PutSummaries(ctx, []AlbumSummaryForUsers{
		{
			AlbumSummary: AlbumSummary{
				AlbumId:    album.AlbumId,
				MediaCount: counts[album.AlbumId],
				Name:       album.Name,
				Start:      album.Start,
				End:        album.End,
				Covers:     covers,
			},
			Users: []Availability{VisitorAvailability(userId)},
		},
	})
}

func (v *AlbumView) AlbumUnShared(ctx context.Context, albumId catalog.AlbumId, userId usermodel.UserId) error {
	return v.Repository.DeleteRow(ctx, VisitorAvailability(userId), albumId)
}

func (v *AlbumView) OnMediasInserted(ctx context.Context, event catalog.MediasInserted) error {
	if len(event.Inserted) == 0 && len(event.Covers) == 0 {
		return nil
	}

	if len(event.Inserted) > 0 {
		diffs := make([]AlbumCountDiff, 0, len(event.Inserted))
		for albumId, mediaIds := range event.Inserted {
			diffs = append(diffs, AlbumCountDiff{
				AlbumId:        albumId,
				MediaCountDiff: len(mediaIds),
			})
		}
		if err := v.Repository.IncrementCountForAllViewers(ctx, diffs); err != nil {
			return err
		}
	}

	return v.applyCoverUpdates(ctx, event.Covers)
}

// OnCoverBackfilled denormalises the cover sets produced by an admin backfill into every
// viewer's cover row. The backfill is the only path that writes covers outside of a
// lifecycle event, so it needs its own observer hook.
func (v *AlbumView) OnCoverBackfilled(ctx context.Context, coversByAlbumId map[catalog.AlbumId][]catalog.Cover) error {
	return v.applyCoverUpdates(ctx, coversByAlbumId)
}

func (v *AlbumView) applyCoverUpdates(ctx context.Context, covers map[catalog.AlbumId][]catalog.Cover) error {
	for albumId, albumCovers := range covers {
		if len(albumCovers) == 0 {
			if err := v.Repository.DeleteCoversForAllViewers(ctx, albumId); err != nil {
				return err
			}
			continue
		}
		if err := v.Repository.PutCoversForAllViewers(ctx, albumId, albumCovers); err != nil {
			return err
		}
	}
	return nil
}

func (v *AlbumView) recountAlbums(ctx context.Context, albumIds []catalog.AlbumId) error {
	if len(albumIds) == 0 {
		return nil
	}

	counts, err := v.MediaCounterPort.CountMedia(ctx, albumIds...)
	if err != nil {
		return err
	}

	updates := make([]AlbumCount, 0, len(albumIds))
	for _, albumId := range albumIds {
		updates = append(updates, AlbumCount{
			AlbumId:    albumId,
			MediaCount: counts[albumId],
		})
	}
	return v.Repository.SetCountForAllViewers(ctx, updates)
}
