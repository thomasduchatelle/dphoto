# 10 — Carry covers on `AlbumCreated` and `AlbumDatesAmended`

Status: ready
Phase: 2
Layer: catalog domain — `pkg/catalog` + `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 09 (for the upgraded `RandomizeCoversPort`)

## Description

Mirror the `MediasInserted` propagation (issue 09) for the two other lifecycle events whose outcome
can change covers: `AlbumCreated` (when `TransferredMedias` brings images into the new album) and
`AlbumDatesAmended` (when `TransferredMedias` moves medias between albums — the source may lose covers
that referenced the moved medias; the destination may need its empty slots filled).

Touches different files than 09 and 11, so these three issues can be worked **in parallel** after 09
lands the primitive. 10's AlbumView writes do **not** need the atomic
`ADD Count :d SET Covers = :c` method from 09 — amend-dates is a rare operation and the existing
display-field write + a second `SetCoversForAllViewers` write is fine (disjoint attributes, no
clobbering).

See `../design.md` → Event-payload matrix (rows `AlbumCreated`, `AlbumDatesAmended`), §Handling
transfers-out, and §Cover projection in the view.

## Acceptance criteria

### `AlbumCreated`

- `catalog.AlbumCreated` event payload gains `Covers []catalog.Cover` (always present; empty slice when
  nothing to show).
- Create-album use case:
  - If `TransferredMedias` brings images into the new album, call
    `RandomizeCoversPort.RandomizeCoversFromCandidates(ctx, newAlbumId, transferredImages)` and
    attach the resulting `[]Cover` to the event.
  - If no medias are transferred in, `event.Covers = []`.
  - The source albums impacted by the transfer-out are handled below (same semantics as the
    `AlbumDatesAmended` source handling).
- `AlbumView.OnAlbumCreated` writes the full row in one shot: `PutSummaries` with
  `Name/Start/End/Count/Covers` from the event.

### `AlbumDatesAmended`

- `catalog.AlbumDatesAmended` event payload gains `Covers map[catalog.AlbumId][]catalog.Cover`
  carrying the resulting cover set for **every album** affected by the transfer (the amended album and
  every source/destination touched).
- Amend-dates use case, per affected album:
  - **Destination (medias moved in)**: `RandomizeCoversPort.RandomizeCoversFromCandidates(ctx,
    albumId, movedInImages)` and attach the result.
  - **Source (medias moved out)**: load the current covers via `FindCoversByAlbum`; strip any cover
    whose `MediaId` is in the "moved out" set for this album; persist the pruned set via
    `CoverRepository.SaveCovers`; attach the pruned set to the event. **No re-randomisation on
    sources** — empties stay empty until the next time medias land in the album.
- `AlbumView.OnAlbumDatesAmended`:
  - Updates display fields on the amended album (as today).
  - For each entry in `event.Covers`, writes the new covers to every viewer row via the existing
    `SetCoversForAllViewers` repository method. Two writes per row (display fields + covers), on
    disjoint attributes — no clobbering concern.
  - Count recounts for transferred source/destination albums (unchanged — already done via
    `SetCounts` in issue 01-05).

### Transfer-out on `AlbumCreated` (same pattern as `AlbumDatesAmended` source)

- When `AlbumCreated.TransferredMedias` moves medias out of pre-existing source albums, the create
  use case must also strip stale covers from those sources and attach the pruned sets. Carry them via
  the event: a `SourceAlbumsCovers map[AlbumId][]Cover` field is one option; another is to put
  everything into a single `Covers map[AlbumId][]Cover` field (consistent with
  `AlbumDatesAmended`). **Decision: use the single map field in both events for consistency** — the
  new album's covers live under its own `AlbumId` key.
  - That makes `AlbumCreated.Covers` a `map[AlbumId][]Cover` rather than a flat `[]Cover`. Update
    the acceptance criteria for `AlbumCreated` above accordingly.
  - PR documents the final shape; the test matrix below assumes the map shape.

### Wiring

- `pkg/pkgfactory` updates `CreateAlbumCase` and `AmendAlbumDatesCase` to inject `RandomizeCoversPort`
  (same instance as issue 09) and `CoverRepository` (for the source-strip read).
- The AlbumView adapters (`AlbumViewCreateAlbumObserver`, `AlbumViewAmendDatesObserver`) forward the
  new event shapes to the updated handlers.

### Tests

- `pkg/catalog/album_create_test.go`:
  - `it_should_randomize_covers_on_the_new_album_when_medias_are_transferred_in`.
  - `it_should_strip_covers_from_source_albums_whose_medias_were_moved_out`.
  - `it_should_emit_empty_covers_when_no_transfer_is_involved`.
- `pkg/catalog/album_amend_dates_test.go`:
  - `it_should_randomize_destination_covers_from_moved_in_medias`.
  - `it_should_strip_source_covers_whose_media_moved_out`.
  - `it_should_leave_covers_untouched_when_no_transfer_is_involved`.
- `pkg/catalogviews/albums_view_events_test.go`:
  - `TestAlbumView_AlbumCreated/it_should_write_covers_on_the_owner_row` — dispatch with 3 covers;
    assert the row has them.
  - `TestAlbumView_AlbumDatesAmended/it_should_update_covers_for_each_affected_album` — seed owner
    + visitor on two albums; dispatch with new covers for both; assert all four rows updated;
    display fields untouched.
- Acceptance suite (`01-09`) extended: after `CreateAlbum` with transferred medias, `ListAlbums`
  returns the new album with the expected covers.
- **Dead-code cleanup — conditional**: `AlbumView.OnAlbumCoversChanged` is dead once issues 09, 10
  and 11 have all landed. If this is the last of the three to merge
  (`grep OnAlbumCoversChanged` shows no production caller), delete the method and its tests from
  `pkg/catalogviews/albums_view.go` and `pkg/catalogviews/albums_view_events_test.go` as part of
  this PR. The repository method `SetCoversForAllViewers` is kept.

## Out of scope

- `MediasInserted` — issue 09.
- `AlbumShared` — issue 11.
- Lifecycle maintenance for `AlbumDeleted`, `AlbumRenamed` folder-change — issue 12.

## References

- `../spec.md` (Automatic cover maintenance).
- `../design.md` (Event-payload matrix; Handling transfers-out).
- Load skills: `go`, `architecture`.
