# 09 — Carry covers on `MediasInserted` + Randomize primitive upgrade

Status: ready
Phase: 2
Layer: catalog domain — `pkg/catalog` + `pkg/catalogadapters/catalogdynamo` + `pkg/catalogviews` + `pkg/catalogviewsadapters` + `pkg/pkgfactory` + `cmd/dphotops`
Depends on: —

## Description

Two intertwined changes shipped in one vertical slice:

1. **Upgrade the Randomize primitive**: rename `CompleteCovers` →
   `RandomizeCovers`, change its semantics to "drop every `RANDOM`, keep every `CHERRY_PICKED`, fill
   empties up to 4", and make both methods return the resulting `[]Cover` so callers can attach it
   to lifecycle events. The admin backfill (`BackfillCovers` + `dphotops covers backfill`) switches
   to the new names and gets a CLI help-text refresh.
2. **Propagate covers through `MediasInserted`**: the media-insert use case calls the upgraded
   primitive per affected album with the just-inserted images as candidates, emits the
   `MediasInserted` event carrying the resulting covers per album, and `AlbumView.OnMediasInserted`
   denormalises covers alongside the count update in a single per-row write.

The primitive is bundled with the first consumer so there is no standalone schema-only ticket; after
this ticket lands, issues 10 and 11 can be worked in parallel against the upgraded primitive.

The backup domain (`pkg/backup`, `cmd/dphoto`) is **not** touched. This replaces issue 05
(`wontdo`).

See `../design.md` → Randomize operation; Event-payload matrix (row `MediasInserted`); §Cover
projection in the view.

## Acceptance criteria

### Randomize primitive upgrade

- Rename in `pkg/catalog`:
  - `CompleteCovers` struct → `RandomizeCovers`.
  - `NewCompleteCovers(...)` → `NewRandomizeCovers(...)`.
  - `CompleteCovers.CompleteCovers(ctx, albumId) error` → `RandomizeCovers.RandomizeCovers(ctx, albumId) ([]Cover, error)`.
  - `CompleteCovers.CompleteCoversFromCandidates(ctx, albumId, candidates) error` →
    `RandomizeCovers.RandomizeCoversFromCandidates(ctx, albumId, candidates) ([]Cover, error)`.
  - `CompleteCoversPort` interface → `RandomizeCoversPort`, both methods now return `([]Cover, error)`.
