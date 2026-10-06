# 09 — Covers on MediasInserted + split view projection + Refresh/Stabilise service

Status: done
Phase: 2
Layer: `pkg/catalog` + `pkg/catalogviews` + `pkg/catalogviewsadapters/catalogviewsdynamodb` + `pkg/pkgfactory` + `cmd/dphotops`
Depends on: —

## Description

Deliver the cover-maintenance service and the view-side storage split together with the first
consumer (`MediasInserted`). After this story:

- Inserting medias into an album results in covers being drawn uniformly at random from the album's
  full image set, with `CHERRY_PICKED` covers preserved and existing `RANDOM` covers redrawn.
- The album-list view surfaces the covers of every album for every viewer, served from a single
  Query.
- The admin backfill (`dphotops covers backfill`) uses the same primitive and keeps its semantics
  (safe to re-run, `CHERRY_PICKED` always preserved).

See `../design.md` → Cover-maintenance service, View projection, per-operation strategy.

## Acceptance criteria

### Cover reconciliation strategies

- The catalog exposes two reconciliation strategies, both operating on the canonical `#COVERS`
  record of a single album:
  - **Refresh**: every `RANDOM` cover is dropped, every `CHERRY_PICKED` cover is kept, any kept
    cover whose `MediaId` is in the "removed" set is stripped, then empty slots are filled up to 4
    by picking uniformly at random from the album's full `IMAGE` set.
  - **Stabilise**: every existing cover is kept (both origins) except those in the "removed" set,
    then empty slots are filled up to 4 by picking uniformly at random from the album's full
    `IMAGE` set.
- Both strategies cap the resulting set at 4, select only `MediaType IMAGE` candidates, and never
  pick a media that is already a cover.
- A reconciliation that leaves the cover set unchanged does not rewrite the canonical record and
  does not surface the album to downstream consumers.

### MediasInserted cover update

- Inserting medias into one or more albums reconciles each affected album's covers with the
  **Refresh** strategy before the `MediasInserted` event is fired. The resulting cover sets travel
  on the event payload.
- Inserting 3 `IMAGE` medias into an album with no covers results in `ListAlbums` returning that
  album with 3 `RANDOM` covers for the owner and every viewer.
- Inserting medias into an album whose cover set is already full of `CHERRY_PICKED` leaves the
  canonical record and the viewer rows unchanged.
- Inserting medias into an album whose cover set has 2 `CHERRY_PICKED` and 1 `RANDOM` cover ends
  with the 2 `CHERRY_PICKED` preserved and up to 2 fresh `RANDOM` entries drawn from the album's
  full image set.
- An insertion that does not change the cover set of an album (e.g. only videos, or already-full
  `CHERRY_PICKED` set) does not result in a cover-row write for that album.

### Split view projection

- The album-list view holds covers in a sibling row rather than as a field on the summary. Each
  album has two SK rows per viewer in the user's `USER#{EMAIL}#ALBUMS_VIEW` partition: one for
  identity + display fields + count, one for the cover list.
- `ListAlbums` is still served by a **single Query** on the user's partition. A missing cover row
  for a known album means "no covers" (empty list). An orphan cover row (no matching summary row)
  is discarded.
- Cover writes never touch the main summary row (no clobbering of count or display fields).

### Admin backfill

- `dphotops covers backfill` applies the **Refresh** strategy to every album of every owner.
  Running it on a mix of empty, partially-covered, and fully-`CHERRY_PICKED` albums leaves every
  album with a legal cover set (≤ 4, `CHERRY_PICKED` preserved, `RANDOM` possibly redrawn).
- The CLI help text documents that re-running the backfill MAY redraw existing `RANDOM` covers
  (fresh draw) and that `CHERRY_PICKED` covers are always preserved.

### Documentation

- `DATA_MODEL.md` reflects the view-side split: two SK rows per album per viewer in the
  `USER#{EMAIL}#ALBUMS_VIEW` partition.

## Out of scope

- Covers on `AlbumCreated`, `AlbumDatesAmended`, `AlbumDeleted`, `AlbumRenamed`, `AlbumShared` —
  issues 10, 11, 12, 16, 17.
- Owner-triggered re-randomise — Phase 3.

## References

- `../spec.md` (Automatic cover maintenance → medias added).
- `../design.md` (Cover-maintenance service; View projection).
- Load skills: `go`, `architecture`.
