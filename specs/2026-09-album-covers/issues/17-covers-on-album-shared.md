# 17 — AlbumShared / AlbumUnshared cover propagation

Status: ready
Phase: 2
Layer: `pkg/acl/catalogacl` + `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 09 (reconciliation primitive, split view projection)

## Description

Sharing an album must make the album appear on the visitor's `ListAlbums` with its current covers in
a single Query (no delayed propagation, no reliance on the next catalog mutation to populate the
visitor's cover row). Unsharing must remove the album and its cover row from the visitor's view.

Covers themselves are not changed by sharing — the canonical `#COVERS` record is owner-defined and
shared across viewers. This story is about the view-side cover row for the newly-granted (or
revoked) visitor.

See `../design.md` → View projection; per-operation reconciliation inputs (table — not covered,
sharing writes the visitor row directly).

## Acceptance criteria

- Sharing an album with a visitor results in the visitor's next `ListAlbums` returning the album
  with the owner's current covers, served from the single-partition Query (no fall-back read).
- Sharing an album that has no covers results in the album appearing on the visitor's `ListAlbums`
  with an empty cover list; no cover row is written.
- Unsharing an album from a visitor results in the visitor's next `ListAlbums` not returning the
  album, and no residual cover row for that album in the visitor's partition.
- Sharing or unsharing an album does not modify the canonical `#COVERS` record, does not modify any
  other viewer's cover row, and does not change the covers served to the owner.
- Sharing the same album twice with the same visitor is a no-op on the cover row (idempotent).

## Out of scope

- Re-propagating covers to existing visitors on cover change — handled by issue 09's cover
  projection observer, not by sharing.
- Any change to how `AlbumShared` / `AlbumUnshared` write the main album-summary row (that part is
  already correct from Phase 1).

## References

- `../spec.md` (Automatic cover maintenance → album shared).
- `../design.md` (View projection).
- Load skills: `go`, `architecture`.
