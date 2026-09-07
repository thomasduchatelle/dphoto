# Stories — Album covers

Draft of every issue, grouped by phase. Phases ship in order; each issue names its phase and treats later
phases as out of scope. Phase 1–2 issues are written out under `issues/`; phases 3–5 are drafted here.

See `spec.md` (what) and `design.md` (cross-issue technical direction).

## Phase 1 — Foundation

- **01 — Album-list view becomes a complete read model** _(catalogviews + adapter; ADR)_
  - Extend the view record with Name/Start/End; write them on album create/rename/amend-dates.
  - Collapse the read to a single per-user Query; drop the owned/shared metadata re-fetch.
  - Drift-rebuild covers the new fields. `GET /albums` output unchanged.

## Phase 2 — Covers populated & displayed (read-only)

- **02 — Cover model & single-record persistence + random completion** _(pkg/catalog + catalogdynamo)_
  - `Cover`/`CoverOrigin`, single `…#COVERS` record per album (max 4), completion operation with uniform
    random IMAGE selection (candidate and query variants).
- **03 — Covers in the album-list read model** _(catalogviews)_
  - `Covers` in the projection; re-denormalise on change; drift-rebuild from canonical.
- **04 — Expose covers on `GET /albums`** _(api)_
  - `covers[]` on `AlbumDTO`.
- **05 — Backup completes covers for touched albums** _(pkg/backup + cmd)_
  - End-of-batch completion from the just-inserted images; full sets untouched.
- **06 — Backfill covers for all owners/albums** _(CLI / tools)_
  - One-off idempotent completion sweep across every album.
- **07 — Render album covers on the album list** _(web-nextjs)_
  - Album card shows real covers (0–4) via the image loader.

## Phase 3 — Re-randomise

- **08 — Re-randomise operation** _(pkg/catalog)_
  - Replace all `RANDOM` covers with a fresh random pick, keep `CHERRY_PICKED`, then complete empties.
- **09 — Re-randomise endpoint** _(api)_
  - `POST …/covers/refresh`, owner-edit permission; re-denormalises to the view.
- **10 — Re-randomise action on the album page** _(web-nextjs)_
  - Control on the album page (grid of all pictures) triggering the endpoint and refreshing covers.

## Phase 4 — Manually pick

- **11 — Pick operation** _(pkg/catalog)_
  - Add a `CHERRY_PICKED` cover; evict a `RANDOM` when full; refuse when all 4 are `CHERRY_PICKED`.
- **12 — Pick endpoint** _(api)_
  - `PUT …/covers/{mediaId}`, owner-edit permission; surfaces the refusal error.
- **13 — Pick on the fullscreen media page** _(web-nextjs)_
  - "Set as cover" on the fullscreen media; snackbar on refusal.

## Phase 5 — Unpick

- **14 — Unpick operation** _(pkg/catalog)_
  - Remove a cover, leave the slot empty (no auto-refill).
- **15 — Unpick endpoint** _(api)_
  - `DELETE …/covers/{mediaId}`, owner-edit permission.
- **16 — Unpick on the fullscreen media page** _(web-nextjs)_
  - "Unset as cover" on the fullscreen media.
