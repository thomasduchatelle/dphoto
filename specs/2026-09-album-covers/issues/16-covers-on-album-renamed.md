# 16 — Covers on AlbumRenamed

Status: done
Phase: 2
Layer: `pkg/catalog` + `pkg/catalogadapters/catalogdynamo` + `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 09 (cover-maintenance service, split view projection)

## Description

Renaming an album must preserve its covers. When the rename is in-place (name change only, same
folder), nothing on the covers side needs to happen. When the folder changes, the canonical cover
record moves to the new SK and every viewer's cover row follows the album to its new identity,
with the same ordered cover list and the same origins.

See `../design.md` → Cover canonical record; per-operation strategy (`AlbumRenamed`).

## Acceptance criteria

- Renaming an album in-place results in no cover-record write, no change to any viewer's cover
  row, and the renamed album still showing its original covers on the next `ListAlbums`.
- Renaming an album with a folder change results in:
  - The canonical cover record present at the new identity with the same ordered cover list; the
    old identity carrying no cover record.
  - Every viewer's cover row for the album following to the new identity, for owner and all
    visitors.
  - `ListAlbums` returning the renamed album with its preserved covers on the next call.
- A folder-change rename on an album with no covers is a no-op on both the canonical record and
  the viewer rows.
- `CHERRY_PICKED` and `RANDOM` origins are preserved verbatim across the rename — no
  reconciliation, no redraw.
- `DATA_MODEL.md` notes that the canonical cover record moves to the new identity on folder-change
  rename and is unchanged on in-place rename.

## Out of scope

- `MediasInserted`, `AlbumCreated`, `AlbumDatesAmended`, `AlbumDeleted`, `AlbumShared` — issues
  09, 10, 11, 12, 17.

## References

- `../spec.md` (Automatic cover maintenance → album renamed).
- `../design.md` (Cover canonical record; per-operation strategy — `AlbumRenamed`).
- Load skills: `go`, `architecture`.
