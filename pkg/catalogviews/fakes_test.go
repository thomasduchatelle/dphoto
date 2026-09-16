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

func stubFindAlbumsByIdsPort(albums ...*catalog.Album) FindAlbumsByIdsFunc {
	byId := make(map[catalog.AlbumId]*catalog.Album, len(albums))
	for _, album := range albums {
		byId[album.AlbumId] = album
	}
	return func(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error) {
		var result []*catalog.Album
		for _, id := range ids {
			if album, ok := byId[id]; ok {
				result = append(result, album)
			}
		}
		return result, nil
	}
}
