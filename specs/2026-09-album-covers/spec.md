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

### Cherry-pick a cover
From a media, the owner marks it as a cover for its album.
- If a slot is free (fewer than 4 covers), the media is added as a `CHERRY_PICKED` cover.
- If all 4 slots are full but at least one is `RANDOM`, a `RANDOM` cover is evicted to make room (no prompt).
- If all 4 covers are `CHERRY_PICKED`, the action is refused with an error asking the owner to remove one
  cover first.

### Remove a cover
The owner unpicks a cover; the slot is left **empty** (no automatic refill). Empty slots are filled later by
re-randomising or by backup adding new medias.

### Re-randomise covers
The owner triggers a re-randomisation for an album. All `RANDOM` covers are replaced by a fresh random
selection of eligible images; `CHERRY_PICKED` covers are preserved. Empty (previously unpicked) slots are
filled up to 4.

### Backup fills missing covers
When backup indexes new medias into an album, any empty cover slots are filled with `RANDOM` covers drawn
from the album's eligible images. Existing covers (random or cherry-picked) are left untouched.

### Backfill existing albums
A one-off operation randomly picks covers for every album of every owner that has empty slots, so existing
albums get covers without waiting for new medias.

## Scope

- Catalog: model covers, persist them, and the logic to pick/replace/re-randomise/backfill.
- Catalogviews: covers denormalised into the per-user album-list view for read-optimised retrieval.
- API: `list-albums` returns covers per album; endpoints to cherry-pick, remove, and re-randomise.
- Backup: fill empty cover slots when indexing new medias.
- Web (`web-nextjs`): render real covers on album cards; UI to cherry-pick, remove, re-randomise.
- A backfill entry point for all albums of all owners.

## Out of scope

- Videos as covers.
- More than 4 covers, or configurable cover count.
- Ordering/arrangement UI beyond the natural cover order.
- Per-user or per-viewer cover preferences.
- Changes to the image resize/format pipeline (owned by image-delivery-optimisation).

## Open questions

_None open._

## Decisions

- **Complete read model** — the album-list view record is promoted to a complete album projection
  (`{albumId, name, start, end, count, covers[]}`). This is done as a **refactoring story that lands before
  covers** (covers then ride on the reshaped view rather than a record we'd immediately rework). Warrants an
  ADR. See the "list-albums review" comment for rationale.
- **Sharing grid stays out** — the owner's "shared with" badges remain a separate query; denormalising them
  crosses into ACL and is out of scope.
- **Covers live in a single record** — the whole cover set of an album is one DynamoDB item
  (`PK = {OWNER}#ALBUM`, `SK = ALBUM#{FOLDER_NAME}#COVERS`), not one item per cover. Covers are never
  accessed individually; every write reads the set, enforces the max of 4, and rewrites it. This supersedes
  the per-cover key floated in grilling round 1 (Q1).

## Phasing

The feature ships in ordered phases. Each issue names its phase and treats later phases as out of scope.

1. **Phase 1 — Foundation**: the album-list view becomes a complete read model (refactor, no user-facing
   change).
2. **Phase 2 — Covers populated & displayed (read-only)**: cover model & persistence, random completion of
   empty slots, denormalisation into the view, `GET /albums` returns covers, backup completes touched
   albums, admin backfill, and the album list renders covers. No user editing yet.
3. **Phase 3 — Re-randomise**: replace `RANDOM` covers on demand — backend, API, UI.
4. **Phase 4 — Manually pick**: cherry-pick a media as a cover — backend, API, UI.
5. **Phase 5 — Unpick**: remove a cover, leaving the slot empty — backend, API, UI.

## Comments

### Grilling — round 1 (2026-09-06)

- **Q1 Source of truth**: canonical covers stored as their own items under the album partition —
  `PK = {OWNER}#ALBUM`, `SK = ALBUM#{FOLDER_NAME}#COVER#{mediaId}`. Denormalised into the per-user view
  records for read.
- **Q2 Cover contents & order**: order does not matter (preference: capture date ascending, but any order is
  fine). Cover stores `mediaId` + `filename` + `origin`; no dimensions/URL.
- **Q3 When covers are selected**: three operations — (a) an initialisation/backfill **batch function** (to
  be written), (b) on **backup**: any album touched by an addition gets its covers re-checked and completed,
  (c) on demand from the user: pick, unpick, re-randomise.
- **Q4 Backup fill source**: cheap option accepted. At the **end of a backup batch**, only albums that
  received a new addition are checked; if their cover set is not full it is completed. Full sets are left
  untouched (cheap even if backup works image by image).
- **Q5 Authorization**: only users who can edit the album (owner-level). Visitors read-only. _(recommended,
  assumed accepted — confirm if wrong.)_
- **Q6 Backfill invocation**: a CLI entry point — either a new `dphoto` command or reviving
  `tools/dphotoops`. Exact form decided during implementation of the story.
- **Q7 UI placement**: the UI is being rewritten in parallel; do not read the current UI. Assume the
  **re-randomise** action on the album page (grid of all pictures) and **pick/unpick** on the media page
  (fullscreen media).

### list-albums read path review (2026-09-06)

Current `GET /albums` makes ~5 DynamoDB round-trips, and it is **not** a true 1+P: shared album metadata is
fetched with a single `BatchGetItem`, and media counts come in-memory from the view. The redundancy: the
per-user view record stores **only the count**, so album display fields (Name/Start/End) are re-fetched
through two divergent paths (owned = Query on `{OWNER}#ALBUM`, shared = BatchGet by ids) and merged. The
owned/shared provider split exists only because those fields are absent from the view.

Proposed direction (candidate ADR): promote the album-list view record to a **complete read model**
`{albumId, name, start, end, count, covers[]}`. Reads then collapse to 1 Query (`GetAvailabilitiesByUser`,
owned + shared) + 1 Query for the owner's sharing grid. Album create/rename/amend-dates must then propagate
display fields to viewer view records (reusing the existing count-propagation + drift-rebuild path).
Denormalising the sharing grid is out of scope (crosses into ACL).
</content>
</invoke>
