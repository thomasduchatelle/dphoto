# 01-06 — Event: AlbumDeleted

Status: ready
Phase: 1
Layer: catalog domain — `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 01-01, 01-02

## Description

Fill in `AlbumView.AlbumDeleted` so that deleting an album wipes every viewer row of that album from
the projection, and wire it into `catalog.DeleteAlbum` via an adapter.

This closes a pre-existing gap: today the owner's summary row survives an album deletion (see the
`TODO Everything about the album should be deleted if the album is deleted` on
`pkg/catalogviews/command_handler_counts.go:196`). Only the visitor rows were removed indirectly by
the unshare path.

## Acceptance criteria

- `AlbumView.AlbumDeleted(ctx, albumId catalog.AlbumId) error` removes every row for `albumId`,
  regardless of `Availability`, via `Repository.DeleteAllRowsForAlbum` (which itself may fan out from
  `ListUsersWhoCanAccessAlbum` or scan the projection — repository's choice, delivered in `01-01`).
- New adapter `AlbumViewDeleteAlbumObserver` implementing `catalog.DeleteAlbumObserver` by forwarding
  to `AlbumView.AlbumDeleted`.
- Factory `SimpleCatalogFactory.CreateAlbumDeleteCase` wires the adapter alongside the existing
  observers.
- `AlbumView` tests (in-memory fake):
  - `TestAlbumView_AlbumDeleted/it_should_remove_the_owner_row` — seed the owner row, delete, assert
    `ListAlbums` returns no album.
  - `TestAlbumView_AlbumDeleted/it_should_remove_all_visitor_rows` — seed owner + two visitors,
    delete, assert all three rows are gone.
  - `TestAlbumView_AlbumDeleted/it_should_leave_other_albums_untouched` — seed two albums for the
    same user, delete one, assert the other is still returned.
- The `TODO` on line 196 of `command_handler_counts.go` is removed as part of this ticket (or in
  `01-09`, whichever comes first).

## Out of scope

- Anything about medias transferred out of the album by the delete flow (already handled by `01-04`'s
  `MediasTransferred`).
- Rename with folder change — its "delete the old row" step is fulfilled by this ticket.
- Drift.

## References

- `../spec.md`, `../design.md`
- Parent ticket `01-album-list-read-model.md` (interface contract).
- Load skills: `go`, `architecture`.
