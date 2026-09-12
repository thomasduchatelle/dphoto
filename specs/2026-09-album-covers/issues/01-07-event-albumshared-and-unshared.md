# 01-07 — Events: AlbumShared & AlbumUnshared

Status: ready
Phase: 1
Layer: catalog domain — `pkg/catalogviews` + `pkg/acl/catalogacl` + `pkg/pkgfactory`
Depends on: 01-01, 01-02

## Description

Fill in the two sharing-related events. `AlbumShared` must write a *complete* visitor row (display
fields + current count) so that the visitor's next `ListAlbums` renders the card without any extra
lookup.

To avoid an extra DDB round-trip in the handler, change the signature of
`catalogacl.AlbumSharedObserver` to receive the whole `catalog.Album` — `ShareAlbumCase` already
looks it up before notifying observers, so this is free.

## Acceptance criteria

- `AlbumView.AlbumShared(ctx, album catalog.Album, userId usermodel.UserId) error`:
  - Reads the current count via `MediaCounterPort.CountMedia`.
  - Writes one full visitor row via `Repository.PutSummaries` — `Name/Start/End` from `album`,
    `Count` from the port, `Availability = VisitorAvailability(userId)`.
- `AlbumView.AlbumUnshared(ctx, albumId catalog.AlbumId, userId usermodel.UserId) error` calls
  `Repository.DeleteRow(VisitorAvailability(userId), albumId)`.
- `catalogacl.AlbumSharedObserver` interface changes:
  - Old: `AlbumShared(ctx, albumId catalog.AlbumId, userEmail usermodel.UserId) error`.
  - New: `AlbumShared(ctx, album catalog.Album, userEmail usermodel.UserId) error`.
  - `ShareAlbumCase.ShareAlbumWith` updated to pass the `*catalog.Album` already fetched at the top
    of the case (`FindAlbumPort.FindAlbum(ctx, albumId)`).
- `catalogacl.AlbumUnSharedObserver` interface unchanged (albumId + user is enough).
- Two adapters in `pkg/catalogviews`:
  - `AlbumViewAlbumSharedObserver` implementing the **new** `catalogacl.AlbumSharedObserver`.
  - `AlbumViewAlbumUnSharedObserver` implementing `catalogacl.AlbumUnSharedObserver`.
- Factory `pkg/pkgfactory/factory_acl.go` updated to wire the new adapters (replacing the current
  `CommandHandlerAlbumSize` observer for these two events).
- Tests:
  - `pkg/acl/catalogacl/case_share_album_test.go` updated: `AlbumSharedObserverFake` signature
    changed and assertions verify the observer receives the full `Album`, not just its id.
  - `AlbumView` tests (in-memory fake):
    - `TestAlbumView_AlbumShared/it_should_add_a_visitor_row_with_display_fields_and_count` — seed
      no rows, share, assert `ListAlbums` for the visitor returns the album with `Name/Start/End`
      and the correct count from `MediaCounterPort`.
    - `TestAlbumView_AlbumUnshared/it_should_remove_the_visitor_row_only` — seed owner + visitor,
      unshare, assert the visitor row is gone and the owner row untouched.

## Out of scope

- Reverse-reader (`ListUsersWhoCanAccessAlbum`) semantics — unchanged.
- Drift.
- Deletion of `CommandHandlerAlbumSize.AlbumShared`/`AlbumUnShared` (they become dead once this
  ticket is wired; final removal happens in `01-09`).

## References

- `../spec.md`, `../design.md`
- Parent ticket `01-album-list-read-model.md` (interface contract).
- Load skills: `go`, `architecture`.
