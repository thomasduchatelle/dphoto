# 03 — Covers in the album-list read model

Status: ready
Phase: 2
Layer: catalog domain — `pkg/catalogviews` + `pkg/catalogviewsadapters/catalogviewsdynamodb`
Depends on: 01, 02

## Description

Denormalise an album's cover set into the per-user album-list view record so covers are returned with the
album list in a single read, and keep them in sync when covers change.

See `../design.md` (Cover projection in the view).

## Acceptance criteria

- The view record (`AlbumSizeRecord`) carries a `Covers` attribute mirroring the canonical set — each entry
  `{MediaId, Filename, Origin}`.
- When an album's covers change, the set is re-denormalised to every viewer's record (same viewer fan-out as
  counts and display fields).
- `VisibleAlbum` carries the covers, so `AlbumView.ListAlbums` returns them without extra queries.
- Drift control rebuilds `Covers` from the canonical `…#COVERS` record.
- Tests follow the Go testing strategy; view tests cover an album with 0 and with 4 covers.

## Out of scope

- Exposing covers over REST (issue 04).
- Producing/altering covers (issues 02, 05, 06 and later phases).

## References

- `../spec.md`, `../design.md`
- Load skills: `go`, `architecture`.
