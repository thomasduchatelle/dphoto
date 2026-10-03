# 16 — AlbumRenamed cover retention

Status: ready
Phase: 2
Layer: `pkg/catalog` + `pkg/catalogadapters/catalogdynamo` + `pkg/pkgfactory`
Depends on: 09 (reconciliation primitive, split view projection)

## Description

Renaming an album must preserve its covers. When the rename is in-place (name change only, same
folder), nothing needs to happen. When the folder changes, the canonical cover record moves to the
new SK and every viewer's cover row follows the album to its new identity.

Shipped as a cover observer on the `AlbumRenamed` lifecycle event, plus a repository extension to
move the canonical `#COVERS` record between SKs.

See `../design.md` → Cover canonical record; per-operation reconciliation inputs (`AlbumRenamed`).

## Acceptance criteria

- Renaming an album in-place (same `AlbumId`, name change only) results in no cover-record write,
  no change to any viewer's cover row, and the renamed album still showing its original covers on
  the next `ListAlbums`.
- Renaming an album with a folder change results in:
  - The canonical `#COVERS` record present at the new SK with the same ordered cover list; the old
    SK containing no cover record.
  - Every viewer's cover row for the album following to the new identity (old row gone, new row
    present with the same covers), for owner and all visitors.
  - `ListAlbums` returning the renamed album with its preserved covers on the next call.
- A folder-change rename on an album with no covers is a no-op on the canonical record (idempotent)
  and still leaves no cover row at either old or new identity.
- `CHERRY_PICKED` and `RANDOM` origins are preserved verbatim across the rename — no reconciliation,
  no redraw.
- The `AlbumRenamed` event payload is not changed. The rename-album use case is not modified.
- `DATA_MODEL.md` notes that the `#COVERS` canonical record moves to the new SK on folder-change
  rename and is unchanged on in-place rename.

## Out of scope

- `MediasInserted`, `AlbumCreated`, `AlbumDatesAmended`, `AlbumDeleted`, `AlbumShared` — issues 09,
  10, 11, 12, 17.
- Any change to `AlbumView.OnAlbumRenamed`.

## References

- `../spec.md` (Automatic cover maintenance → album renamed).
- `../design.md` (Cover canonical record; per-operation reconciliation inputs — `AlbumRenamed`).
- Load skills: `go`, `architecture`.
