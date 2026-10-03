package catalog

import (
	"context"
	"math/rand"

	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// CoverRepository persists an album's cover set as a single record.
type CoverRepository interface {
	FindCoversByAlbum(ctx context.Context, albumId AlbumId) ([]Cover, error)
	SaveCovers(ctx context.Context, albumId AlbumId, covers []Cover) error
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

// CoverMaintenance reconciles the cover set of an album, picking randomly from the
// album's full image set. It offers two strategies:
//
//   - Refresh: drop every RANDOM cover, keep every CHERRY_PICKED cover, strip any kept
//     cover whose MediaId is in `removed`, fill empty slots up to 4. Use when the cover
//     set should reflect what is newest or most interesting (media insertion, admin
//     backfill, owner-triggered re-randomise).
//   - Stabilise: keep every existing cover (both origins) except those in `removed`,
//     fill empty slots up to 4. Use when stability is preferred over freshness (album
//     creation with transferred medias, dates amended, album deleted neighbours).
//
// Both strategies cap the resulting set at MaxCoversPerAlbum, select only MediaType
// IMAGE candidates, and never re-pick a media that is already a cover. Both methods
// return the resulting set and a `changed` boolean set to true when the cover set
// actually differs from the one loaded from the canonical record.
type CoverMaintenance struct {
	CoverRepository     CoverRepository
	MediaReadRepository MediaReadRepository
	Randomiser          Randomiser
}

func NewCoverMaintenance(coverRepository CoverRepository, mediaReadRepository MediaReadRepository) *CoverMaintenance {
	return &CoverMaintenance{
		CoverRepository:     coverRepository,
		MediaReadRepository: mediaReadRepository,
		Randomiser:          DefaultRandomiser,
	}
}

// Refresh applies the Refresh strategy: see CoverMaintenance.
func (c *CoverMaintenance) Refresh(ctx context.Context, albumId AlbumId, removed []MediaId) ([]Cover, bool, error) {
	return c.reconcile(ctx, albumId, removed, true)
}

// Stabilise applies the Stabilise strategy: see CoverMaintenance.
func (c *CoverMaintenance) Stabilise(ctx context.Context, albumId AlbumId, removed []MediaId) ([]Cover, bool, error) {
	return c.reconcile(ctx, albumId, removed, false)
}

func (c *CoverMaintenance) reconcile(ctx context.Context, albumId AlbumId, removed []MediaId, dropRandom bool) ([]Cover, bool, error) {
	existing, err := c.CoverRepository.FindCoversByAlbum(ctx, albumId)
	if err != nil {
		return nil, false, errors.Wrapf(err, "cover reconciliation of %s failed to load current covers", albumId)
	}

	removedSet := make(map[MediaId]bool, len(removed))
	for _, id := range removed {
		removedSet[id] = true
	}

	var kept []Cover
	for _, cover := range existing {
		if dropRandom && cover.Origin == CoverOriginRandom {
			continue
		}
		if removedSet[cover.MediaId] {
			continue
		}
		kept = append(kept, cover)
	}

	slotsToFill := MaxCoversPerAlbum - len(kept)
	result := kept
	if slotsToFill > 0 {
		alreadyCovered := make(map[MediaId]bool, len(kept))
		for _, cover := range kept {
			alreadyCovered[cover.MediaId] = true
		}

		medias, err := c.MediaReadRepository.FindMedias(ctx, NewFindMediaRequest(albumId.Owner).WithAlbum(albumId.FolderName))
		if err != nil {
			return nil, false, errors.Wrapf(err, "cover reconciliation of %s failed to list medias", albumId)
		}

		var eligible []*MediaMeta
		for _, candidate := range medias {
			if candidate == nil || candidate.Type != MediaTypeImage {
				continue
			}
			if alreadyCovered[candidate.Id] {
				continue
			}
			eligible = append(eligible, candidate)
		}

		if len(eligible) > 0 {
			picks := slotsToFill
			if picks > len(eligible) {
				picks = len(eligible)
			}
			indices := c.Randomiser.SampleIndices(len(eligible), picks)
			for _, idx := range indices {
				media := eligible[idx]
				result = append(result, Cover{MediaId: media.Id, Filename: media.Filename, Origin: CoverOriginRandom})
			}
		}
	}

	if coversEqual(existing, result) {
		return existing, false, nil
	}

	if err := c.CoverRepository.SaveCovers(ctx, albumId, result); err != nil {
		return nil, false, errors.Wrapf(err, "cover reconciliation of %s failed to save covers", albumId)
	}
	return result, true, nil
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

type FindAlbumByOwnerPort interface {
	FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*Album, error)
}

// RefreshCoversPort is the slice of CoverMaintenance exposed to callers that only
// trigger a Refresh (admin backfill, Phase 3 owner-triggered re-randomise, and the
// InsertMedias use case).
type RefreshCoversPort interface {
	Refresh(ctx context.Context, albumId AlbumId, removed []MediaId) ([]Cover, bool, error)
}

// StabiliseCoversPort is the slice of CoverMaintenance exposed to callers that only
// trigger a Stabilise (AlbumCreated, AlbumDatesAmended, AlbumDeleted use cases).
type StabiliseCoversPort interface {
	Stabilise(ctx context.Context, albumId AlbumId, removed []MediaId) ([]Cover, bool, error)
}

// BackfillCoversViewUpdater propagates the new covers of a single album to the
// album-list view. BackfillCovers calls it after every Refresh that actually
// changed the cover set.
type BackfillCoversViewUpdater interface {
	UpdateCovers(ctx context.Context, albumId AlbumId, covers []Cover) error
}

type BackfillCoversViewUpdaterFunc func(ctx context.Context, albumId AlbumId, covers []Cover) error

func (f BackfillCoversViewUpdaterFunc) UpdateCovers(ctx context.Context, albumId AlbumId, covers []Cover) error {
	return f(ctx, albumId, covers)
}

// BackfillCovers reconciles the cover set of every album of an owner and
// propagates the resulting set to the album-list view. Used by the
// administrative CLI to seed existing albums that pre-date the covers feature.
// A failure on a single album is reported and the sweep continues: Refresh
// preserves CHERRY_PICKED covers and caps the set at MaxCoversPerAlbum, so
// re-running the backfill is safe (though RANDOM covers may be redrawn on each
// pass).
type BackfillCovers struct {
	FindAlbumByOwnerPort      FindAlbumByOwnerPort
	RefreshCoversPort         RefreshCoversPort
	BackfillCoversViewUpdater BackfillCoversViewUpdater
}

type BackfillReport struct {
	Albums   int
	Failures []BackfillFailure
}

type BackfillFailure struct {
	AlbumId AlbumId
	Err     error
}

func (b *BackfillCovers) BackfillForOwner(ctx context.Context, owner ownermodel.Owner) (BackfillReport, error) {
	albums, err := b.FindAlbumByOwnerPort.FindAlbumsByOwner(ctx, owner)
	if err != nil {
		return BackfillReport{}, errors.Wrapf(err, "BackfillCovers(%s) failed to list albums", owner)
	}

	report := BackfillReport{Albums: len(albums)}
	for _, album := range albums {
		covers, changed, err := b.RefreshCoversPort.Refresh(ctx, album.AlbumId, nil)
		if err != nil {
			report.Failures = append(report.Failures, BackfillFailure{AlbumId: album.AlbumId, Err: err})
			continue
		}
		if !changed {
			continue
		}
		if err := b.BackfillCoversViewUpdater.UpdateCovers(ctx, album.AlbumId, covers); err != nil {
			report.Failures = append(report.Failures, BackfillFailure{AlbumId: album.AlbumId, Err: err})
		}
	}
	return report, nil
}
