# 11 — AlbumDatesAmended cover maintenance

Status: ready
Phase: 2
Layer: `pkg/catalog` + `pkg/pkgfactory`
Depends on: 09 (reconciliation primitive, split view projection)

## Description

Amending an album's dates can transfer medias between albums. Covers of every album affected by that
transfer must be reconciled: destinations get their empty slots filled from the moved-in medias;
sources get their stale covers stripped and refilled from survivors.

Shipped as a cover observer on the `AlbumDatesAmended` lifecycle event. The amend-dates use case is
not modified.

See `../design.md` → Cover-maintenance pattern; per-operation reconciliation inputs
(`AlbumDatesAmended`).

## Acceptance criteria

- Amending dates with no media transfer (same date range in practice, or a shift that moves no
  media) results in no cover-record write and no cover row change on any viewer.
- Amending dates that transfers medias from album A to album B results in:
  - For B (destination): empty cover slots filled with `RANDOM` covers drawn from the moved-in
    medias, up to 4. `CHERRY_PICKED` and already-set `RANDOM` covers on B remain untouched (modulo
    the "drop RANDOM, redraw" step applied by the primitive — the invariant).
  - For A (source): any cover whose `MediaId` was moved out is stripped; the slot is refilled from
    A's remaining `IMAGE` medias (fallback query). `CHERRY_PICKED` covers on A whose medias stayed
    are preserved.
  - Both reconciliations visible on every viewer's next `ListAlbums` for A and B (owner and any
    visitors).
- The amended album itself, when it is also a transfer source or destination, is reconciled under
  the same rules (one reconciliation per affected album, no double-run).
- Amending dates on an album whose medias all stay (no transfer anywhere) results in no cover
  changes.
- The `AlbumDatesAmended` event payload is not changed. The amend-dates use case is not modified.

## Out of scope

- `MediasInserted`, `AlbumCreated`, `AlbumDeleted`, `AlbumRenamed`, `AlbumShared` — issues 09, 10,
  12, 16, 17.
- Any change to `AlbumView.OnAlbumDatesAmended`.

## References

- `../spec.md` (Automatic cover maintenance → medias moved out, medias added).
- `../design.md` (per-operation reconciliation inputs — `AlbumDatesAmended`).
- Load skills: `go`, `architecture`.
