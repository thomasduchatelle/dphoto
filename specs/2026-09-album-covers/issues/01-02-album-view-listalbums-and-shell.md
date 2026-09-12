# 01-02 — New AlbumView shell and single-query ListAlbums

Status: ready
Phase: 1
Layer: catalog domain — `pkg/catalogviews`
Depends on: 01-01

## Description

Replace the `Provider`-based `AlbumView` with a single struct that:

- serves `ListAlbums` from the projection alone (one `ListSummariesForUser` Query, plus the owner's
  sharing-grid Query — unchanged), and
- exposes an event method per domain event (`AlbumCreated`, `AlbumRenamedInPlace`, `AlbumDatesAmended`,
  `AlbumDeleted`, `AlbumShared`, `AlbumUnshared`, `MediasInserted`, `MediasTransferred`).

In this ticket the event methods are **stubbed** to `return nil`. Each subsequent ticket
(`01-03` … `01-07`) fills in one or two of them together with its adapter into the calling domain.

See parent ticket `01-album-list-read-model.md` for the full interface contract.

## Acceptance criteria

- New `AlbumView` struct with the interface agreed in the parent ticket:
  - `ListAlbums(ctx, user, filter) ([]*VisibleAlbum, error)`.
  - Stubbed event methods: `AlbumCreated`, `AlbumRenamedInPlace`, `AlbumDatesAmended`, `AlbumDeleted`,
    `AlbumShared`, `AlbumUnshared`, `MediasInserted`, `MediasTransferred`.
- `ListAlbums` uses **exactly one** `AlbumSummaryRepository.ListSummariesForUser` call for both owned and
  shared albums. When `user.Owner != nil` it decorates the owned rows with
  `GetAlbumSharingGrid(*user.Owner)` (unchanged).
- The `OnlyDirectlyOwned` filter is honoured by filtering `Availability.AsOwner == true` on the
  in-memory result — no repository-level change.
- Sort order unchanged: reverse chronological on `Start`, then on `End`.
- `NewAlbumView(...)` constructor updated: takes the repository, the sharing-grid port, and the
  media-counter port + find-albums-by-ids port (kept as fields so the event methods can be filled in by
  the next tickets without further signature churn).
- Deleted: `albums_view_owned_provider.go`, `albums_view_shared_provider.go`,
  `albums_view_provider.go`, `MediaCounterInjector`, `MediaCounterFromView`,
  `OwnedAlbumListProvider`, `SharedAlbumListProvider`. Any port that becomes unused
  (`FindAlbumByOwnerPort` on the read path) is removed.
- Factory `pkg/pkgfactory/factory_catalog_view.go::AlbumView` updated to the new constructor.
- Tests in `albums_view_test.go` are rewritten around the new shape:
  - `TestAlbumView_ListAlbums/it_should_serve_all_display_fields_from_the_view` — seed the in-memory
    fake with owner + visitor rows carrying `Name/Start/End/Count`, assert the full `VisibleAlbum` slice.
  - `TestAlbumView_ListAlbums/it_should_only_return_directly_owned_when_filter_is_set`.
  - `TestAlbumView_ListAlbums/it_should_decorate_owned_rows_with_the_sharing_grid`.
  - `TestAlbumView_ListAlbums/it_should_return_empty_when_the_user_has_no_rows`.
  - `TestAlbumView_ListAlbums/it_should_sort_by_start_desc_then_end_desc`.
  - The old `TestNewAlbumViewAcceptance` (owned/shared provider aggregation) is deleted with a one-line
    comment on the removal commit — it encodes the split we're eliminating.
- All existing callers of `AlbumView` (the `list-albums` lambda) continue to compile and work; the
  DTO mapping stays unchanged.

## Out of scope

- Any actual event-handling logic (the eight event methods are `return nil` stubs, wired by later
  tickets).
- Repository-level changes (already delivered by `01-01`).
- `CommandHandlerAlbumSize` deletion — it is left in place until `01-09` since it still observes the
  catalog events for now (its methods will be superseded one by one).
- Drift reconciler changes (`01-08`).

## References

- `../spec.md`, `../design.md`
- Parent ticket `01-album-list-read-model.md` (interface contract).
- Load skills: `go`, `architecture`.
