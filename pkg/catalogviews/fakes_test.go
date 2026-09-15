package catalogviews

import (
	"context"

	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
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
