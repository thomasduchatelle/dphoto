# Design — Album covers

Cross-issue technical direction. Issues link here rather than repeat it. The "what" is in `spec.md`.

## Read model (Phase 1 — done)

Captured in ADR-0004. The per-user album-list view record (`AlbumSummary`) carries identity, display
fields (`Name`, `Start`, `End`), `Count`, and `Covers`. Writes are split so count updates and
display-field updates never clobber each other. Reads are served by a single Query per user.

## Cover canonical record

The whole cover set of an album is a **single** DynamoDB item:

- PK `{OWNER}#ALBUM`, SK `ALBUM#{FOLDER_NAME}#COVERS`.
- Holds an ordered list of at most 4 covers, each `{MediaId, Filename, Origin}` where `Origin` is
  `RANDOM` or `CHERRY_PICKED`.
- Every write reads the set, enforces the max of 4, and rewrites the whole item.
- Deleted when the album is deleted; migrated (same owner, new folder-name SK) when the album is renamed
  with a folder change.

Documented in `DATA_MODEL.md`.

## Randomize operation

A single operation drives every automatic cover refresh. **Renames and supersedes `CompleteCovers`.**

Semantics:

1. Load the current cover set from the canonical record.
2. **Drop every `RANDOM` cover**; keep every `CHERRY_PICKED` cover.
3. Compute `slotsToFill = 4 − len(kept)`.
4. From the supplied eligible candidates (`MediaType IMAGE`, not already a kept cover), pick
   `slotsToFill` uniformly at random (or all of them if fewer are available) and tag them `RANDOM`.
5. Write the resulting set back and return it.

Two entry shapes (same semantics, different candidate source):

- `RandomizeCoversFromCandidates(ctx, albumId, candidates) ([]Cover, error)` — fills from a supplied set
  of medias, no extra query. Used by catalog handlers for `MediasInserted`, `AlbumCreated`,
  `AlbumDatesAmended` — the use case passes the medias it just inserted/transferred.
- `RandomizeCovers(ctx, albumId) ([]Cover, error)` — queries the album's `IMAGE` medias then randomises.
  Used by the admin backfill and by the Phase 3 owner-triggered re-randomise.

Returning the resulting `[]Cover` is essential: the caller attaches it to the lifecycle event so the view
denormalises in the same pass.

### Idempotence note

Because `RANDOM` covers are always dropped and redrawn, calling `Randomize` on a set made exclusively of
`CHERRY_PICKED` covers is a no-op. Calling it on a full set with at least one `RANDOM` cover **does**
change the covers (fresh draw). The admin backfill is therefore not strictly idempotent in the
"byte-for-byte identical output" sense; it is **safe to re-run** (`CHERRY_PICKED` always preserved,
count always `≤ 4`), which is what matters operationally.

## Event-payload matrix

Covers propagate through existing lifecycle events. Each event carries `Covers` iff its outcome can
change them. There is **no** standalone `CoversChanged` event, and no standalone `OnAlbumCoversChanged`
method on `AlbumView` (the one shipped in issue 03 is removed in Phase 2 cleanup).

| Event                 | Carries `Covers`? | Why                                                                              |
|-----------------------|-------------------|----------------------------------------------------------------------------------|
| `AlbumCreated`        | yes               | `TransferredMedias` may bring images in → `Randomize` fills empty slots.         |
| `MediasInserted`      | yes (per album)   | Inserts may fill empty slots → `Randomize` with the just-inserted images.        |
| `AlbumDatesAmended`   | yes (per album)   | Transfers in/out: stale covers stripped, empties re-filled on destinations.      |
| `AlbumRenamed`        | no                | Covers survive untouched; the handler migrates the canonical record's SK.        |
| `AlbumDeleted`        | no                | Covers are wiped; the handler deletes the canonical record.                      |
| `AlbumShared`         | yes               | New visitor row needs the current covers written in one shot.                    |

### Who computes the covers inside the use case

Each catalog use case that may change covers takes a new port:

```go
type RandomizeCoversPort interface {
    RandomizeCoversFromCandidates(ctx, albumId, candidates []*MediaMeta) ([]Cover, error)
}
```

