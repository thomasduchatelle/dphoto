# 01-05 — Events: AlbumRenamed & AlbumDatesAmended

Status: ready
Phase: 1
Layer: catalog domain — `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 01-01, 01-02

## Description

Fill in the two display-field-mutation event handlers and wire them into the corresponding
catalog flows.

The `catalog.AlbumRenamed` event was reshaped (base branch refactor) to carry `ExistingAlbum`,
`RenamedAlbum`, and `TransferredMedias`. It fires on **both** the in-place branch (name only,
folder unchanged) and the folder-changing branch of `catalog.RenameAlbum`. The `AlbumView`
handler decides what to do based on whether the folder changed:

- name-only change → `SetDisplayFields` on the existing row.
- folder change → treat as a "delete old row + create new row" (the semantics are equivalent to
  `AlbumDeleted(ExistingAlbum.AlbumId) + AlbumCreated(RenamedAlbum)`), plus recount the transfer
  destinations/sources.

`catalog.AlbumDatesAmended` similarly carries `DatesUpdate` and `TransferredMedias`, so the
handler both updates the display fields and recounts the affected albums when medias moved.

Both events must operate on display attributes / counts as separate operations so display fields
and count are never clobbered by a single write.

## Acceptance criteria

- `AlbumView.AlbumRenamed(ctx, event catalog.AlbumRenamed) error` (replacing the 01-02 stub
  `AlbumRenamedInPlace(albumId, newName)`; adjust the shell signature accordingly):
  - When `event.ExistingAlbum.AlbumId == event.RenamedAlbum.AlbumId` (folder unchanged): look up
    viewers via `ListUsersWhoCanAccessAlbumPort` and call `Repository.SetDisplayFields` with
    `Name/Start/End` from `event.RenamedAlbum`. Start/End must be preserved / updated correctly
    even if only the name changed.
  - When the folder changed: `Repository.DeleteAllRowsForAlbum(event.ExistingAlbum.AlbumId)` and
    write a full row for `event.RenamedAlbum` for every viewer (same logic as `01-03`
    `AlbumCreated`). The count of the new row is derived from `event.TransferredMedias` or
    queried via `MediaCounterPort.CountMedia`.
  - If `event.TransferredMedias` affects other albums (sources), recount them via
    `Repository.SetCounts`.
- `AlbumView.AlbumDatesAmended(ctx, event catalog.AlbumDatesAmended) error` (replacing the 01-02
  stub `AlbumDatesAmended(update)`; adjust the shell signature accordingly):
  - Look up viewers, call `Repository.SetDisplayFields` with `Name`, and the new `Start/End` from
    `event.DatesUpdate.UpdatedAlbum` (or the equivalent field on the event).
  - If `event.TransferredMedias` is non-empty, recount the affected source/destination albums via
    `Repository.SetCounts`.
- Two adapters in `pkg/catalogviews`:
  - `AlbumViewRenameObserver` implementing `catalog.AlbumRenamedObserver`
    (`OnAlbumRenamed(ctx, event)`) by forwarding to `AlbumView.AlbumRenamed`.
  - `AlbumViewAmendDatesObserver` implementing `catalog.AlbumDatesAmendedObserver`
    (`OnAlbumDatesAmended(ctx, event)`) by forwarding to `AlbumView.AlbumDatesAmended`.
- Factories: `RenameAlbumCase` and `AmendAlbumDatesCase` in `pkg/pkgfactory/factory_catalog.go`
  updated to pass the new adapters alongside the existing `CommandHandlerAlbumSize`-based
  observers (`CommandHandlerAlbumSize` retired in `01-09`).
- `AlbumView` tests (in-memory fake):
  - `TestAlbumView_AlbumRenamed/it_should_update_the_name_on_all_viewer_rows_when_folder_unchanged`
    — seed owner + visitor, rename with `ExistingAlbum.AlbumId == RenamedAlbum.AlbumId`, assert
    both rows show the new name and same Start/End/Count.
  - `TestAlbumView_AlbumRenamed/it_should_replace_the_rows_when_folder_changes` — seed owner
    row for the old albumId, rename with new folder, assert the old row is gone and a new row
    exists for the new albumId with correct display fields.
  - `TestAlbumView_AlbumDatesAmended/it_should_update_start_and_end_on_all_viewer_rows` — seed
    owner + visitor, amend dates without transfer, assert both rows show new dates and same
    Name/Count.
  - `TestAlbumView_AlbumDatesAmended/it_should_recount_transferred_albums` — seed source album
    with `Count=5`, amend dates with `TransferredMedias` moving medias out, assert source count
    reflects `MediaCounterPort` and display fields survive.

## Out of scope

- Direct `MediasInserted` (`01-04`).
- Delete / share / unshare / drift.
- Deletion of `CommandHandlerAlbumSize` (`01-09`).

## References

- `../spec.md`, `../design.md`
- Parent ticket `01-album-list-read-model.md` (interface contract).
- Load skills: `go`, `architecture`.
