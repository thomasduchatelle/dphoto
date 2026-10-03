# 12 — AlbumDeleted cover maintenance

Status: ready
Phase: 2
Layer: `pkg/catalog` + `pkg/catalogadapters/catalogdynamo` + `pkg/pkgfactory`
Depends on: 09 (reconciliation primitive, split view projection)

## Description

When an album is deleted, its canonical cover record must be removed and every viewer's cover row for
that album must disappear. Any neighbouring album that absorbed the deleted album's medias must have
its covers reconciled from the moved-in medias.

Shipped as a cover observer on the `AlbumDeleted` lifecycle event, plus a repository extension to
remove the canonical `#COVERS` record.

See `../design.md` → Cover canonical record; per-operation reconciliation inputs (`AlbumDeleted`).

## Acceptance criteria

- Deleting an album results in:
  - Its canonical `#COVERS` record being removed (idempotent — no error if the record is already
    absent, e.g. the album never had covers).
  - Every viewer's cover row for that album being removed from the album-list view (owner and all
    visitors).
  - Each destination album that absorbed transferred medias having its covers reconciled with the
    moved-in medias (same rules as `AlbumDatesAmended` destinations: fill empty slots up to 4,
    preserve `CHERRY_PICKED`).
- Deleting an album with no media transfer (no neighbour absorbed its medias — e.g. empty album)
  still removes the canonical record and the viewer rows; no other album is touched.
- After deletion, `ListAlbums` returns neither the deleted album nor any cover row referencing it
  for any viewer.
- The `AlbumDeleted` event payload is not changed. The delete-album use case is not modified.
- `DATA_MODEL.md` notes that the `#COVERS` canonical record is deleted when the album is deleted.

## Out of scope

- `MediasInserted`, `AlbumCreated`, `AlbumDatesAmended`, `AlbumRenamed`, `AlbumShared` — issues 09,
  10, 11, 16, 17.
- Any change to `AlbumView.OnAlbumDeleted`.

## References

- `../spec.md` (Automatic cover maintenance → album deleted).
- `../design.md` (Cover canonical record; per-operation reconciliation inputs — `AlbumDeleted`).
- Load skills: `go`, `architecture`.
