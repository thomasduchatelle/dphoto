package catalogviews

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type MediaCounterPortFake map[catalog.AlbumId]int

func (m MediaCounterPortFake) CountMedia(ctx context.Context, album ...catalog.AlbumId) (map[catalog.AlbumId]int, error) {
	counts := make(map[catalog.AlbumId]int)
	for _, a := range album {
		counts[a] = m[a]
	}
	return counts, nil
}

func stubFindAlbumByOwnerPort(albums ...*catalog.Album) FindAlbumByOwnerFunc {
	return func(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
		return albums, nil
	}
}

func stubOwnerUserIdPort(owner ownermodel.Owner, userId usermodel.UserId) OwnerUserIdFunc {
	return func(ctx context.Context, o ownermodel.Owner) (usermodel.UserId, error) {
		if o == owner {
			return userId, nil
		}
		return "", nil
	}
}

type ListUserWhoCanAccessAlbumPortFake struct {
	Values map[catalog.AlbumId][]Availability
}

func (l *ListUserWhoCanAccessAlbumPortFake) ListUsersWhoCanAccessAlbum(ctx context.Context, albumId ...catalog.AlbumId) (map[catalog.AlbumId][]Availability, error) {
	if albumId == nil {
		return nil, errors.Errorf("ListUsersWhoCanAccessAlbum(nil): albumId should not be nil")
	}

	result := make(map[catalog.AlbumId][]Availability)
	for _, id := range albumId {
		result[id] = l.Values[id]
	}
	return result, nil
}
