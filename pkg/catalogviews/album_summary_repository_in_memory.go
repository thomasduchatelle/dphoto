package catalogviews

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"slices"
	"time"
)

// UserAlbumSummary is the read-projection: a summary as seen from a specific user's row.
type UserAlbumSummary struct {
	AlbumSummary AlbumSummary
	Availability Availability
}

func (s UserAlbumSummary) ToSummaryForUsers() AlbumSummaryForUsers {
	return AlbumSummaryForUsers{AlbumSummary: s.AlbumSummary, Users: []Availability{s.Availability}}
}

type AlbumSummaryInMemoryRepository struct {
	Summaries []UserAlbumSummary
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
			userSummary := UserAlbumSummary{AlbumSummary: summary.AlbumSummary, Availability: user}

			index := slices.IndexFunc(r.Summaries, func(current UserAlbumSummary) bool {
				return current.AlbumSummary.AlbumId == summary.AlbumId && current.Availability.UserId == user.UserId
			})
			if index >= 0 {
				r.Summaries[index] = userSummary
			} else {
				r.Summaries = append(r.Summaries, userSummary)
			}
		}
	}

	return nil
}

func (r *AlbumSummaryInMemoryRepository) SetDisplayFields(ctx context.Context, albumId catalog.AlbumId, users []Availability, name string, start, end time.Time) error {
	for _, user := range users {
		index := slices.IndexFunc(r.Summaries, func(current UserAlbumSummary) bool {
			return current.AlbumSummary.AlbumId == albumId && current.Availability.UserId == user.UserId
		})
		if index >= 0 {
			r.Summaries[index].AlbumSummary.Name = name
			r.Summaries[index].AlbumSummary.Start = start
			r.Summaries[index].AlbumSummary.End = end
		} else {
			r.Summaries = append(r.Summaries, UserAlbumSummary{
				AlbumSummary: AlbumSummary{AlbumId: albumId, Name: name, Start: start, End: end},
				Availability: user,
			})
		}
	}

	return nil
}

func (r *AlbumSummaryInMemoryRepository) IncrementCounts(ctx context.Context, updates []AlbumMediaCountDiff) error {
	if updates == nil {
		return errors.Errorf("IncrementCounts(nil): updates should not be nil")
	}

	for _, update := range updates {
		for _, user := range update.Users {
			index := slices.IndexFunc(r.Summaries, func(current UserAlbumSummary) bool {
				return current.AlbumSummary.AlbumId == update.AlbumId && current.Availability.UserId == user.UserId
			})
			if index >= 0 {
				if r.Summaries[index].Availability != user {
					return errors.Errorf("availability cannot be updated during a IncrementCounts (%s != %s)", r.Summaries[index].Availability, user)
				}
				r.Summaries[index].AlbumSummary.MediaCount = r.Summaries[index].AlbumSummary.MediaCount + update.MediaCountDiff
			} else {
				r.Summaries = append(r.Summaries, UserAlbumSummary{
					AlbumSummary: AlbumSummary{AlbumId: update.AlbumId, MediaCount: update.MediaCountDiff},
					Availability: user,
				})
			}
		}
	}

	return nil
}

func (r *AlbumSummaryInMemoryRepository) SetCounts(ctx context.Context, updates []AlbumMediaCountForUsers) error {
	if updates == nil {
		return errors.Errorf("SetCounts(nil): updates should not be nil")
	}

	for _, update := range updates {
		for _, user := range update.Users {
			index := slices.IndexFunc(r.Summaries, func(current UserAlbumSummary) bool {
				return current.AlbumSummary.AlbumId == update.AlbumId && current.Availability.UserId == user.UserId
			})
			if index >= 0 {
				if r.Summaries[index].Availability != user {
					return errors.Errorf("availability cannot be updated during a SetCounts (%s != %s)", r.Summaries[index].Availability, user)
				}
				r.Summaries[index].AlbumSummary.MediaCount = update.MediaCount
			} else {
				r.Summaries = append(r.Summaries, UserAlbumSummary{
					AlbumSummary: AlbumSummary{AlbumId: update.AlbumId, MediaCount: update.MediaCount},
					Availability: user,
				})
			}
		}
	}

	return nil
}

func (r *AlbumSummaryInMemoryRepository) DeleteRow(ctx context.Context, availability Availability, albumId catalog.AlbumId) error {
	index := slices.IndexFunc(r.Summaries, func(current UserAlbumSummary) bool {
		return current.AlbumSummary.AlbumId == albumId && current.Availability == availability
	})
	if index >= 0 {
		r.Summaries = append(r.Summaries[:index], r.Summaries[index+1:]...)
	}

	return nil
}

func (r *AlbumSummaryInMemoryRepository) DeleteAllRowsForAlbum(ctx context.Context, albumId catalog.AlbumId) error {
	r.Summaries = slices.DeleteFunc(r.Summaries, func(current UserAlbumSummary) bool {
		return current.AlbumSummary.AlbumId == albumId
	})
	return nil
}
