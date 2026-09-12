# 01-04 — Events: MediasInserted & MediasTransferred

Status: ready
Phase: 1
Layer: catalog domain — `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 01-01, 01-02, 01-03

## Description

Fill in the two count-management event methods and wire them into the catalog / backup flows
through adapters. Both methods must operate on `Count` alone so that they never clobber the
display fields set by `01-03` and the tickets that follow.

- `MediasInserted` is the atomic per-insert diff (currently `CommandHandlerAlbumSize.OnMediasInserted`).
- `MediasTransferred` is the recount fired when medias move between albums as a side-effect of
  create-with-transfer / rename-with-folder-change / amend-dates / delete (currently
  `CommandHandlerAlbumSize.OnTransferredMedias`).

## Acceptance criteria

- `AlbumView.MediasInserted(ctx, medias map[catalog.AlbumId][]catalog.MediaId) error` looks up viewers
  via `ListUsersWhoCanAccessAlbumPort` and calls `Repository.IncrementCounts` — atomic `ADD Count :d`,
  no other attribute touched. `d = len(medias[albumId])`.
- `AlbumView.MediasTransferred(ctx, transfers catalog.TransferredMedias) error`:
  - Collects the touched albums (destinations from `transfers.Transfers` + sources from
    `transfers.FromAlbums`).
  - Re-queries the canonical count via `MediaCounterPort.CountMedia`.
  - Writes the new counts via `Repository.SetCounts` (SET Count = :c), **not** via `PutSummaries`,
    so display fields survive the operation.
  - On an album whose row does not yet exist (source album gone empty and never had a row), the
    `SetCounts` upsert may leave a row with empty display fields. Acceptable in this ticket because
    the album is either being deleted (handled by `01-06`) or already had display fields written by
    `01-03`.
- Two adapters in `pkg/catalogviews` (or `adapters.go`):
  - `AlbumViewMediasInsertedObserver` — matches the port that `catalog.InsertMedias` calls today
    (currently `CommandHandlerAlbumSize.OnMediasInserted`).
  - `AlbumViewTimelineMutationObserver` — implements `catalog.TimelineMutationObserver` by forwarding
    to `AlbumView.MediasTransferred`.
- Factories wire both adapters alongside (or replacing) the existing `CommandHandlerAlbumSize`
  observer in: `InsertMediasCase`, `CreateAlbumCase`, `CreateAlbumDeleteCase`, `RenameAlbumCase`,
  `AmendAlbumDatesCase`, and the backup factory.
- `AlbumView` tests (in-memory fake):
  - `TestAlbumView_MediasInserted/it_should_increment_the_count_for_all_viewers` — pre-seed rows for
    owner + visitor, insert 3 medias, assert both `Count` bumped by 3, display fields intact.
  - `TestAlbumView_MediasTransferred/it_should_recount_the_source_and_destination_albums` — seed
    two albums with different counts, transfer, assert counts reflect `MediaCounterPort`, display
    fields untouched.
  - `TestAlbumView_MediasTransferred/it_should_not_clobber_display_fields` — seed with
    `Name/Start/End`, run transfer, assert display fields still present.

## Out of scope

- Rename / amend-dates / delete / share (their own tickets).
- Anything about videos vs images (the count is over all medias, unchanged from today).
- Drift (`01-08`).

## References

- `../spec.md`, `../design.md`
- Parent ticket `01-album-list-read-model.md` (interface contract).
- Load skills: `go`, `architecture`.
