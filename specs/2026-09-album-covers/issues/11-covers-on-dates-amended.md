# 11 — Covers on AlbumDatesAmended

Status: ready
Phase: 2
Layer: `pkg/catalog` + `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 09 (cover-maintenance service, split view projection)

## Description

Amending an album's dates can transfer medias between albums. Covers of every album affected by
that transfer must be reconciled with **stable** semantics: existing covers are preserved whenever
possible, only those whose medias moved out are stripped, and empty slots are refilled from the
album's remaining images.

See `../design.md` → per-operation strategy (`AlbumDatesAmended`).

## Acceptance criteria

- Amending dates with no media transfer (same effective date range) results in no cover-record
  write and no cover row change on any viewer.
- Amending dates that transfers medias from album A to album B results in:
  - For B (destination): existing covers kept, empty slots filled from B's full image set (which
    now includes the moved-in medias).
  - For A (source): any cover whose `MediaId` was moved out is stripped; the slot is refilled from
    A's remaining images. `CHERRY_PICKED` covers on A whose medias stayed are preserved.
  - Both reconciliations visible on every viewer's next `ListAlbums` for A and B.
- The amended album itself, when it is also a transfer source or destination, is reconciled once
  under the same rules.
- Albums whose cover set is unchanged by the amendment do not get a cover-row rewrite.

## Out of scope

- `MediasInserted`, `AlbumCreated`, `AlbumDeleted`, `AlbumRenamed`, `AlbumShared` — issues 09, 10,
  12, 16, 17.

## References

- `../spec.md` (Automatic cover maintenance → medias moved out, medias added).
- `../design.md` (per-operation strategy — `AlbumDatesAmended`).
- Load skills: `go`, `architecture`.
