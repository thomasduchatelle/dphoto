# 01-03 — Event: AlbumCreated

Status: done
Phase: 1
Layer: catalog domain — `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 01-01, 01-02

## Description

Fill in the `AlbumView` handler for the `catalog.AlbumCreated` event so that creating an album
writes a full summary row for the owner (the only viewer at creation time), and updates the counts
on any albums whose medias were transferred into the new one at creation time.

The `catalog.AlbumCreated` event carries both `CreatedAlbum` and `TransferredMedias`, so the same
handler owns both the "create the new row" step and the "recount the source albums" step. This
keeps the projection consistent within one observer call and removes the need for a separate
`MediasTransferred` shell method on `AlbumView`.

The old `CommandHandlerAlbumSize` observer stays wired for now — it is retired in `01-09`.

## Acceptance criteria

- `AlbumView` exposes `AlbumCreated(ctx, event catalog.AlbumCreated) error` (replacing the
  01-02 stub `AlbumCreated(ctx, catalog.Album)`; adjust the shell signature accordingly). The
  method:
  - Writes a full summary row for the owner via `Repository.PutSummaries` (or `SetDisplayFields`
    upsert — implementation choice, whichever composes best with the 01-01 primitives).
    `Name/Start/End` come from `event.CreatedAlbum`, `Count=0`, `Availability =
    OwnerAvailability(user of the owner)`. If the album already has visitors at creation time
    (edge case, none expected today), a row is written for each.
  - If `event.TransferredMedias` is non-empty, recount the affected source albums (from
    `event.TransferredMedias.FromAlbums`) via `MediaCounterPort.CountMedia` and persist the new
    counts via `Repository.SetCounts` — SET Count = :c, display fields untouched. The destination
    album's count reflects the medias just inserted; since the row is being freshly created here,
    the count is set via the `PutSummaries` above (compute it from `event.TransferredMedias` or
    ask the port — implementation choice).
- The owner's `usermodel.UserId` is derived from the album's owner via
  `ListUsersWhoCanAccessAlbumPort` (kept as a port on `AlbumView`).
- New adapter `AlbumViewCreateAlbumObserver` in `pkg/catalogviews` (or a dedicated `adapters.go`)
  implementing `catalog.AlbumCreatedObserver` (`OnAlbumCreated(ctx, event)`) by forwarding to
  `AlbumView.AlbumCreated`.
- Factory `SimpleCatalogFactory.CreateAlbumCase` wires the new adapter alongside the current
  `CommandHandlerAlbumSize`-based observers (`CommandHandlerAlbumSize` is retired in `01-09`).
- `AlbumView` tests (against the in-memory `AlbumSummaryRepository` fake):
  - `TestAlbumView_AlbumCreated/it_should_make_the_album_visible_to_the_owner` — dispatch
    `AlbumCreated` with an empty `TransferredMedias`, then `ListAlbums`, assert the album is
    returned with all display fields and `MediaCount=0`.
  - `TestAlbumView_AlbumCreated/it_should_not_shadow_an_existing_row` — pre-seed a row with
    `Count=5`; dispatch again with empty `TransferredMedias`; expectation depends on the
    primitive chosen (idempotent-preserving-count if `SetDisplayFields`, overwritten if
    `PutSummaries` — either is acceptable, documented in the PR).
  - `TestAlbumView_AlbumCreated/it_should_recount_transferred_source_albums` — pre-seed a
    source album row with `Count=10`; dispatch `AlbumCreated` with a `TransferredMedias`
    moving 3 medias out of that source; assert the source count reflects `MediaCounterPort`
    and its display fields are untouched.

## Out of scope

- Rename / amend-dates / delete / share / unshare (each in their own ticket).
- Direct `MediasInserted` events (`01-04`).
- Drift changes (`01-08`).
- Deletion of `CommandHandlerAlbumSize` (`01-09`).

## References

- `../spec.md`, `../design.md`
- Parent ticket `01-album-list-read-model.md` (interface contract).
- Load skills: `go`, `architecture`.