Wiring (via `pkg/pkgfactory`) injects the existing `RandomizeCovers` struct. The use case calls the port
after it has inserted/transferred the medias, then emits the lifecycle event with the resulting `Covers`
attached. `AlbumShared` is slightly different: it reads the current covers via
`FindCoversByAlbum` (the album's canonical set never changes on share) and attaches them.

### Handling transfers-out

When medias move out of an album (`AlbumDatesAmended`, `AlbumRenamed` folder-change), the source
album's cover set may contain covers whose `MediaId` is no longer in the album. The use case must:

1. Load the source album's current covers.
2. Strip any cover whose `MediaId` is in the "moved out" set.
3. Attach the pruned set to the event (so the view denormalises). No re-randomisation here — empties
   stay empty until the next time medias land in the album or until the admin backfill.

The destination album is handled as a transfer-in: `Randomize` with the moved-in medias as candidates.

## Cover projection in the view

`AlbumSummary` carries `Covers`. The repository exposes `SetCoversForAllViewers(ctx, albumId, covers)`
for the standalone case, but Phase 2 moves it inside the per-event handlers:

- `OnAlbumCreated` writes a full row (name/start/end/count/covers) via `PutSummaries`.
- `OnMediasInserted` combines the count increment with a cover update: a single call per-viewer path
  that `SET`s `Covers` **in addition to** `ADD Count :d` — attribute footprints stay disjoint from
  display fields (no clobbering).
- `OnAlbumDatesAmended` already writes display fields; it also writes `Covers` for each affected album.
- `OnAlbumRenamed` folder-change: delete old rows + write new rows with the preserved covers.
- `OnAlbumRenamed` in-place: no cover change, no cover write.
- `OnAlbumDeleted`: `DeleteAllRowsForAlbum` as today.
- `OnAlbumShared`: `PutSummaries` for the visitor row with covers included (no second call).
- `OnAlbumUnshared`: `DeleteRow` as today.

The standalone `AlbumView.OnAlbumCoversChanged` method is deleted. Drift reconciliation already writes
full rows via `PutSummaries` (`Covers` included) — unchanged.

### Writing covers alongside counts

`AlbumSummaryRepository` gains a repository method (or the existing one gains a `Covers` argument) that
performs a single `UpdateItem` with `ADD Count :d SET Covers = :c`. Both attributes are updated
atomically per row; display fields stay untouched.

Chosen shape — new method:
```go
IncrementCountAndSetCoversForAllViewers(ctx, []AlbumCountAndCoversDiff) error
```
with `AlbumCountAndCoversDiff{AvailabilityType, AlbumId, Diff int, Covers []Cover}`. Covers is set
unconditionally (so empty-after-strip is a valid write). Existing `IncrementCountForAllViewers` stays
for the (now rare) case where the use case has no covers to propagate (shouldn't happen in production
Phase 2 — see below).

## REST contract

### Phase 2 (done in issue 04)

`GET /api/v1/albums` — each album in the response has `covers[]` (0–4 entries):

```json
"covers": [
  { "mediaId": "...", "filename": "...", "origin": "RANDOM" }
]
```

`origin` ∈ `RANDOM | CHERRY_PICKED`. The frontend builds image URLs from `owner` + `mediaId` +
`filename` via its existing image loader — no URL is stored or returned.

### Phase 3 — owner-triggered re-randomise

| Phase | Action        | Route                                                              |
|-------|---------------|--------------------------------------------------------------------|
| 3     | Re-randomise  | `POST /api/v1/owners/{owner}/albums/{folderName}/covers/refresh`   |

- Authorisation: owner-edit permission on the album (same rule as rename / amend-dates).
- Body: empty.
- Response: `200 OK` with the new covers list (so the UI can refresh without re-querying the album list).
- Side effect: emits an event carrying the new covers → view updates via the standard path.

### Phase 4+ (deferred)

A single `PUT /api/v1/owners/{owner}/albums/{folderName}/covers` endpoint will replace the per-media
pick/unpick from the earlier draft, taking the full ordered list of up to 4 covers. Not specified here.
