# 12 — Canonical cover record lifecycle (delete + folder-rename migration)

Status: ready
Phase: 2
Layer: catalog domain — `pkg/catalog` + `pkg/catalogadapters/catalogdynamo` + `pkg/pkgfactory`
Depends on: —

## Description

The `#COVERS` canonical DynamoDB record (one per album) must track its album's lifecycle:

- When the album is **deleted**, the record is deleted too (today it's left orphaned).
- When the album is **renamed with a folder change**, the record moves to the new SK (today the old
  record is left behind and the new album has no covers). In-place rename (same folder, name only) is
  a no-op — the SK is unchanged.

Both are small operations on the same repository, triggered by two sibling lifecycle events, so they
ship together: one PR, one test file, one repository extension.

The AlbumView's own handling of these events (deleting viewer rows on `AlbumDeleted`, deleting+writing
rows on folder-rename) is already correct from Phase 1 — this ticket only touches the canonical
`#COVERS` record and the view's cover attribute on the new rows written by folder-rename.

See `../design.md` → Cover canonical record; Event-payload matrix (rows `AlbumDeleted`, `AlbumRenamed`).

## Acceptance criteria

### Repository extension (`pkg/catalog` + `pkg/catalogadapters/catalogdynamo`)

- `CoverRepository` interface extended with:
  - `DeleteCovers(ctx, albumId AlbumId) error` — idempotent (no error when the record is absent).
  - `MoveCovers(ctx, from AlbumId, to AlbumId) error` — move the covers record from one SK to
    another under the same owner. Idempotent when the source is absent (no-op, no error). Asserts
    `from.Owner == to.Owner` (defensive — rename is same-owner by construction).
- DynamoDB adapter:
  - `DeleteCovers` → single `DeleteItem` on `CoverPrimaryKey(albumId)`.
  - `MoveCovers` → prefer a `TransactWriteItems` with a `Put` (new SK, same covers list,
    `AlbumFolderName` field updated) + `Delete` (old SK). If the source is absent, no-op. Document
    the chosen primitive in the PR body.
- In-memory fake implements both methods.

### Observers (`pkg/catalog`)

- `CoverRepositoryDeleteObserver` implementing `AlbumDeletedObserver` — on `OnAlbumDeleted`, calls
  `DeleteCovers(event.DeletedAlbumId)`.
- `CoverRepositoryRenameObserver` implementing `AlbumRenamedObserver` — on `OnAlbumRenamed`:
  - If `event.ExistingAlbum.AlbumId == event.RenamedAlbum.AlbumId` (in-place rename), no-op.
  - Otherwise, `MoveCovers(event.ExistingAlbum.AlbumId, event.RenamedAlbum.AlbumId)`.

### View handling of folder-rename covers (`pkg/catalogviews`)

- `AlbumView.OnAlbumRenamed` folder-change path (introduced in issue 01-05) currently writes the new
  rows with `Covers = nil`. Fix it: read the current covers via a new `FindCoversByAlbumPort`
  dependency and include them in the `PutSummaries` call that writes the new rows.
- **Observer ordering matters**: the cover-move observer must run **before** the AlbumView observer
  so the covers are at the new SK when the view handler reads them. Register them in that order in
  the factory, and document the guarantee in the PR body.
- Alternative, if the ordering guarantee is awkward: the AlbumView adapter reads from the **old** SK
  before the move happens, by reading in the adapter and buffering the result. PR documents the
  chosen approach.

### Wiring (`pkg/pkgfactory`)

- `DeleteAlbumCase` registers `CoverRepositoryDeleteObserver` alongside the existing AlbumView
  observer.
- `RenameAlbumCase` registers `CoverRepositoryRenameObserver` before the AlbumView observer; the
  AlbumView adapter gains the `FindCoversByAlbumPort` dependency.

### Tests

- `pkg/catalog/covers_test.go` (or `covers_lifecycle_test.go`):
  - `TestCoverRepositoryDeleteObserver/it_should_delete_the_canonical_record` — seed covers; dispatch
    `AlbumDeleted`; assert `FindCoversByAlbum` returns empty.
  - `TestCoverRepositoryDeleteObserver/it_should_be_idempotent_on_an_album_with_no_covers`.
  - `TestCoverRepositoryRenameObserver/it_should_be_a_noop_on_in_place_rename`.
  - `TestCoverRepositoryRenameObserver/it_should_move_covers_on_folder_change`.
  - `TestCoverRepositoryRenameObserver/it_should_be_a_noop_when_album_has_no_covers`.
- DynamoDB adapter test for `DeleteCovers` and `MoveCovers`:
  - Seed, operate, assert via `FindCoversByAlbum` on both old and new SKs.
  - Absent-source cases return `nil`.
- `pkg/catalogviews/albums_view_events_test.go`:
  - `TestAlbumView_AlbumRenamed/it_should_preserve_covers_on_folder_change` — seed canonical covers;
    rename with new folder; assert the new owner + visitor rows carry the preserved covers.
- Acceptance suite (`01-09`) extended with one case per operation: after `DeleteAlbum`, the covers
  record is gone; after a folder-change `RenameAlbum`, `ListAlbums` returns the renamed album with
  the preserved covers.

### Documentation

- `DATA_MODEL.md`: add a one-line note on the `#COVERS` row that the item is deleted on album
  deletion and moved to the new SK on folder-change rename.

## Out of scope

- Any event-payload change (both events stay as they are).
- In-place rename cover handling (no cover change by construction).
- Any new behaviour on `AlbumRenamed` for the view's display fields (that's already correct from
  Phase 1).

## References

- `../spec.md` (Automatic cover maintenance → album deleted / renamed).
- `../design.md` (Cover canonical record; Event-payload matrix).
- Load skills: `go`, `architecture`.
