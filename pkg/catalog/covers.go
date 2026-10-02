package catalog

import (
	"context"
	"math/rand"

	"github.com/pkg/errors"
)

// CoverRepository persists an album's cover set as a single record.
type CoverRepository interface {
	FindCoversByAlbum(ctx context.Context, albumId AlbumId) ([]Cover, error)
	SaveCovers(ctx context.Context, albumId AlbumId, covers []Cover) error
}

// CompleteCovers fills the empty slots of an album's cover set by drawing at random from
// its eligible medias (MediaType IMAGE, not already a cover). See CompleteCoversService.
type CompleteCovers struct {
	CoverRepository     CoverRepository
	MediaReadRepository MediaReadRepository
	Randomiser          Randomiser
}

// Randomiser picks n indices in [0, upperBound) uniformly at random, without replacement.
// The default implementation uses math/rand; tests substitute a deterministic one.
type Randomiser interface {
	SampleIndices(upperBound, n int) []int
}

type RandomiserFunc func(upperBound, n int) []int

func (f RandomiserFunc) SampleIndices(upperBound, n int) []int {
	return f(upperBound, n)
}

// DefaultRandomiser is a Randomiser backed by the global math/rand source. It returns n
// distinct indices in [0, upperBound), or all of them when n >= upperBound.
var DefaultRandomiser Randomiser = RandomiserFunc(func(upperBound, n int) []int {
	if n >= upperBound {
		indices := make([]int, upperBound)
		for i := range indices {
			indices[i] = i
		}
		return indices
	}
	return rand.Perm(upperBound)[:n]
})

func NewCompleteCovers(coverRepository CoverRepository, mediaReadRepository MediaReadRepository) *CompleteCovers {
	return &CompleteCovers{
		CoverRepository:     coverRepository,
		MediaReadRepository: mediaReadRepository,
		Randomiser:          DefaultRandomiser,
	}
}

// CompleteCovers queries the album's medias and fills any empty cover slot with
// uniform-random RANDOM covers. No-op when the set is already full.
func (c *CompleteCovers) CompleteCovers(ctx context.Context, albumId AlbumId) error {
	medias, err := c.MediaReadRepository.FindMedias(ctx, NewFindMediaRequest(albumId.Owner).WithAlbum(albumId.FolderName))
	if err != nil {
		return errors.Wrapf(err, "CompleteCovers(%s) failed to list medias", albumId)
	}
	return c.CompleteCoversFromCandidates(ctx, albumId, medias)
}

// CompleteCoversFromCandidates fills the empty slots of the album's cover set with
// RANDOM covers drawn uniformly at random from the supplied candidates. Only candidates
// of MediaType IMAGE that are not already covers are eligible. No-op when the set is
// already full, or when no eligible candidate is available.
func (c *CompleteCovers) CompleteCoversFromCandidates(ctx context.Context, albumId AlbumId, candidates []*MediaMeta) error {
	existing, err := c.CoverRepository.FindCoversByAlbum(ctx, albumId)
	if err != nil {
		return errors.Wrapf(err, "CompleteCovers(%s) failed to load current covers", albumId)
	}

	slotsToFill := MaxCoversPerAlbum - len(existing)
	if slotsToFill <= 0 {
		return nil
	}

	alreadyCovered := make(map[MediaId]bool, len(existing))
	for _, cover := range existing {
		alreadyCovered[cover.MediaId] = true
	}

	var eligible []*MediaMeta
	for _, candidate := range candidates {
		if candidate == nil || candidate.Type != MediaTypeImage {
			continue
		}
		if alreadyCovered[candidate.Id] {
			continue
		}
		eligible = append(eligible, candidate)
	}

	if len(eligible) == 0 {
		return nil
	}

	picks := slotsToFill
	if picks > len(eligible) {
		picks = len(eligible)
	}

	indices := c.Randomiser.SampleIndices(len(eligible), picks)
	completed := append([]Cover(nil), existing...)
	for _, idx := range indices {
		media := eligible[idx]
		completed = append(completed, Cover{
			MediaId:  media.Id,
			Filename: media.Filename,
			Origin:   CoverOriginRandom,
		})
	}

	return c.CoverRepository.SaveCovers(ctx, albumId, completed)
}
