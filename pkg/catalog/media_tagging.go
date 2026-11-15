package catalog

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// TagRepository persists tag items as siblings of the media and lists medias by tag. A tag
// item is stored on the media partition and carries the AlbumId, so adding or removing a
// tag is a single idempotent write/delete and listing medias of a tag resolves their album
// without an extra lookup.
type TagRepository interface {
	// AddTag stores the (media, tag) item ; it is idempotent (adding twice is a no-op).
	AddTag(ctx context.Context, albumId AlbumId, mediaId MediaId, tag Tag) error
	// RemoveTag deletes the (media, tag) item ; it is idempotent (removing an absent tag succeeds).
	RemoveTag(ctx context.Context, owner ownermodel.Owner, mediaId MediaId, tag Tag) error
	// FindMediasByTag returns every media carrying the tag for the owner, across albums.
	FindMediasByTag(ctx context.Context, owner ownermodel.Owner, tag Tag) ([]*MediaWithAlbum, error)
}

// TagService adds and removes tags on medias and lists medias by tag. Writes are owner-only:
// the media must belong to the owner, which is proven by resolving its current album.
type TagService struct {
	MediaReadRepository MediaReadRepository
	TagRepository       TagRepository
}

func NewTagService(mediaReadRepository MediaReadRepository, tagRepository TagRepository) *TagService {
	return &TagService{
		MediaReadRepository: mediaReadRepository,
		TagRepository:       tagRepository,
	}
}

// AddTag attaches the tag to the media owned by owner. It is idempotent.
func (s *TagService) AddTag(ctx context.Context, owner ownermodel.Owner, mediaId MediaId, tag Tag) error {
	albumId, err := s.MediaReadRepository.FindMediaCurrentAlbum(ctx, owner, mediaId)
	if err != nil {
		return errors.Wrapf(err, "cannot add tag %s to media %s/%s", tag, owner, mediaId)
	}

	return s.TagRepository.AddTag(ctx, *albumId, mediaId, tag)
}

// RemoveTag detaches the tag from the media owned by owner. It is idempotent.
func (s *TagService) RemoveTag(ctx context.Context, owner ownermodel.Owner, mediaId MediaId, tag Tag) error {
	return s.TagRepository.RemoveTag(ctx, owner, mediaId, tag)
}

// ListMediasByTag returns every media of the owner carrying the tag, across all albums.
func (s *TagService) ListMediasByTag(ctx context.Context, owner ownermodel.Owner, tag Tag) ([]*MediaWithAlbum, error) {
	return s.TagRepository.FindMediasByTag(ctx, owner, tag)
}
