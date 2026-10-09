# Feature Spec — Manage album covers (UI)

**Status:** draft
**Author:** Arch
**Date:** 2026-10-10

## Intent

Give the album owner a UI to curate the covers of an album: cherry-pick individual medias as covers,
unpick them, and re-randomise. This is the user-facing counterpart to the catalog cover model
already built in `specs/archived/2026-09-album-covers/` (Phase 2 — read model, propagation, API —
and the Phase 3 backend operation `Randomize` with its endpoint).

The feature closes the album-covers story: owners currently have no way to influence which images
represent their albums.

## Background — what already exists

- Catalog: canonical cover record (`{OWNER}#ALBUM / ALBUM#{FOLDER_NAME}#COVERS`) holding up to 4
  covers, each `{MediaId, Filename, Origin}` with `Origin ∈ {RANDOM, CHERRY_PICKED}`.
- Catalog operation `Randomize(albumId)`: replaces every `RANDOM` cover, preserves `CHERRY_PICKED`,
  refills empties up to 4. Called inline by lifecycle events and by the Phase 3 endpoint.
- REST: `GET /api/v1/albums` returns `covers[]` per album; `POST /api/v1/owners/{owner}/albums/{folderName}/covers/refresh`
  triggers `Randomize` and returns the new set.
- Web (`web-nextjs`): the album list already renders real covers (issue 07 of the previous feature).
- **No UI yet** to trigger a re-randomise or to pick/unpick a cover.
- **No `SetCovers` catalog operation yet, and no REST endpoint** to persist cherry-picks — Phase 4+
  of the previous spec was deferred and is in scope here.

## Concepts

Covers, cover origin, and the max-of-4 rule are defined in `docs/catalog/CONTEXT.md`. This feature
adds no new domain concept.

## User journeys

### Pick a media as a cover (owner)

From the album page (grid of all medias of an album), the owner marks an image as a cover. The
image is stored as `CHERRY_PICKED` in the album's cover set. If the album already has 4 covers, the
oldest `RANDOM` cover is evicted to make room; if there are no `RANDOM` covers left, the pick is
rejected with a clear message ("Unpick a cover first").

### Unpick a cover (owner)

From the same page, the owner removes a media from the cover set. The slot is left empty — it is
not auto-refilled by `RANDOM` on unpick (an explicit re-randomise does that).

### Re-randomise covers (owner)

The owner triggers a re-randomisation for the album. All `RANDOM` covers are redrawn from the
album's images; `CHERRY_PICKED` covers are preserved; empty slots are filled up to 4. The REST
endpoint and catalog operation already exist; this feature adds the UI control.

### See which medias are covers (owner & visitor)

On the album page, each cover media carries a visible indicator. The indicator distinguishes
`CHERRY_PICKED` from `RANDOM` for the owner (so they understand what re-randomise will change);
visitors see no distinction (or no indicator at all — out of scope to decide here).

## Scope

- **Catalog**: a `SetCovers` operation (name TBD during design) that persists a cherry-pick or
  unpick for one media, enforcing the max of 4 and the eviction rule above. Emits whatever
  lifecycle update is needed for the view to stay in sync.
- **REST**: endpoint(s) to pick and unpick a single media as a cover. Exact shape in `design.md`.
- **ACL**: a dedicated `CanManageAlbumCovers` role (name TBD) — see "Access control" below.
- **Web (`web-nextjs`)**: album-page controls for pick, unpick, re-randomise; owner-only visibility;
  optimistic update of the album list cache so album cards reflect the new covers without a reload.
- **Tests & stories**: Storybook stories covering owner / visitor views, loading & error states for
  each action.

## Out of scope

- Videos as covers.
- Ordering/arrangement of covers beyond the natural order.
- Per-user or per-viewer cover preferences.
- Backfill or any admin tooling (already shipped).
- Any change to the archive / image-delivery pipeline.
- Any change to lifecycle propagation (`AlbumCreated`, `MediasInserted`, …) — already shipped.

## Access control (fix)

The current `POST /covers/refresh` endpoint is wired in the Lambda Authorizer to `CanAmendAlbumDates`
(see `api/lambdas/authorizer/main.go:154`). This is **incorrect**: cover management is a distinct
capability and must not piggy-back on an unrelated permission.

- Introduce a **dedicated role** — provisional name `CanManageAlbumCovers` — in `pkg/acl/catalogacl`.
- Re-wire the existing refresh endpoint and the new pick/unpick endpoint(s) to use the new role.
- Default grant: owners always have it on their own albums. Granting it to non-owners is **not in
  scope** of this feature (no sharing-UI change); the role simply exists and is checked.
- Migration: no stored permission uses the old mapping — the fix is a code-only change.

## Decisions

- **Owner-only, Phase 1.** No sharing of cover-management rights to visitors. The role exists so
  the capability can be delegated later without a second refactor, but no UI exposes it.
- **Single-media pick/unpick, not bulk `SetCovers`.** The deferred Phase 4+ sketch in the previous
  design proposed `PUT …/covers` with the full ordered list. The pick/unpick-per-media model is
  simpler for the UI (tap a tile to toggle) and matches the mental model better. Bulk replace stays
  deferred.
- **Max-of-4 enforcement with implicit eviction of `RANDOM`.** When the owner picks a 5th cover and
  the set has at least one `RANDOM`, the oldest `RANDOM` is evicted. This keeps the pick action
  one-click; the alternative (reject + ask the owner to unpick first) is worse UX in the common
  case. If no `RANDOM` is evictable, the pick is rejected (the owner is in full-curation mode and
  must unpick explicitly). _Design to confirm the eviction order._
- **Unpick does not auto-refill.** Leaving the slot empty is predictable; the owner uses
  re-randomise if they want it filled.
- **Image URL identity unchanged.** The UI still builds image URLs from `owner` + `mediaId` +
  `filename` via the existing loader; covers continue to store identity only.

## Open questions

- Exact REST shape for pick/unpick: `PUT /…/covers/{mediaId}` + `DELETE /…/covers/{mediaId}`, or a
  single toggle endpoint? _To settle in `design.md`._
- Visual language for the cover indicator on the album-page grid (badge? border? corner icon?).
  _Design in the UI story, grounded in `ui-components`._
- Does the album-list cache update happen client-side from the pick/unpick response, or via a
  re-fetch of `GET /albums`? _To settle in `design.md`._

## References

- Previous feature (archived): `specs/archived/2026-09-album-covers/spec.md`,
  `specs/archived/2026-09-album-covers/design.md`.
- Catalog language: `docs/catalog/CONTEXT.md`.
- ACL authorizer wiring: `api/lambdas/authorizer/main.go`.
- Existing refresh handler: `api/lambdas/randomize-album-covers/`.
