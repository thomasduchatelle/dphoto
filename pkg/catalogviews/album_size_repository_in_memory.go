package catalogviews

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"slices"
)

type UserAlbumSize struct {
	AlbumSize    AlbumSize
	Availability Availability
}

func (s UserAlbumSize) ToMultiUser() MultiUserAlbumSize {
	return MultiUserAlbumSize{AlbumSize: s.AlbumSize, Users: []Availability{s.Availability}}
}

type AlbumSizeInMemoryRepository struct {
	Sizes []UserAlbumSize
}

func (r *AlbumSizeInMemoryRepository) GetAvailabilitiesByUser(ctx context.Context, userId usermodel.UserId) ([]UserAlbumSize, error) {
	var sizes []UserAlbumSize
	for _, size := range r.Sizes {
		if size.Availability.UserId == userId {
			sizes = append(sizes, size)
		}
	}

	return sizes, nil
}

func (r *AlbumSizeInMemoryRepository) GetAlbumSizes(ctx context.Context, userId usermodel.UserId, owner ...ownermodel.Owner) ([]UserAlbumSize, error) {
	var sizes []UserAlbumSize
	for _, size := range r.Sizes {
		if size.Availability.UserId == userId && slices.Contains(owner, size.AlbumSize.AlbumId.Owner) {
			sizes = append(sizes, size)
		}
	}

	return sizes, nil
}

func (r *AlbumSizeInMemoryRepository) InsertAlbumSize(ctx context.Context, albumSize []MultiUserAlbumSize) error {
	if albumSize == nil {
		return errors.Errorf("InsertAlbumSize(nil): albumSize should not be nil")
	}

	for _, size := range albumSize {
		for _, user := range size.Users {
			userAlbumSize := UserAlbumSize{AlbumSize: size.AlbumSize, Availability: user}

			index := slices.IndexFunc(r.Sizes, func(current UserAlbumSize) bool {
				return current.AlbumSize.AlbumId == size.AlbumSize.AlbumId && current.Availability.UserId == user.UserId
			})
			if index >= 0 {
				r.Sizes[index] = userAlbumSize
			} else {
				r.Sizes = append(r.Sizes, userAlbumSize)
			}
		}
	}

	return nil

}

func (r *AlbumSizeInMemoryRepository) UpdateAlbumSize(ctx context.Context, albumCountUpdates []AlbumSizeDiff) error {
	if albumCountUpdates == nil {
		return errors.Errorf("UpdateAlbumSize(nil): albumCountUpdates should not be nil")
	}

	for _, update := range albumCountUpdates {
		for _, user := range update.Users {
			index := slices.IndexFunc(r.Sizes, func(current UserAlbumSize) bool {
				return current.AlbumSize.AlbumId == update.AlbumId && current.Availability.UserId == user.UserId
			})
			if index >= 0 {
				if r.Sizes[index].Availability != user {
					return errors.Errorf("availability cannot be updated during a UpdateAlbumSize (%s != %s)", r.Sizes[index].Availability, user)
				}
				r.Sizes[index].AlbumSize.MediaCount = r.Sizes[index].AlbumSize.MediaCount + update.MediaCountDiff
			} else {
				r.Sizes = append(r.Sizes, UserAlbumSize{AlbumSize: AlbumSize{AlbumId: update.AlbumId, MediaCount: update.MediaCountDiff}, Availability: user})
			}
		}
	}

	return nil
}

func (r *AlbumSizeInMemoryRepository) DeleteAlbumSize(ctx context.Context, availability Availability, albumId catalog.AlbumId) error {
	index := slices.IndexFunc(r.Sizes, func(current UserAlbumSize) bool {
		return current.AlbumSize.AlbumId == albumId && current.Availability == availability
	})
	if index >= 0 {
		r.Sizes = append(r.Sizes[:index], r.Sizes[index+1:]...)
	}

	return nil
}
