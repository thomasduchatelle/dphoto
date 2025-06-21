package catalog

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

type InMemoryMedia struct {
	MediaMeta
	AlbumId AlbumId
}

func NewInMemoryMedia(id MediaId, albumId AlbumId) InMemoryMedia {
	return InMemoryMedia{
		MediaMeta: MediaMeta{
			Id: id,
		},
		AlbumId: albumId,
	}
}

// MediaQueriesInMemory is a copy of MediaQueries with an in-memory implementation
type MediaQueriesInMemory struct {
	Medias []InMemoryMedia
}

func (q *MediaQueriesInMemory) ListMedias(ctx context.Context, albumId AlbumId) ([]*MediaMeta, error) {
	var medias []*MediaMeta
	for _, media := range q.Medias {
		if media.AlbumId == albumId {
			medias = append(medias, &media.MediaMeta)
		}
	}

	return medias, nil
}

// FindMediaOwnership returns the folderName containing the media, or AlbumNotFoundErr.
func (q *MediaQueriesInMemory) FindMediaOwnership(ctx context.Context, owner ownermodel.Owner, mediaId MediaId) (*AlbumId, error) {
	for _, media := range q.Medias {
		if media.Id == mediaId {
			return &media.AlbumId, nil
		}
	}

	return nil, MediaNotFoundError
}
