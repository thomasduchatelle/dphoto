package catalog

import (
	"context"
	"math/rand"
	"slices"

	"github.com/pkg/errors"
)

type CoverServicePort interface {
	ForcedRandomise(ctx context.Context, albumIds ...AlbumId) (map[AlbumId][]Cover, error)
	StableRandomise(ctx context.Context, albumIds ...AlbumId) (map[AlbumId][]Cover, error)
	ApplyTransfer(ctx context.Context, transferred TransferredMedias, deletedAlbumIds ...AlbumId) (map[AlbumId][]Cover, error)
}

// CoverRepository persists an album's cover set as a single record.
type CoverRepository interface {
	FindCoversByAlbum(ctx context.Context, albumId AlbumId) ([]Cover, error)
	FindCoversByAlbums(ctx context.Context, albumIds ...AlbumId) (map[AlbumId][]Cover, error)
	SaveCovers(ctx context.Context, albumId AlbumId, covers []Cover) error
	MoveCovers(ctx context.Context, from, to AlbumId) error
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

type CoverService struct {
	CoverRepository     CoverRepository
	MediaReadRepository MediaReadRepository
	Randomiser          Randomiser
}

func NewCoverService(coverRepository CoverRepository, mediaReadRepository MediaReadRepository) *CoverService {
	return &CoverService{
		CoverRepository:     coverRepository,
		MediaReadRepository: mediaReadRepository,
		Randomiser:          DefaultRandomiser,
	}
}

// ForcedRandomise refreshes the cover set of every requested album by stripping covers
// whose media is no longer in the album, dropping every RANDOM cover so that newly-added
// medias get a chance to appear as a cover, and filling the empty slots at random from
// the album's full image set. CHERRY_PICKED covers whose media is still in the album are
// always kept. The returned map carries one entry per album whose cover set actually
// changed; unchanged albums are omitted. Returns a nil map when nothing changed.
func (c *CoverService) ForcedRandomise(ctx context.Context, albumIds ...AlbumId) (map[AlbumId][]Cover, error) {
	return c.randomise(ctx, func(cover Cover) bool {
		return cover.Origin == CoverOriginRandom
	}, albumIds)
}

// StableRandomise refreshes the cover set of every requested album by stripping covers
// whose media is no longer in the album and filling the empty slots at random from the
// album's full image set. Every cover whose media is still in the album is kept
// untouched. The returned map carries one entry per album whose cover set actually
// changed; unchanged albums are omitted. Returns a nil map when nothing changed.
func (c *CoverService) StableRandomise(ctx context.Context, albumIds ...AlbumId) (map[AlbumId][]Cover, error) {
	return c.randomise(ctx, func(Cover) bool { return false }, albumIds)
}

// randomise walks every requested album, drops the covers whose media is gone AND the
// ones for which `dropExisting` returns true, then fills the empty slots at random from
// the album's full image set. Only albums whose cover set actually changed are returned
// and persisted.
func (c *CoverService) randomise(ctx context.Context, dropExisting func(Cover) bool, albumIds []AlbumId) (map[AlbumId][]Cover, error) {
	if len(albumIds) == 0 {
		return nil, nil
	}

	coversByAlbumId, err := c.CoverRepository.FindCoversByAlbums(ctx, albumIds...)
	if err != nil {
		return nil, errors.Wrapf(err, "randomise failed to load covers for albums %v", albumIds)
	}

	var changed map[AlbumId][]Cover
	for _, albumId := range albumIds {
		original := coversByAlbumId[albumId]

		medias, err := c.MediaReadRepository.FindMedias(ctx, NewFindMediaRequest(albumId.Owner).WithAlbum(albumId.FolderName))
		if err != nil {
			return nil, errors.Wrapf(err, "randomise failed to list medias for %s", albumId)
		}

		kept := slices.DeleteFunc(slices.Clone(original), func(cover Cover) bool {
			stillInAlbum := slices.ContainsFunc(medias, func(meta *MediaMeta) bool {
				return meta.Id == cover.MediaId
			})
			if !stillInAlbum {
				return true
			}
			return dropExisting(cover)
		})

		filled, _ := c.fillCoversWithMedias(kept, medias)
		if coversEqual(original, filled) {
			continue
		}
		if changed == nil {
			changed = make(map[AlbumId][]Cover)
		}
		changed[albumId] = filled
		if err := c.CoverRepository.SaveCovers(ctx, albumId, filled); err != nil {
			return nil, errors.Wrapf(err, "randomise failed to save covers for %s", albumId)
		}
	}

	return changed, nil
}

func coversEqual(a, b []Cover) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type affectedAlbum struct {
	albumId     AlbumId
	removed     map[MediaId]interface{}
	added       bool
	addedMedias []MediaId
}

func affectedAlbums(transferred TransferredMedias) []affectedAlbum {
	var affected []affectedAlbum

	strikeThroughAlbums := make(map[AlbumId]interface{})
	for _, sourceAlbumId := range transferred.FromAlbums {
		added := false
		medias := make(map[MediaId]interface{})
		var transferredMedias []MediaId
		for targetAlbumId, addedMedias := range transferred.Transfers {
			if sourceAlbumId == targetAlbumId {
				added = true
				transferredMedias = addedMedias
				strikeThroughAlbums[targetAlbumId] = nil
			} else {
				for _, media := range addedMedias {
					medias[media] = nil
				}
			}
		}

		affected = append(affected, affectedAlbum{
			albumId:     sourceAlbumId,
			removed:     medias,
			added:       added,
			addedMedias: transferredMedias,
		})
	}

	for targetAlbumId, addedMedias := range transferred.Transfers {
		if _, present := strikeThroughAlbums[targetAlbumId]; !present {
			affected = append(affected, affectedAlbum{
				albumId:     targetAlbumId,
				removed:     nil,
				added:       true,
				addedMedias: addedMedias,
			})
		}
	}

	return affected
}

// ApplyTransfer processes every album affected by `transferred`, plus every album
// listed in `deletedAlbumIds` whose canonical cover record must be erased even when no
// media was moved out of it (e.g. empty album deleted):
//
//   - For each source album, covers whose media has moved away are stripped.
//   - CHERRY_PICKED covers that moved to a destination album are inherited by the
//     destination (as CHERRY_PICKED), displacing the oldest RANDOM cover when needed.
//     Destinations already full of CHERRY_PICKED covers silently ignore the incoming
//     one.
//   - Every affected album has its empty slots backfilled from its full image set.
//   - Every album in `deletedAlbumIds` has its canonical cover record erased
//     unconditionally.
//
// Returns one entry per album whose cover set actually changed; deleted albums map to
// nil. Returns a nil map when no album changed.
func (c *CoverService) ApplyTransfer(ctx context.Context, transferred TransferredMedias, deletedAlbumIds ...AlbumId) (map[AlbumId][]Cover, error) {
	albums := affectedAlbums(transferred)

	affected := make([]AlbumId, len(albums), len(albums))
	for i, album := range albums {
		affected[i] = album.albumId
	}

	coversByAlbum, err := c.CoverRepository.FindCoversByAlbums(ctx, affected...)
	if err != nil {
		return nil, errors.Wrapf(err, "ApplyTransfer failed to load covers for albums %v", affected)
	}

	pickedMedias := make(map[MediaId]Cover)
	for _, covers := range coversByAlbum {
		for _, cover := range covers {
			if cover.Origin == CoverOriginCherryPicked {
				pickedMedias[cover.MediaId] = cover
			}
		}
	}

	var updatedCovers map[AlbumId][]Cover
	for _, album := range albums {
		covers, _ := coversByAlbum[album.albumId]

		covers, hasBeenFiltered := c.filterOutRemovedMedias(covers, album)

		hasBeenAltered, covers := c.forceAlreadyCherryPickedCovers(covers, album, transferred, pickedMedias)

		hasBeenFilled := false
		if hasBeenFiltered || len(album.addedMedias) > 0 {
			covers, hasBeenFilled, err = c.fillCovers(ctx, covers, album.albumId)
			if err != nil {
				return nil, err
			}
		}

		if hasBeenFiltered || hasBeenAltered || hasBeenFilled {
			if updatedCovers == nil {
				updatedCovers = make(map[AlbumId][]Cover)
			}
			updatedCovers[album.albumId] = covers
			err = c.CoverRepository.SaveCovers(ctx, album.albumId, covers)
			if err != nil {
				return nil, errors.Wrapf(err, "ApplyTransfer failed to store new covers for %s", album.albumId)
			}
		}
	}

	for _, deletedAlbumId := range deletedAlbumIds {
		if err := c.CoverRepository.SaveCovers(ctx, deletedAlbumId, nil); err != nil {
			return nil, errors.Wrapf(err, "ApplyTransfer failed to delete covers for %s", deletedAlbumId)
		}
		if updatedCovers == nil {
			updatedCovers = make(map[AlbumId][]Cover)
		}
		updatedCovers[deletedAlbumId] = nil
	}

	return updatedCovers, nil
}

func (c *CoverService) filterOutRemovedMedias(covers []Cover, album affectedAlbum) ([]Cover, bool) {
	size := len(covers)

	covers = slices.DeleteFunc(covers, func(cover Cover) bool {
		_, removed := album.removed[cover.MediaId]
		return removed
	})
	hasBeenFiltered := size != len(covers)

	if len(covers) == 0 {
		covers = nil
	}
	return covers, hasBeenFiltered
}

func (c *CoverService) forceAlreadyCherryPickedCovers(covers []Cover, album affectedAlbum, _ TransferredMedias, pickedMedias map[MediaId]Cover) (bool, []Cover) {
	hasBeenAltered := false
	for _, addedMedia := range album.addedMedias {
		if sourceCover, picked := pickedMedias[addedMedia]; picked {
			if len(covers) < MaxCoversPerAlbum {
				hasBeenAltered = true
				covers = append(covers, sourceCover)
			} else {
				for i, cover := range covers {
					if cover.Origin != CoverOriginCherryPicked {
						hasBeenAltered = true
						covers[i] = sourceCover
						break
					}
				}
			}
		}
	}
	return hasBeenAltered, covers
}

func (c *CoverService) fillCovers(ctx context.Context, covers []Cover, albumId AlbumId) ([]Cover, bool, error) {
	hasBeenFilled := false
	if len(covers) < MaxCoversPerAlbum {
		// Refill the covers if there is a chance the albums has more medias (i.e. if it wasn't full and no medias got added, no point trying to find new ones)
		medias, err := c.MediaReadRepository.FindMedias(ctx, NewFindMediaRequest(albumId.Owner).WithAlbum(albumId.FolderName))
		if err != nil {
			return nil, false, errors.Wrapf(err, "ApplyTransfer failed to list medias for %s", albumId)
		}

		covers, hasBeenFilled = c.fillCoversWithMedias(covers, medias)
	}
	return covers, hasBeenFilled, nil
}

func (c *CoverService) fillCoversWithMedias(covers []Cover, medias []*MediaMeta) ([]Cover, bool) {
	medias = slices.DeleteFunc(medias, func(meta *MediaMeta) bool {
		return meta.Type != MediaTypeImage || slices.ContainsFunc(covers, func(cover Cover) bool {
			return cover.MediaId == meta.Id
		})
	})

	hasBeenFilled := false
	for _, indice := range c.Randomiser.SampleIndices(len(medias), MaxCoversPerAlbum-len(covers)) {
		hasBeenFilled = true
		covers = append(covers, Cover{
			MediaId:  medias[indice].Id,
			Filename: medias[indice].Filename,
			Origin:   CoverOriginRandom,
		})
	}
	return covers, hasBeenFilled
}
