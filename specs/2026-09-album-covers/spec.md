# Feature Spec — Album covers

**Status:** draft
**Author:** Arch
**Date:** 2026-09-06

## Intent

Let each album be represented on the album list by up to four featured photos — its **covers** — served
directly by the backend alongside the album list. Today the website fetches the last four medias of every
album with a separate request per album: the idea landed well with users but is slow. This feature makes it
a first-class, backend-supported concept: covers are chosen (randomly or by the owner), stored, and returned
in the same album-list response so the page paints without N per-album media queries.

See the catalog language for the terms **Cover** and **CoverOrigin** (`docs/catalog/CONTEXT.md`).

## Background — how it works today

- The old `web/` computes covers client-side by requesting the last four medias of each album. `web-nextjs`
  already exposes a placeholder `thumbnails?: string[]` on its `Album` type ("test purposes").
- `list-albums` is served by **catalogviews**: per-user view records (`USER#{EMAIL}#ALBUMS_VIEW / …#COUNT`)
  holding the media count, not the album metadata partition (`{OWNER}#ALBUM / ALBUM#{FOLDER_NAME}`).
- The image-delivery-optimisation feature (in flight) is changing the image URL scheme and formats
  (WebP, a 4-tier width ladder, LQIP). The frontend builds image URLs from `owner/mediaId/filename?w=…`
  through its loader. Covers therefore store identity (`mediaId` + `filename`), never a URL.

## Concepts

- An album has an ordered set of **0 to 4 covers**.
- A cover references a media by `mediaId` + `filename` (enough for the frontend loader to build the URL).
  No URL, format, or width is persisted — those are the archive's concern and are changing.
- Only medias of type `IMAGE` are eligible.
- Each cover has a **CoverOrigin**: `RANDOM` (auto-picked) or `CHERRY_PICKED` (chosen by the owner).
- Covers are **owner-defined and shared with every viewer** of the album (not per-user).

## User journeys

### See an album's covers
Opening the album list shows each album card with its up-to-four covers, returned in the same response as
the album list (no per-album request).

### Re-randomise covers (Phase 3)
The owner triggers a re-randomisation for an album. All `RANDOM` covers are replaced by a fresh random
selection of eligible images; `CHERRY_PICKED` covers are preserved. Empty (previously unpicked) slots are
filled up to 4.

### Automatic cover maintenance (invisible to the user)
Covers stay in sync with the album's medias without any user action. Whenever an album's media set changes,
the catalog decides whether covers must be refreshed:

- **Medias added to an album** (via `MediasInserted`, `AlbumCreated` with transferred medias, or
  `AlbumDatesAmended` with transferred medias): if the album has empty cover slots, the catalog fills them
  with `RANDOM` covers drawn from the newly-inserted images. `CHERRY_PICKED` and existing `RANDOM` covers
  are preserved; already-full sets are untouched.
