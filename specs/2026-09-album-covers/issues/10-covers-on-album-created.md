# 10 — AlbumCreated cover maintenance

Status: ready
Phase: 2
Layer: `pkg/catalog` + `pkg/pkgfactory`
Depends on: 09 (reconciliation primitive, split view projection)

## Description

When an album is created, its covers must be initialised from the medias transferred into it, and any
pre-existing album that lost medias to the new one must have its covers reconciled too (stale covers
stripped, empties refilled from survivors).

Shipped as a cover observer on the `AlbumCreated` lifecycle event. The create-album use case is not
modified; it continues to emit the event as it does today.

See `../design.md` → Cover-maintenance pattern; per-operation reconciliation inputs (`AlbumCreated`).

## Acceptance criteria

- Creating an album with no transferred-in medias results in no cover record for that album and no
  cover row on any viewer's `ListAlbums` entry.
- Creating an album that overlaps and absorbs medias from pre-existing albums results in:
  - The new album having up to 4 `RANDOM` covers drawn from the transferred-in medias (fewer than 4
    if fewer eligible `IMAGE` medias were transferred).
  - Each source album whose cover referenced a transferred-out media having that cover stripped and
    its slot refilled from its remaining `IMAGE` medias (fallback query). `CHERRY_PICKED` covers on
    source albums are preserved.
  - Source albums whose covers were untouched by the transfer having no cover-record write.
- All of the above is visible on every viewer's next `ListAlbums` (owner and any visitor who already
  had access to the source albums).
- The `AlbumCreated` event payload is not changed. The create-album use case is not modified.
- Observer wiring: the cover observer is registered on the `AlbumCreated` event alongside the
  existing view observer; order between the two is not load-bearing (covers propagate through their
  own chain).

## Out of scope

- `MediasInserted`, `AlbumDatesAmended`, `AlbumDeleted`, `AlbumRenamed`, `AlbumShared` — issues 09,
  11, 12, 16, 17.
- Any change to `AlbumView.OnAlbumCreated`.

## References

- `../spec.md` (Automatic cover maintenance → medias added, medias moved out).
- `../design.md` (per-operation reconciliation inputs — `AlbumCreated`).
- Load skills: `go`, `architecture`.