- Behaviour upgrade (both methods):
  - Drop every `RANDOM` cover before filling (today's version keeps them).
  - Keep every `CHERRY_PICKED` cover untouched.
  - Compute `slotsToFill = 4 − len(kept)` and fill from eligible candidates (`MediaType IMAGE`, not
    already a kept cover) uniformly at random, tagged `RANDOM`.
  - Return the resulting ordered `[]Cover`.
- `BackfillCovers` keeps its name, retargets to `RandomizeCoversPort`, and updates its method call.
- Factory `pkg/pkgfactory/factory_catalog.go::CompleteCoversCase` renamed to `RandomizeCoversCase`.
- CLI help text in `cmd/dphotops/cmd/covers.go` updated: re-running `dphotops covers backfill` may
  **change** existing `RANDOM` covers; `CHERRY_PICKED` ones remain preserved.
- Tests in `pkg/catalog/covers_test.go` rewritten against the new names/signatures, plus:
  - `it_should_drop_existing_RANDOM_covers_and_redraw`.
  - `it_should_keep_only_the_kept_set_when_no_eligible_candidate`.
  - `it_should_be_a_noop_when_the_kept_set_is_already_full_of_cherry_picked`.
  - `it_should_return_the_resulting_set_from_both_methods`.

### MediasInserted propagation

- `pkg/catalog`:
  - The media-insert use case (`pkg/catalog/medias_insert.go`) gains a dependency on
    `RandomizeCoversPort`.
  - After the medias are inserted and the indexing side-effects complete, for each affected
    `AlbumId` the use case calls `RandomizeCoversPort.RandomizeCoversFromCandidates(ctx, albumId,
    insertedMediasForThatAlbum)` and collects the returned `[]Cover` per album.
  - The `MediasInserted` event payload (whatever shape the current `InsertMediasObserver` expects —
    today a `map[AlbumId][]MediaId`) is extended with `Covers map[AlbumId][]catalog.Cover` carrying
    the result for every album whose covers were (re)computed. Albums whose cover set is unchanged
    (e.g. randomize was a no-op because the set was already full of `CHERRY_PICKED`) may be omitted
    or carry the unchanged set — document the choice in the PR.
  - Order of operations is: insert medias → randomize covers → notify observers (count + covers in
    one event).
- `pkg/catalogviews`:
  - `AlbumView.OnMediasInserted` method signature updated to accept the covers map alongside the
    inserts.
  - For every affected album, viewers are fanned out via `ListUsersWhoCanAccessAlbumPort`; the
    implementation performs a single per-row `UpdateItem` that does
    `ADD Count :d SET Covers = :c` — attribute footprints stay disjoint from display fields.
  - A new repository method on `AlbumSummaryRepository` is introduced:
    `IncrementCountAndSetCoversForAllViewers(ctx, []AlbumCountAndCoversDiff) error`, with
    `AlbumCountAndCoversDiff{AvailabilityType, AlbumId, Diff int, Covers []Cover}`. The existing
    `IncrementCountForAllViewers` is kept for callers that have no covers to propagate (none in
    production after this ticket, but kept available for drift / tests).
  - Both the in-memory fake (`AlbumSummaryInMemoryRepository`) and the DynamoDB adapter implement
    the new method.
  - DynamoDB adapter: single `UpdateItem` per row, expression
    `ADD Count :d SET Covers = :c`. Covers is set unconditionally (empty list is a valid write).
- `pkg/pkgfactory`:
  - `InsertMediasCase` wires the `RandomizeCoversPort` into the use case.
  - The AlbumView adapter (`AlbumViewMediasInsertedObserver`) forwards the new event shape to the
    updated `AlbumView.OnMediasInserted`.
- `pkg/backup` and `cmd/dphoto` are **not changed**. Grep the diff — any touched file in these paths
  is a bug in the PR.
- Tests:
  - `pkg/catalog/medias_insert_test.go` extended:
    - `it_should_randomize_covers_for_each_affected_album` — insert medias into 2 albums with empty
      covers; assert `RandomizeCoversPort` was called once per album with the correct candidates;
      assert the event carries covers per album.
    - `it_should_leave_full_cover_sets_unchanged_but_still_report_them` — insert medias into an
      album whose set is full of `CHERRY_PICKED`; assert the event's `Covers[albumId]` equals the
      pre-existing set (or is omitted — PR documents choice).
  - `pkg/catalogviews/albums_view_events_test.go`:
    - `TestAlbumView_MediasInserted/it_should_increment_count_and_set_covers_in_one_write_per_viewer`
      — seed owner + visitor with `Count=1` and 1 cover; dispatch with 3 medias inserted and a new
      4-cover set; assert both rows show `Count=4` and the new covers, display fields intact.
    - `TestAlbumView_MediasInserted/it_should_do_nothing_when_the_event_is_empty`.
  - DynamoDB adapter test for `IncrementCountAndSetCoversForAllViewers`:
    - Seed a row with `Count=5` and display fields; call with `Diff=2` and new covers; assert
      `Count=7`, covers updated, display fields untouched.
- The `pkg/` acceptance suite (`01-09`) extended: after a backup/insert of 3 images into an empty
  album, `ListAlbums` returns the album with `MediaCount=3` **and** 3 covers.
- **Dead-code cleanup — conditional**: `AlbumView.OnAlbumCoversChanged` is dead once issues 09, 10
  and 11 have all landed. If this is the last of the three to merge
  (`grep OnAlbumCoversChanged` shows no production caller), delete the method and its tests from
  `pkg/catalogviews/albums_view.go` and `pkg/catalogviews/albums_view_events_test.go` as part of
  this PR. The repository method `SetCoversForAllViewers` is kept.

## Out of scope

- `AlbumCreated`, `AlbumDatesAmended`, `AlbumShared` — see issues 10, 11.
- Lifecycle maintenance for `AlbumDeleted`, `AlbumRenamed` folder-change — issue 12.

## References

- `../spec.md` (Automatic cover maintenance → medias added).
- `../design.md` (Event-payload matrix; Cover projection in the view; §Writing covers alongside counts).
- Load skills: `go`, `architecture`.