- **Medias moved out of an album** (via `AlbumDatesAmended` or `AlbumRenamed` with transferred medias): any
  cover whose media is no longer in the album is stripped; empty slots created this way are left empty
  (they'll be refilled next time medias are added or by a re-randomise).
- **Album deleted**: the canonical cover record is deleted.
- **Album renamed (folder changes)**: the canonical cover record moves with the album; covers survive.
- **Album shared**: the new visitor's view row is written with the current cover set, so the visitor's
  next album-list request shows the album with its covers.

### Backfill existing albums
A one-off operation randomly picks covers for every album of every owner that has empty slots, so existing
albums get covers without waiting for new medias.

### Cherry-pick / unpick (Phase 4+ — out of scope here)
The owner picks or unpicks an individual media as a cover. Deferred.

## Scope

**Phase 2 — Covers populated, propagated, displayed (read-only)**:

- Catalog: cover model, persistence, and the `Randomize` operation (keeps `CHERRY_PICKED`, replaces
  `RANDOM`, fills empties up to 4).
- Catalog use cases emit covers on the events whose outcome can change them (`AlbumCreated`,
  `MediasInserted`, `AlbumDatesAmended`). Lifecycle maintenance for `AlbumDeleted`, `AlbumRenamed`, and
  `AlbumShared` is handled by the catalog without touching the backup domain.
- Catalogviews: covers denormalised into the per-user album-list view, driven by the lifecycle events
  above.
- API: `list-albums` returns covers per album.
- Web (`web-nextjs`): render real covers on album cards.
- Admin backfill via `dphotops covers backfill` for existing albums.

**Phase 3 — Owner-triggered re-randomise**:

- Catalog operation that re-randomises the covers of one album.
- REST endpoint, owner-edit permission.
- UI control on the album page.

## Out of scope

- Videos as covers.
- More than 4 covers, or configurable cover count.
- Ordering/arrangement UI beyond the natural cover order.
- Per-user or per-viewer cover preferences.
- Changes to the image resize/format pipeline (owned by image-delivery-optimisation).
- Cherry-pick / unpick / `SetCovers` (Phase 4+).
- **Any change to the backup domain** — backup must stay untouched; covers are a catalog concern,
  driven by catalog events.

## Decisions

- **Complete read model** _(Phase 1, done)_ — the album-list view record was promoted to a complete album
  projection (`{albumId, name, start, end, count, covers[]}`). See ADR-0004.
- **Covers live in a single canonical record** — the whole cover set of an album is one DynamoDB item
  (`PK = {OWNER}#ALBUM`, `SK = ALBUM#{FOLDER_NAME}#COVERS`), not one item per cover. Every write reads
  the set, enforces the max of 4, and rewrites it.
- **Covers propagate on existing lifecycle events — no dedicated `CoversChanged` event.** Each catalog
  event carries `Covers` when its outcome can change them (`AlbumCreated`, `MediasInserted`,
  `AlbumDatesAmended`, `AlbumShared`). Events that cannot change them do not carry them
  (`AlbumRenamed` keeps the same covers in a migrated record; `AlbumDeleted` wipes them).
- **`Randomize` is the single completion primitive.** It replaces `RANDOM` covers, preserves
  `CHERRY_PICKED`, and fills empties up to 4. It is the operation called whenever covers may need to be
  refreshed — media insertion, dates amendment with transfer-in, administrative backfill, and (Phase 3)
  the owner-triggered re-randomise.
- **Drift reconciliation is exceptional.** `dphotops drift --apply` is a bug-recovery tool, not a
  propagation path. Normal operation keeps the view consistent through events; drift exists to repair
  divergence after an incident.
- **Backup domain stays untouched.** When medias land in an album, the catalog decides whether covers
  need to change based on the `MediasInserted` event (which it already owns). The backup project does
  not depend on the covers logic.
- **Sharing grid stays out** — the owner's "shared with" badges remain a separate query; denormalising
  them crosses into ACL and is out of scope.

## Phasing

The feature ships in ordered phases. Each issue names its phase and treats later phases as out of scope.

1. **Phase 1 — Foundation** _(done)_: the album-list view is a complete read model.
2. **Phase 2 — Covers populated, propagated, displayed (read-only)**: cover model & persistence, the
   `Randomize` operation, catalog events carry covers so the view stays in sync, lifecycle maintenance
   (`AlbumDeleted`, `AlbumRenamed`, `AlbumShared`), `GET /albums` returns covers, admin backfill, album
   list renders covers. **No user editing yet.**
3. **Phase 3 — Owner-triggered re-randomise**: operation, API endpoint, UI control.
4. **Phase 4+ — Cherry-pick / unpick / `SetCovers`** _(deferred, not specified)_.

## Open questions

_None open._

## Comments

### Grilling — round 1 (2026-09-06)

- **Q2 Cover contents & order**: order does not matter (preference: capture date ascending). Cover
  stores `mediaId` + `filename` + `origin`; no dimensions/URL.
- **Q5 Authorization**: only users who can edit the album (owner-level). Visitors read-only.
- **Q6 Backfill invocation**: CLI entry point `dphotops covers backfill`.
- **Q7 UI placement**: Phase 3 action **re-randomise** on the album page.

### list-albums read path review (2026-09-06, Phase 1 done)

Captured in ADR-0004. Kept here for traceability.
