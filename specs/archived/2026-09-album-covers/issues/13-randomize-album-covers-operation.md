# 13 — Owner-triggered randomize operation

Status: done
Phase: 3
Layer: catalog domain — `pkg/catalog` + `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 09

## Description

Expose a catalog use case that lets the owner re-randomise an album's covers on demand. Reuses the
`Randomize` operation shipped by issue 09 (drop `RANDOM`, keep `CHERRY_PICKED`, fill empties) and propagates
the new covers through a lifecycle event so the view updates via the standard path.

The user-facing effect: the random-origin covers of the album are redrawn; cherry-picked covers stay in
place. Called by the REST endpoint (issue 14) and by the UI control (issue 15).

See `../design.md` → Randomize operation; Event-payload matrix.

## Acceptance criteria

- `pkg/catalog`:
  - New use case `RandomizeAlbumCovers` (struct with its port dependencies — `RandomizeCoversPort`
    plus whatever observer registry pattern the catalog uses for the chosen event).
  - Method signature: `RandomizeAlbumCovers.Randomize(ctx, albumId AlbumId) ([]Cover, error)`.
    Returns the new covers so the API layer can echo them in the response.
  - Behaviour:
    1. Calls `RandomizeCoversPort.RandomizeCovers(ctx, albumId)` (the query-variant, which lists the
       album's `IMAGE` medias and randomises).
    2. Emits a lifecycle event carrying the new covers. **Event choice**: a new dedicated event
       `AlbumCoversRandomised{AlbumId, Covers []Cover}`. Rationale: no other state changed, so
       re-using `MediasInserted`/`AlbumDatesAmended` would misrepresent what happened. Keep the
       payload minimal.
  - No permission check inside the use case — the API layer (issue 14) authorises before calling.
    This keeps the catalog use case usable from any trusted context (future scripts, admin tools).
- `pkg/catalogviews`:
  - `AlbumView.OnAlbumCoversRandomised(ctx, event)` writes the new covers to every viewer row via
    the per-cover repository method (`SetCoversForAllViewers`, kept from issue 03).
  - A new adapter `AlbumViewCoversRandomisedObserver` implementing the new observer interface,
    wired via the factory.
- `pkg/pkgfactory`:
  - New factory `RandomizeAlbumCoversCase(ctx)` returning the use case with the observer registered.
- Tests:
  - `pkg/catalog/randomize_album_covers_test.go`:
    - `it_should_call_randomize_and_emit_the_event_with_the_new_covers`.
    - `it_should_return_the_new_covers_to_the_caller`.
    - `it_should_propagate_errors_from_the_randomize_port`.
  - `pkg/catalogviews/albums_view_events_test.go`:
    - `TestAlbumView_AlbumCoversRandomised/it_should_update_covers_on_all_viewer_rows`.
    - `TestAlbumView_AlbumCoversRandomised/it_should_leave_count_and_display_fields_untouched`.
- Acceptance suite extended: after calling `RandomizeAlbumCovers`, `ListAlbums` returns the new
  covers for the owner and every visitor.

## Out of scope

- REST endpoint — issue 14.
- UI — issue 15.
- Permission checks — enforced at the API layer (issue 14).

## References

- `../spec.md` (Phase 3; User journeys → Re-randomise covers).
- `../design.md` (Randomize operation; §Phase 3 REST contract).
- Load skills: `go`, `architecture`.
