# 09 — Cover reconciliation primitive + split view projection + MediasInserted

Status: ready
Phase: 2
Layer: `pkg/catalog` + `pkg/catalogviews` + `pkg/catalogviewsadapters/catalogviewsdynamodb` + `pkg/pkgfactory` + `cmd/dphotops`
Depends on: —

## Description

Ship the cover-maintenance pattern end-to-end, plus its first consumer (`MediasInserted`). After this
story:

- A single primitive, `CoverMaintenance.Reconcile`, applies the full cover invariant and emits a
  `CoversChanged` signal when the set actually changes.
- The album-list view holds covers as a sibling row (`#COVERS` SK) in the user's partition, kept in
  sync by a single projection observer listening to `CoversChanged`.
- Inserting medias into an album causes its covers to be reconciled from the newly-inserted images,
  with the result visible on every viewer's next `ListAlbums`.
- The admin backfill (`dphotops covers backfill`) uses the new primitive and keeps its semantics
  (safe to re-run, `CHERRY_PICKED` always preserved).

See `../design.md` → Cover-maintenance pattern, View projection, per-operation reconciliation inputs.

## Acceptance criteria

### Reconciliation primitive

- A primitive owned by the catalog takes `(albumId, added []*MediaMeta, removed []MediaId)` and
  applies the invariant below in one call:
  - Every `RANDOM` cover is dropped; every `CHERRY_PICKED` cover is kept.
  - Any kept cover whose `MediaId` is in `removed` is stripped.
  - Empty slots are filled up to 4 by picking uniformly at random from eligible candidates
    (`MediaType IMAGE`, not already a kept cover): first from `added`, then — only if `added` is
    insufficient to fill the slots — from the album's remaining `IMAGE` medias queried on demand.
  - The resulting set is persisted to the canonical `#COVERS` record.
- The primitive emits a `CoversChanged` notification carrying the new set **only when the set
  differs** from the one loaded at step 1. No notification is emitted when the reconciliation is a
  no-op (e.g. full `CHERRY_PICKED` set, no `removed` match, no empty slot to fill).
- Called with `added = nil` and `removed = nil`, the primitive still applies the full invariant (it
  triggers the fallback query to refill empty slots).
- Admin backfill (`BackfillCovers` + `dphotops covers backfill`) uses the primitive. Its CLI help
  text reflects the new semantics: re-running the backfill may **change** existing `RANDOM` covers
  (fresh draw); `CHERRY_PICKED` covers are preserved; the operation is safe to re-run on any album.

### MediasInserted cover maintenance

- Inserting medias into one or more albums reconciles each affected album's covers using the
  just-inserted images as `added` (no `removed`).
- The `MediasInserted` event payload is **not changed**. The media-insert use case does not call
  the primitive directly; the observer wiring attaches to the existing event bus.
- `AlbumView.OnMediasInserted` is **not changed**: it continues to update counts only. Covers
  propagate independently through the view's cover-projection observer.

### Split view projection

- The album-list view no longer stores covers on `AlbumSummary`. Instead, each album has two SK
  rows per viewer in the user's `USER#{EMAIL}#ALBUMS_VIEW` partition: one for identity + display
  fields + count, one for the cover list.
- `ListAlbums` is still served by a **single Query** on the user's partition. The DynamoDB adapter
  merges the two SK families per album before returning. A missing cover row for a known album
  means "no covers" (empty list). An orphan cover row (no matching summary row) is discarded.
- A single observer in `pkg/catalogviews` reacts to `CoversChanged` by writing (or deleting, if the
  new set is empty) the `#COVERS` SK row for every viewer who can access the album.
- `AlbumSummary.Covers` and `AlbumView.OnAlbumCoversChanged` are removed. The repository method
  `SetCoversForAllViewers` is removed. The in-memory fake and DynamoDB adapter both implement the
  new cover-row storage.
- Drift reconciliation includes the `#COVERS` SK row in its rebuild path so rebuilt viewer rows
  carry the current covers.

### End-to-end behaviour

- Inserting 3 `IMAGE` medias into an album with no covers results, after the operation returns, in a
  canonical `#COVERS` record with 3 `RANDOM` entries and in `ListAlbums` returning the album with
  those 3 covers for the owner and every viewer.
- Inserting medias into an album whose cover set is already full of `CHERRY_PICKED` leaves the
  canonical record, the viewer rows, and the `ListAlbums` output unchanged.
- Inserting medias into an album whose cover set has 2 `CHERRY_PICKED` and 1 `RANDOM` cover ends
  with 2 `CHERRY_PICKED` preserved and up to 2 fresh `RANDOM` entries drawn from eligible candidates
  (preferring the inserted medias).
- Running `dphotops covers backfill` on an owner with a mix of empty, partially-covered, and
  fully-`CHERRY_PICKED` albums leaves every album with a legal cover set (≤ 4, `CHERRY_PICKED`
  preserved, `RANDOM` possibly redrawn).

### Documentation

- `DATA_MODEL.md` reflects the view-side split: two SK rows per album per viewer in the
  `USER#{EMAIL}#ALBUMS_VIEW` partition.

## Out of scope

- Reconciliation on `AlbumCreated`, `AlbumDatesAmended`, `AlbumDeleted`, `AlbumRenamed`,
  `AlbumShared`, `AlbumUnshared` — issues 10, 11, 12, 16, 17.
- Any change to catalog use-case code or event payloads (reconciliation hooks in reactively via the
  observer pattern).
- Owner-triggered re-randomise — Phase 3.

## References

- `../spec.md` (Automatic cover maintenance → medias added).
- `../design.md` (Cover-maintenance pattern; View projection).
- Load skills: `go`, `architecture`.
