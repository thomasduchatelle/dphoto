package catalog

import (
	"context"
	"math/rand"
	"slices"

	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// CoverRepository persists an album's cover set as a single record.
type CoverRepository interface {
	FindCoversByAlbum(ctx context.Context, albumId AlbumId) ([]Cover, error)
	FindCoversByAlbums(ctx context.Context, albumIds ...AlbumId) (map[AlbumId][]Cover, error)
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

// Randomise regenerates every cover that is not CHERRY_PICKED, drops any CHERRY_PICKED
// cover whose media is no longer in the album, then fills up to MaxCoversPerAlbum by
// drawing uniformly at random from the album's full image set. The returned map carries
// one entry per album whose cover set actually changed; unchanged albums are omitted.
// Returns a nil map when no album changed.
func (c *CoverService) Randomise(ctx context.Context, albumIds ...AlbumId) (map[AlbumId][]Cover, error) {
	if len(albumIds) == 0 {
		return nil, nil
	}

	existing, err := c.CoverRepository.FindCoversByAlbums(ctx, albumIds...)
	if err != nil {
		return nil, errors.Wrapf(err, "Randomise failed to load covers for albums %v", albumIds)
	}

	var changed map[AlbumId][]Cover
	for _, albumId := range albumIds {
		medias, err := c.MediaReadRepository.FindMedias(ctx, NewFindMediaRequest(albumId.Owner).WithAlbum(albumId.FolderName))
		if err != nil {
			return nil, errors.Wrapf(err, "Randomise failed to list medias for %s", albumId)
		}

		inAlbum := make(map[MediaId]bool, len(medias))
		for _, media := range medias {
			if media == nil {
				continue
			}
			inAlbum[media.Id] = true
		}

		kept := keepCherryPickedStillInAlbum(existing[albumId], inAlbum)
		result := fillRandomSlots(kept, medias, c.Randomiser)

		if coversEqual(existing[albumId], result) {
			continue
		}
		if err := c.CoverRepository.SaveCovers(ctx, albumId, result); err != nil {
			return nil, errors.Wrapf(err, "Randomise failed to save covers for %s", albumId)
		}
		if changed == nil {
			changed = make(map[AlbumId][]Cover)
		}
		changed[albumId] = result
	}
	return changed, nil
}

// StableRefresh processes every album affected by `transferred`:
//
//   - For each source album, covers whose media has moved away are stripped.
//   - CHERRY_PICKED covers that moved to a destination album are inherited by the
//     destination (as CHERRY_PICKED), displacing the oldest RANDOM cover when needed.
//     Destinations already full of CHERRY_PICKED covers silently ignore the incoming
//     one.
//   - Every affected album has its empty slots backfilled from its full image set.
//
// Returns one entry per album whose cover set actually changed. Unchanged albums are
// omitted. Returns a nil map when no album changed.
func (c *CoverService) StableRefresh(ctx context.Context, transferred TransferredMedias) (map[AlbumId][]Cover, error) {
	if transferred.IsEmpty() {
		return nil, nil
	}

	affected := collectAffectedAlbums(transferred)
	if len(affected) == 0 {
		return nil, nil
	}

	existing, err := c.CoverRepository.FindCoversByAlbums(ctx, affected...)
	if err != nil {
		return nil, errors.Wrapf(err, "StableRefresh failed to load covers for albums %v", affected)
	}

	// Resolve which CHERRY_PICKED covers were moved to a destination so they can be
	// inherited there. Scanning every source album's covers once is enough.
	cherryPickedByDestination := make(map[AlbumId][]Cover)
	for sourceId, sourceCovers := range existing {
		for _, cover := range sourceCovers {
			if cover.Origin != CoverOriginCherryPicked {
				continue
			}
			if destId, ok := destinationOf(cover.MediaId, transferred.Transfers); ok && !destId.IsEqual(sourceId) {
				cherryPickedByDestination[destId] = append(cherryPickedByDestination[destId], cover)
			}
		}
	}

	movedOutOfAlbum := make(map[AlbumId]map[MediaId]bool, len(affected))
	for _, albumId := range affected {
		set := make(map[MediaId]bool)
		for destId, mediaIds := range transferred.Transfers {
			if destId.IsEqual(albumId) {
				continue
			}
			for _, mediaId := range mediaIds {
				set[mediaId] = true
			}
		}
		movedOutOfAlbum[albumId] = set
	}

	var changed map[AlbumId][]Cover
	for _, albumId := range affected {
		medias, err := c.MediaReadRepository.FindMedias(ctx, NewFindMediaRequest(albumId.Owner).WithAlbum(albumId.FolderName))
		if err != nil {
			return nil, errors.Wrapf(err, "StableRefresh failed to list medias for %s", albumId)
		}

		kept := stripMovedAway(existing[albumId], movedOutOfAlbum[albumId])
		kept = inheritCherryPicked(kept, cherryPickedByDestination[albumId])
		result := fillRandomSlots(kept, medias, c.Randomiser)

		if coversEqual(existing[albumId], result) {
			continue
		}
		if err := c.CoverRepository.SaveCovers(ctx, albumId, result); err != nil {
			return nil, errors.Wrapf(err, "StableRefresh failed to save covers for %s", albumId)
		}
		if changed == nil {
			changed = make(map[AlbumId][]Cover)
		}
		changed[albumId] = result
	}
	return changed, nil
}

func collectAffectedAlbums(transferred TransferredMedias) []AlbumId {
	var affected []AlbumId
	for destId := range transferred.Transfers {
		if !containsAlbumId(affected, destId) {
			affected = append(affected, destId)
		}
	}
	for _, sourceId := range transferred.FromAlbums {
		if !containsAlbumId(affected, sourceId) {
			affected = append(affected, sourceId)
		}
	}
	return affected
}

func containsAlbumId(ids []AlbumId, target AlbumId) bool {
	return slices.ContainsFunc(ids, target.IsEqual)
}

func destinationOf(mediaId MediaId, transfers map[AlbumId][]MediaId) (AlbumId, bool) {
	for destId, mediaIds := range transfers {
		for _, id := range mediaIds {
			if id == mediaId {
				return destId, true
			}
		}
	}
	return AlbumId{}, false
}

func keepCherryPickedStillInAlbum(covers []Cover, inAlbum map[MediaId]bool) []Cover {
	var kept []Cover
	for _, cover := range covers {
		if cover.Origin != CoverOriginCherryPicked {
			continue
		}
		if !inAlbum[cover.MediaId] {
			continue
		}
		kept = append(kept, cover)
	}
	return kept
}

func stripMovedAway(covers []Cover, movedAway map[MediaId]bool) []Cover {
	if len(covers) == 0 {
		return nil
	}
	kept := make([]Cover, 0, len(covers))
	for _, cover := range covers {
		if movedAway[cover.MediaId] {
			continue
		}
		kept = append(kept, cover)
	}
	return kept
}

// inheritCherryPicked adds every incoming CHERRY_PICKED cover to the kept set,
// displacing the oldest RANDOM cover when the set is already full. When no RANDOM
// cover can be displaced (the set is full of CHERRY_PICKED), the incoming cover is
// silently dropped. Incoming covers already in the kept set (same MediaId) are
// skipped to avoid duplicates.
func inheritCherryPicked(kept, incoming []Cover) []Cover {
	for _, cover := range incoming {
		if containsMediaId(kept, cover.MediaId) {
			continue
		}
		if len(kept) < MaxCoversPerAlbum {
			kept = append(kept, cover)
			continue
		}
		displaceAt := slices.IndexFunc(kept, func(c Cover) bool {
			return c.Origin == CoverOriginRandom
		})
		if displaceAt < 0 {
			continue
		}
		kept = append(kept[:displaceAt], kept[displaceAt+1:]...)
		kept = append(kept, cover)
	}
	return kept
}

func containsMediaId(covers []Cover, mediaId MediaId) bool {
	return slices.ContainsFunc(covers, func(c Cover) bool { return c.MediaId == mediaId })
}

// fillRandomSlots fills the empty slots of `kept` (up to MaxCoversPerAlbum) by
// drawing uniformly at random from `medias`. Only MediaType IMAGE candidates not
// already in `kept` are eligible.
func fillRandomSlots(kept []Cover, medias []*MediaMeta, randomiser Randomiser) []Cover {
	slotsToFill := MaxCoversPerAlbum - len(kept)
	if slotsToFill <= 0 {
		return kept
	}

	alreadyCovered := make(map[MediaId]bool, len(kept))
	for _, cover := range kept {
		alreadyCovered[cover.MediaId] = true
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
	if len(eligible) == 0 {
		return kept
	}

	picks := slotsToFill
	if picks > len(eligible) {
		picks = len(eligible)
	}
	indices := randomiser.SampleIndices(len(eligible), picks)
	result := kept
	for _, idx := range indices {
		media := eligible[idx]
		result = append(result, Cover{MediaId: media.Id, Filename: media.Filename, Origin: CoverOriginRandom})
	}
	return result
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

// RandomiseCoversPort is the slice of CoverService exposed to callers that need to
// regenerate the cover set of one or more albums (admin backfill, Phase 3
// owner-triggered re-randomise, and the InsertMedias use case).
type RandomiseCoversPort interface {
	Randomise(ctx context.Context, albumIds ...AlbumId) (map[AlbumId][]Cover, error)
}

// StableRefreshCoversPort is the slice of CoverService exposed to callers that need
// to patch cover sets after a media transfer (AlbumCreated, AlbumDatesAmended,
// AlbumDeleted use cases).
type StableRefreshCoversPort interface {
	StableRefresh(ctx context.Context, transferred TransferredMedias) (map[AlbumId][]Cover, error)
}

// BackfillCoversViewUpdater propagates the new covers of a single album to the
// album-list view. BackfillCovers calls it after every Randomise that actually
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
// A failure on a single album is reported and the sweep continues: Randomise
// preserves CHERRY_PICKED covers whose media is still in the album and caps the
// set at MaxCoversPerAlbum, so re-running the backfill is safe (though RANDOM
// covers may be redrawn on each pass).
type BackfillCovers struct {
	FindAlbumByOwnerPort      FindAlbumByOwnerPort
	RandomiseCoversPort       RandomiseCoversPort
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
		changed, err := b.RandomiseCoversPort.Randomise(ctx, album.AlbumId)
		if err != nil {
			report.Failures = append(report.Failures, BackfillFailure{AlbumId: album.AlbumId, Err: err})
			continue
		}
		covers, ok := changed[album.AlbumId]
		if !ok {
			continue
		}
		if err := b.BackfillCoversViewUpdater.UpdateCovers(ctx, album.AlbumId, covers); err != nil {
			report.Failures = append(report.Failures, BackfillFailure{AlbumId: album.AlbumId, Err: err})
		}
	}
	return report, nil
}
