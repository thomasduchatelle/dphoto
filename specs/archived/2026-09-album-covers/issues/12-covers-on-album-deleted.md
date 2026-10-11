# 12 — Covers on AlbumDeleted

Status: done
Phase: 2
Layer: `pkg/catalog` + `pkg/catalogadapters/catalogdynamo` + `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 09 (cover-maintenance service, split view projection)

## Description

When an album is deleted, its canonical cover record must disappear along with every viewer's
cover row for that album. Any neighbouring album that absorbed the deleted album's medias must
have its covers reconciled with stable semantics so the new medias may fill empty slots without
disturbing the existing covers.

See `../design.md` → per-operation strategy (`AlbumDeleted`); Cover canonical record.

## Acceptance criteria

- Deleting an album results in:
  - Its canonical cover record being removed (idempotent — no error if the record is already
    absent).
  - Every viewer's cover row for that album being removed from the album-list view (owner and all
    visitors).
  - Each destination album that absorbed transferred medias having its covers reconciled: existing
    covers preserved, empty slots filled from the destination album's full image set.
- Deleting an album with no media transfer still removes the canonical record and the viewer rows;
  no other album is touched.
- Destination albums whose cover set is unaffected by the transfer do not get a cover-row rewrite.
- After deletion, `ListAlbums` returns neither the deleted album nor any cover row referencing it
  for any viewer.
- `DATA_MODEL.md` notes that the canonical cover record is deleted when the album is deleted.

## Out of scope

- `MediasInserted`, `AlbumCreated`, `AlbumDatesAmended`, `AlbumRenamed`, `AlbumShared` — issues
  09, 10, 11, 16, 17.

## References

- `../spec.md` (Automatic cover maintenance → album deleted).
- `../design.md` (Cover canonical record; per-operation strategy — `AlbumDeleted`).
- Load skills: `go`, `architecture`.
