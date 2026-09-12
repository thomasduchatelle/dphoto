# 01-03 — Event: AlbumCreated

Status: ready
Phase: 1
Layer: catalog domain — `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 01-01, 01-02

## Description

Fill in `AlbumView.AlbumCreated` so that creating an album writes a full summary row for the owner
(the only viewer at creation time), and wire it into `catalog.CreateAlbum` through an adapter that
implements `catalog.CreateAlbumObserver`.

The old `CommandHandlerAlbumSize` observer stays wired for now but no longer needs to be the one
handling this event — the create-time write of the summary row moves to `AlbumView.AlbumCreated`.

## Acceptance criteria

- `AlbumView.AlbumCreated(ctx, album catalog.Album) error` writes a full summary row for the owner
  via `Repository.PutSummaries` (or `SetDisplayFields` upsert, whichever composes best with the
  existing `01-01` primitives): `Name/Start/End` populated from `album`, `Count=0`, `Availability =
  OwnerAvailability(user of the owner)`. If the album has visitors at creation time (edge case, none
  expected today), a row is written for each.
- The owner's `usermodel.UserId` is derived from the album's owner via
  `ListUsersWhoCanAccessAlbumPort` (kept as a `AlbumView` port), so the code stays honest about the
  owner→user mapping.
- New adapter `AlbumViewCreateAlbumObserver` in `pkg/catalogviews` (or a dedicated `adapters.go`)
  implementing `catalog.CreateAlbumObserver` by forwarding to `AlbumView.AlbumCreated`.
- Factory `SimpleCatalogFactory.CreateAlbumCase` wires the new adapter alongside (or replacing) the
  current `CommandHandlerAlbumSize` observer.
- `AlbumView` test (against the in-memory `AlbumSummaryRepository` fake):
  `TestAlbumView_AlbumCreated/it_should_make_the_album_visible_to_the_owner` — call
  `AlbumCreated`, then `ListAlbums`, assert the album is returned with all display fields and
  `MediaCount=0`.
- Extra `AlbumView` test:
  `TestAlbumView_AlbumCreated/it_should_not_shadow_an_existing_row` — pre-seed a row with
  `Count=5`; create again is idempotent and does not zero the count (or, if the implementation
  uses `PutSummaries` and does overwrite, this expectation is inverted with a comment justifying it).
- Removed from `CommandHandlerAlbumSize`: nothing yet (it stops being the source of truth for this
  event but stays wired until `01-09`).

## Out of scope

- Rename / amend-dates / delete / share / unshare (each in their own ticket).
- Count events (`01-04`).
- Drift changes (`01-08`).
- Deletion of `CommandHandlerAlbumSize` (`01-09`).

## References

- `../spec.md`, `../design.md`
- Parent ticket `01-album-list-read-model.md` (interface contract).
- Load skills: `go`, `architecture`.
