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

// CoversChangedObserver is notified when the cover set of an album has changed.
// Empty slice means "no covers". Observers are called only when the set actually
// differs from the previous one.
type CoversChangedObserver interface {
	OnCoversChanged(ctx context.Context, albumId AlbumId, covers []Cover) error
}

type CoversChangedObserverFunc func(ctx context.Context, albumId AlbumId, covers []Cover) error

func (f CoversChangedObserverFunc) OnCoversChanged(ctx context.Context, albumId AlbumId, covers []Cover) error {
	return f(ctx, albumId, covers)
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

// CoverMaintenance applies the cover invariant on an album: it keeps CHERRY_PICKED
// covers, drops every RANDOM cover, strips any kept cover whose MediaId is in the
// `removed` set, and fills the empty slots up to MaxCoversPerAlbum from the supplied
// `added` candidates (falling back to a query against the album's medias when the
// `added` set does not provide enough eligible images).
//
// When the resulting set differs from the one loaded from the repository, every
// registered CoversChangedObserver is notified.
type CoverMaintenance struct {
	CoverRepository     CoverRepository
	MediaReadRepository MediaReadRepository
	Randomiser          Randomiser
	Observers           []CoversChangedObserver
}

func NewCoverMaintenance(coverRepository CoverRepository, mediaReadRepository MediaReadRepository, observers ...CoversChangedObserver) *CoverMaintenance {
	return &CoverMaintenance{
		CoverRepository:     coverRepository,
		MediaReadRepository: mediaReadRepository,
		Randomiser:          DefaultRandomiser,
		Observers:           observers,
	}
}

// Reconcile applies the cover invariant on an album.
//
// `added` are candidates to fill the empty slots (typically medias just transferred
// into the album). `removed` is the set of MediaIds that must no longer appear as
// covers (typically medias just transferred out of the album). Both are optional:
// an empty `added` triggers the fallback query when slots need to be filled; an
// empty `removed` leaves kept covers untouched.
func (c *CoverMaintenance) Reconcile(ctx context.Context, albumId AlbumId, added []*MediaMeta, removed []MediaId) error {
	existing, err := c.CoverRepository.FindCoversByAlbum(ctx, albumId)
	if err != nil {
		return errors.Wrapf(err, "Reconcile(%s) failed to load current covers", albumId)
	}

	removedSet := make(map[MediaId]bool, len(removed))
	for _, id := range removed {
		removedSet[id] = true
	}

	var kept []Cover
	for _, cover := range existing {
		if cover.Origin == CoverOriginRandom {
			continue
		}
		if removedSet[cover.MediaId] {
			continue
		}
		kept = append(kept, cover)
	}

	slotsToFill := MaxCoversPerAlbum - len(kept)
	alreadyCovered := make(map[MediaId]bool, len(kept))
	for _, cover := range kept {
		alreadyCovered[cover.MediaId] = true
	}

	result := kept
	if slotsToFill > 0 {
		picked := c.pickFromCandidates(added, alreadyCovered, slotsToFill)
		for _, media := range picked {
			alreadyCovered[media.Id] = true
			result = append(result, Cover{MediaId: media.Id, Filename: media.Filename, Origin: CoverOriginRandom})
		}

		if remaining := slotsToFill - len(picked); remaining > 0 {
			fallback, err := c.MediaReadRepository.FindMedias(ctx, NewFindMediaRequest(albumId.Owner).WithAlbum(albumId.FolderName))
			if err != nil {
				return errors.Wrapf(err, "Reconcile(%s) failed to list medias for the fallback query", albumId)
			}
			extra := c.pickFromCandidates(fallback, alreadyCovered, remaining)
			for _, media := range extra {
				alreadyCovered[media.Id] = true
				result = append(result, Cover{MediaId: media.Id, Filename: media.Filename, Origin: CoverOriginRandom})
			}
		}
	}

	if coversEqual(existing, result) {
		return nil
	}

	if err := c.CoverRepository.SaveCovers(ctx, albumId, result); err != nil {
		return errors.Wrapf(err, "Reconcile(%s) failed to save covers", albumId)
	}

	for _, observer := range c.Observers {
		if err := observer.OnCoversChanged(ctx, albumId, result); err != nil {
			return errors.Wrapf(err, "Reconcile(%s) failed to notify observer", albumId)
		}
	}
	return nil
}

func (c *CoverMaintenance) pickFromCandidates(candidates []*MediaMeta, alreadyCovered map[MediaId]bool, slotsToFill int) []*MediaMeta {
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
	picked := make([]*MediaMeta, 0, len(indices))
	for _, idx := range indices {
		picked = append(picked, eligible[idx])
	}
	return picked
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

// ReconcileCoversPort is the slice of CoverMaintenance exposed to callers that only
// need to trigger a reconciliation (per-event cover observers, admin backfill,
// Phase 3 owner-triggered re-randomise).
type ReconcileCoversPort interface {
	Reconcile(ctx context.Context, albumId AlbumId, added []*MediaMeta, removed []MediaId) error
}

// BackfillCovers reconciles the cover set of every album of an owner. Used by the
// administrative CLI to seed existing albums that pre-date the covers feature. A
// failure on a single album is reported and the sweep continues: Reconcile preserves
// CHERRY_PICKED covers and caps the set at MaxCoversPerAlbum, so re-running the
// backfill is safe.
type BackfillCovers struct {
	FindAlbumByOwnerPort FindAlbumByOwnerPort
	ReconcileCoversPort  ReconcileCoversPort
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
		if err := b.ReconcileCoversPort.Reconcile(ctx, album.AlbumId, nil, nil); err != nil {
			report.Failures = append(report.Failures, BackfillFailure{AlbumId: album.AlbumId, Err: err})
		}
	}
	return report, nil
}
