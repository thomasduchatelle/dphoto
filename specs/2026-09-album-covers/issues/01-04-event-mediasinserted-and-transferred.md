# 01-04 — Event: MediasInserted

Status: ready
Phase: 1
Layer: catalog domain — `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 01-01, 01-02

## Description

Fill in the count-management path for direct media inserts (i.e. medias appended to an existing
album, not moved by a use-case such as create-with-transfer or amend-dates). The method operates
on `Count` alone so it never clobbers the display fields set by the other tickets.

The transfer-driven count updates that were originally in this ticket's scope have moved into the
per-use-case handlers (`01-03` for create, `01-05` for rename & amend-dates, `01-06` for delete)
because each `catalog.Album*` event now carries the `TransferredMedias` field. The
`AlbumView.MediasTransferred` shell stub from `01-02` is therefore removed.

## Acceptance criteria

- `AlbumView` exposes `MediasInserted(ctx, event catalog.MediasInsertedEvent) error` — signature
  aligned to `catalog.InsertMediasObserver` (use whatever event type / arguments that observer
  currently defines; e.g. `map[catalog.AlbumId][]catalog.MediaId` if that's still the shape).
  Implementation:
  - For each `albumId` in the event, look up viewers via `ListUsersWhoCanAccessAlbumPort` and
    build an `AlbumMediaCountDiff{AvailabilityType, AlbumId, Diff: len(medias)}` per viewer.
  - Call `Repository.IncrementCounts(diffs)` — atomic `ADD Count :d`, no other attribute touched.
- New adapter `AlbumViewMediasInsertedObserver` in `pkg/catalogviews` implementing
  `catalog.InsertMediasObserver` (`OnMediasInserted`), forwarding to `AlbumView.MediasInserted`.
- Factory `SimpleCatalogFactory.InsertMediasCase` (and any other factory that wires the current
  `CommandHandlerAlbumSize.OnMediasInserted`) is updated to pass the new adapter alongside the
  existing `CommandHandlerAlbumSize` observer.
- Delete the stub `AlbumView.MediasTransferred` shell method left over from `01-02`. Any factory
  wiring that referenced it is cleaned up. The transfer-driven recount is handled per-use-case
  by the tickets `01-03`, `01-05`, `01-06`.
- `AlbumView` tests (in-memory fake):
  - `TestAlbumView_MediasInserted/it_should_increment_the_count_for_all_viewers` — pre-seed rows
    for owner + visitor with `Count=1` and display fields; insert 3 medias; assert both rows show
    `Count=4` and display fields intact.
  - `TestAlbumView_MediasInserted/it_should_do_nothing_when_the_event_is_empty` — dispatch with
    an empty map; assert repository was not called (or called with an empty diff list — either
    is acceptable, PR body documents the choice).

## Out of scope

- Create / rename / amend-dates / delete / share (their own tickets, each handling its own
  transfer counts).
- Anything about videos vs images (the count is over all medias, unchanged from today).
- Drift (`01-08`).

## References

- `../spec.md`, `../design.md`
- Parent ticket `01-album-list-read-model.md` (interface contract).
- Load skills: `go`, `architecture`.
