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
) *AlbumView {
	return &AlbumView{
		Repository:              repository,
		GetAlbumSharingGridPort: getAlbumSharingGridPort,
		MediaCounterPort:        mediaCounterPort,
		FindAlbumsByIdsPort:     findAlbumsByIdsPort,
		OwnerUserIdPort:         ownerUserIdPort,
	}
}

type AlbumView struct {
	Repository              AlbumSummaryRepository
	GetAlbumSharingGridPort GetAlbumSharingGridPort
	MediaCounterPort        MediaCounterPort
	FindAlbumsByIdsPort     FindAlbumsByIdsPort
	OwnerUserIdPort         OwnerUserIdPort
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

	return v.Repository.PutSummaries(ctx, []AlbumSummaryForUsers{
		{
			AlbumSummary: AlbumSummary{
				AlbumId:    album.AlbumId,
				MediaCount: counts[album.AlbumId],
				Name:       album.Name,
				Start:      album.Start,
				End:        album.End,
			},
			Users: []Availability{VisitorAvailability(userId)},
		},
	})
}

func (v *AlbumView) AlbumUnShared(ctx context.Context, albumId catalog.AlbumId, userId usermodel.UserId) error {
	return v.Repository.DeleteRow(ctx, VisitorAvailability(userId), albumId)
}

func (v *AlbumView) OnMediasInserted(ctx context.Context, medias map[catalog.AlbumId][]catalog.MediaId) error {
	if len(medias) == 0 {
		return nil
	}

	diffs := make([]AlbumCountDiff, 0, len(medias))
	for albumId, mediaIds := range medias {
		diffs = append(diffs, AlbumCountDiff{
			AlbumId:        albumId,
			MediaCountDiff: len(mediaIds),
		})
	}

	return v.Repository.IncrementCountForAllViewers(ctx, diffs)
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
