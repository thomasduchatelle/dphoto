# 10 — Covers on AlbumCreated

Status: done
Phase: 2
Layer: `pkg/catalog` + `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 09 (cover-maintenance service, split view projection)

## Description

When an album is created, its covers must be initialised from the medias transferred into it (if
any), and any pre-existing album that lost medias to the new one must have its covers reconciled
too — stripped of covers whose media was moved out, refilled from the remaining medias so stable
covers persist whenever possible.

See `../design.md` → per-operation strategy (`AlbumCreated`).

## Acceptance criteria

- Creating an album with no transferred-in medias results in no cover record for that album and no
  cover row on any viewer's `ListAlbums` entry.
- Creating an album that absorbs medias from pre-existing albums results in:
  - The new album having up to 4 `RANDOM` covers drawn from the album's full image set (which, for
    a freshly created album, equals the transferred-in medias).
  - Each source album whose cover referenced a transferred-out media having that cover stripped and
    its slot refilled from the album's remaining images (**Stabilise** semantics). Source covers
    that were not affected by the transfer stay in place.
  - `CHERRY_PICKED` covers on source albums are preserved.
  - Source albums whose cover set was unaffected (no transferred-out media was a cover) do not get
    a cover-row rewrite.
- All of the above is visible on every viewer's next `ListAlbums` (owner and any visitor who
  already had access to the source albums).

## Out of scope

- `MediasInserted`, `AlbumDatesAmended`, `AlbumDeleted`, `AlbumRenamed`, `AlbumShared` — issues
  09, 11, 12, 16, 17.

## References

- `../spec.md` (Automatic cover maintenance → medias added, medias moved out).
- `../design.md` (per-operation strategy — `AlbumCreated`).
- Load skills: `go`, `architecture`.
