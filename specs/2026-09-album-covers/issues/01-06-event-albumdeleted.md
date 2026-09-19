# 01-06 — Event: AlbumDeleted

Status: ready
Phase: 1
Layer: catalog domain — `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 01-01, 01-02

## Description

Fill in the `AlbumView` handler for the `catalog.AlbumDeleted` event so that deleting an album
wipes every viewer row of that album from the projection, and updates the counts on any albums
whose medias were transferred out of the deleted one.

The `catalog.AlbumDeleted` event carries both `DeletedAlbumId` and `TransferredMedias`, so the
same handler owns both the "wipe the rows" step and the "recount the destination albums" step.

This also closes a pre-existing gap: before this feature the owner's summary row survived an
album deletion (see the `TODO Everything about the album should be deleted if the album is
deleted` on `pkg/catalogviews/command_handler_counts.go:196`). Only the visitor rows were removed
indirectly by the unshare path.

## Acceptance criteria

- `AlbumView.AlbumDeleted(ctx, event catalog.AlbumDeleted) error` (replacing the 01-02 stub
  `AlbumDeleted(albumId)`; adjust the shell signature accordingly):
  - Removes every row for `event.DeletedAlbumId` regardless of `Availability` via
    `Repository.DeleteAllRowsForAlbum` (delivered by `01-01`).
  - If `event.TransferredMedias` is non-empty, recount the destination albums (from
    `event.TransferredMedias.Transfers`) via `MediaCounterPort.CountMedia` and persist the new
    counts via `Repository.SetCounts` — SET Count = :c, display fields untouched.
- New adapter `AlbumViewDeleteAlbumObserver` implementing `catalog.AlbumDeletedObserver`
  (`OnAlbumDeleted(ctx, event)`) by forwarding to `AlbumView.AlbumDeleted`.
- Factory `SimpleCatalogFactory.CreateAlbumDeleteCase` wires the adapter alongside the existing
  `CommandHandlerAlbumSize`-based observers (`CommandHandlerAlbumSize` retired in `01-09`).
- `AlbumView` tests (in-memory fake):
  - `TestAlbumView_AlbumDeleted/it_should_remove_the_owner_row` — seed the owner row, dispatch
    with empty `TransferredMedias`, assert `ListAlbums` returns no album.
  - `TestAlbumView_AlbumDeleted/it_should_remove_all_visitor_rows` — seed owner + two visitors,
    dispatch, assert all three rows are gone.
  - `TestAlbumView_AlbumDeleted/it_should_leave_other_albums_untouched` — seed two albums for
    the same user, delete one, assert the other is still returned.
  - `TestAlbumView_AlbumDeleted/it_should_recount_transferred_destination_albums` — seed
    destination album with `Count=2`, delete with `TransferredMedias` moving 3 medias into that
    destination, assert destination count reflects `MediaCounterPort` and display fields
    survive.
- The `TODO` on line 196 of `command_handler_counts.go` is removed as part of this ticket (or in
  `01-09`, whichever comes first).

## Out of scope

- Rename / amend-dates / share / drift.
- Direct `MediasInserted` events.
- Deletion of `CommandHandlerAlbumSize` (`01-09`).

## References

- `../spec.md`, `../design.md`
- Parent ticket `01-album-list-read-model.md` (interface contract).
- Load skills: `go`, `architecture`.
