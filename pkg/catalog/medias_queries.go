package catalog

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

type MediaReadRepository interface {
	FindMedias(ctx context.Context, request *FindMediaRequest) (medias []*MediaMeta, err error)
	FindMediaCurrentAlbum(ctx context.Context, owner ownermodel.Owner, mediaId MediaId) (id *AlbumId, err error)
}

type MediaQueries struct {
	MediaReadRepository MediaReadRepository
}

func (q *MediaQueries) ListMedias(ctx context.Context, albumId AlbumId) ([]*MediaMeta, error) {
	return q.MediaReadRepository.FindMedias(ctx, NewFindMediaRequest(albumId.Owner).WithAlbum(albumId.FolderName))
}

// FindMediaOwnership returns the folderName containing the media, or AlbumNotFoundError.
func (q *MediaQueries) FindMediaOwnership(ctx context.Context, owner ownermodel.Owner, mediaId MediaId) (*AlbumId, error) {
	return q.MediaReadRepository.FindMediaCurrentAlbum(ctx, owner, mediaId)
}
