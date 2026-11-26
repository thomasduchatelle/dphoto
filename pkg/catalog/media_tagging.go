package catalog

import (
	"context"

	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// TagRepository persists a media's tags and lists medias by tag. The media's #METADATA row
// holds the tag list as the source of truth ; a sibling index item per tag feeds the
// medias-by-tag GSI. Writes are keyed by owner and media id, so the repository rejects a
// media that is not owned by the owner (MediaNotFoundError) and no album lookup is needed.
type TagRepository interface {
	// AddTag attaches the tag to the owner's media and reconciles the index ; it is idempotent
	// (adding twice is a no-op). It returns MediaNotFoundError if the media is not owned.
	AddTag(ctx context.Context, owner ownermodel.Owner, mediaId MediaId, tag Tag) error
	// RemoveTag detaches the tag from the owner's media and reconciles the index ; it is
	// idempotent (removing an absent tag succeeds).
	RemoveTag(ctx context.Context, owner ownermodel.Owner, mediaId MediaId, tag Tag) error
	// FindMediasByTag returns every media carrying the tag for the owner, across albums.
	FindMediasByTag(ctx context.Context, owner ownermodel.Owner, tag Tag) ([]*MediaWithAlbum, error)
}

// TagService adds and removes tags on medias and lists medias by tag. Writes are owner-only:
// the repository stores tags on the media partition keyed by owner, so tagging a media that
// is not owned by the owner is rejected.
type TagService struct {
	TagRepository TagRepository
}

func NewTagService(tagRepository TagRepository) *TagService {
	return &TagService{
		TagRepository: tagRepository,
	}
}

// AddTag attaches the tag to the media owned by owner. It is idempotent.
func (s *TagService) AddTag(ctx context.Context, owner ownermodel.Owner, mediaId MediaId, tag Tag) error {
	return s.TagRepository.AddTag(ctx, owner, mediaId, tag)
}

// RemoveTag detaches the tag from the media owned by owner. It is idempotent.
func (s *TagService) RemoveTag(ctx context.Context, owner ownermodel.Owner, mediaId MediaId, tag Tag) error {
	return s.TagRepository.RemoveTag(ctx, owner, mediaId, tag)
}

// ListMediasByTag returns every media of the owner carrying the tag, across all albums.
func (s *TagService) ListMediasByTag(ctx context.Context, owner ownermodel.Owner, tag Tag) ([]*MediaWithAlbum, error) {
	return s.TagRepository.FindMediasByTag(ctx, owner, tag)
}
