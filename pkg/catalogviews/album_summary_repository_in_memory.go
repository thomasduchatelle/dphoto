package catalogviews

import (
	"context"
	"slices"
	"time"

	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type UserAlbumSummary struct {
	AlbumSummary AlbumSummary
	Availability Availability
}

func (s UserAlbumSummary) ToSummaryForUsers() AlbumSummaryForUsers {
	return AlbumSummaryForUsers{AlbumSummary: s.AlbumSummary, Users: []Availability{s.Availability}}
}

type AlbumSummaryInMemoryRepository struct {
	Summaries                 []UserAlbumSummary
	LegacyRowsCleanedForUsers []usermodel.UserId
}

func (r *AlbumSummaryInMemoryRepository) ListSummariesForUser(ctx context.Context, userId usermodel.UserId) ([]UserAlbumSummary, error) {
	var summaries []UserAlbumSummary
	for _, summary := range r.Summaries {
		if summary.Availability.UserId == userId {
			summaries = append(summaries, summary)
		}
	}

	return summaries, nil
}

func (r *AlbumSummaryInMemoryRepository) ListSummariesForUserAndOwners(ctx context.Context, userId usermodel.UserId, owner ...ownermodel.Owner) ([]UserAlbumSummary, error) {
	var summaries []UserAlbumSummary
	for _, summary := range r.Summaries {
		if summary.Availability.UserId == userId && slices.Contains(owner, summary.AlbumSummary.AlbumId.Owner) {
			summaries = append(summaries, summary)
		}
	}

	return summaries, nil
}

func (r *AlbumSummaryInMemoryRepository) PutSummaries(ctx context.Context, summaries []AlbumSummaryForUsers) error {
	if summaries == nil {
		return errors.Errorf("PutSummaries(nil): summaries should not be nil")
	}

	for _, summary := range summaries {
		for _, user := range summary.Users {
			index := slices.IndexFunc(r.Summaries, func(current UserAlbumSummary) bool {
				return current.AlbumSummary.AlbumId.IsEqual(summary.AlbumId) && current.Availability.UserId == user.UserId
			})
			newSummary := summary.AlbumSummary
			if len(newSummary.Covers) == 0 && index >= 0 {
				newSummary.Covers = r.Summaries[index].AlbumSummary.Covers
			}
			if index >= 0 {
				r.Summaries[index] = UserAlbumSummary{AlbumSummary: newSummary, Availability: user}
			} else {
				r.Summaries = append(r.Summaries, UserAlbumSummary{AlbumSummary: newSummary, Availability: user})
			}
		}
	}

	return nil
}

func (r *AlbumSummaryInMemoryRepository) SetDisplayFieldsForAllViewers(ctx context.Context, albumId catalog.AlbumId, name string, start, end time.Time) error {
	for i := range r.Summaries {
		if r.Summaries[i].AlbumSummary.AlbumId.IsEqual(albumId) {
			r.Summaries[i].AlbumSummary.Name = name
			r.Summaries[i].AlbumSummary.Start = start
			r.Summaries[i].AlbumSummary.End = end
		}
	}
	return nil
}

func (r *AlbumSummaryInMemoryRepository) PutCoversForAllViewers(ctx context.Context, albumId catalog.AlbumId, covers []catalog.Cover) error {
	for i := range r.Summaries {
		if r.Summaries[i].AlbumSummary.AlbumId.IsEqual(albumId) {
			r.Summaries[i].AlbumSummary.Covers = append([]catalog.Cover(nil), covers...)
		}
	}
	return nil
}

func (r *AlbumSummaryInMemoryRepository) DeleteCoversForAllViewers(ctx context.Context, albumId catalog.AlbumId) error {
	for i := range r.Summaries {
		if r.Summaries[i].AlbumSummary.AlbumId.IsEqual(albumId) {
			r.Summaries[i].AlbumSummary.Covers = nil
		}
	}
	return nil
}

func (r *AlbumSummaryInMemoryRepository) IncrementCountForAllViewers(ctx context.Context, updates []AlbumCountDiff) error {
	if updates == nil {
		return errors.Errorf("IncrementCountForAllViewers(nil): updates should not be nil")
	}

	for _, update := range updates {
		for i := range r.Summaries {
			if r.Summaries[i].AlbumSummary.AlbumId.IsEqual(update.AlbumId) {
				r.Summaries[i].AlbumSummary.MediaCount = r.Summaries[i].AlbumSummary.MediaCount + update.MediaCountDiff
			}
		}
	}

	return nil
}

func (r *AlbumSummaryInMemoryRepository) SetCountForAllViewers(ctx context.Context, updates []AlbumCount) error {
	if updates == nil {
		return errors.Errorf("SetCountForAllViewers(nil): updates should not be nil")
	}

	for _, update := range updates {
		for i := range r.Summaries {
			if r.Summaries[i].AlbumSummary.AlbumId.IsEqual(update.AlbumId) {
				r.Summaries[i].AlbumSummary.MediaCount = update.MediaCount
			}
		}
	}

	return nil
}

func (r *AlbumSummaryInMemoryRepository) RenameAlbum(ctx context.Context, existingId, renamedId catalog.AlbumId, newName string) error {
	ownerIndex := slices.IndexFunc(r.Summaries, func(current UserAlbumSummary) bool {
		return current.AlbumSummary.AlbumId.IsEqual(existingId) && current.Availability.AsOwner
	})
	if ownerIndex < 0 {
		return nil
	}
	source := r.Summaries[ownerIndex].AlbumSummary

	type viewerRow struct {
		availability Availability
		covers       []catalog.Cover
	}
	var viewers []viewerRow
	for _, summary := range r.Summaries {
		if summary.AlbumSummary.AlbumId.IsEqual(existingId) {
			viewers = append(viewers, viewerRow{
				availability: summary.Availability,
				covers:       append([]catalog.Cover(nil), summary.AlbumSummary.Covers...),
			})
		}
	}

	r.Summaries = slices.DeleteFunc(r.Summaries, func(current UserAlbumSummary) bool {
		return current.AlbumSummary.AlbumId.IsEqual(existingId)
	})

	for _, viewer := range viewers {
		r.Summaries = append(r.Summaries, UserAlbumSummary{
			AlbumSummary: AlbumSummary{
				AlbumId:    renamedId,
				Name:       newName,
				Start:      source.Start,
				End:        source.End,
				MediaCount: source.MediaCount,
				Covers:     viewer.covers,
			},
			Availability: viewer.availability,
		})
	}

	return nil
}

func (r *AlbumSummaryInMemoryRepository) DeleteRow(ctx context.Context, availability Availability, albumId catalog.AlbumId) error {
	index := slices.IndexFunc(r.Summaries, func(current UserAlbumSummary) bool {
		return current.AlbumSummary.AlbumId.IsEqual(albumId) && current.Availability == availability
	})
	if index >= 0 {
		r.Summaries = append(r.Summaries[:index], r.Summaries[index+1:]...)
	}

	return nil
}

func (r *AlbumSummaryInMemoryRepository) DeleteAllRowsForAlbum(ctx context.Context, albumId catalog.AlbumId) error {
	r.Summaries = slices.DeleteFunc(r.Summaries, func(current UserAlbumSummary) bool {
		return current.AlbumSummary.AlbumId.IsEqual(albumId)
	})
	return nil
}

func (r *AlbumSummaryInMemoryRepository) DeleteLegacyRowsForUser(_ context.Context, userId usermodel.UserId) error {
	r.LegacyRowsCleanedForUsers = append(r.LegacyRowsCleanedForUsers, userId)
	return nil
}
