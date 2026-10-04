# Stories — Album covers

Forward-looking backlog, grouped by phase. Each issue is a vertical slice that can be owned end-to-end
by one agent. See `spec.md` (what) and `design.md` (cross-issue technical direction).

## Phase 2 — Covers propagated, lifecycle-maintained, displayed

Covers are kept in sync with each album's medias through a cover-maintenance primitive and one
observer per catalog lifecycle event. The album-list view projects covers as a sibling row in the
user's partition so `ListAlbums` stays a single Query.

- **07 — Render album covers on the album list** _(web-nextjs)_
  - Album card shows real covers (0–4) via the image loader. Replaces the placeholder `thumbnails`
    field. Storybook / visual tests for 0, partial, full.
- **09 — Covers on MediasInserted + split view projection + Refresh/Stabilise service** _(pkg/catalog + pkg/catalogviews + pkg/catalogviewsadapters/catalogviewsdynamodb + pkg/pkgfactory + cmd/dphotops)_
  - Introduces the `CoverMaintenance` service with two strategies: **Refresh** (drop RANDOM, keep
    CHERRY_PICKED, strip removed, refill from the album's full image set) and **Stabilise** (keep
    every existing cover except removed, refill from the album's full image set).
  - Splits the view: covers move out of `AlbumSummary` into a sibling `#COVERS` SK row in the same
    user partition; one Query still serves `ListAlbums`.
  - `InsertMedias` calls Refresh inline per affected album and attaches the resulting covers to the
    `MediasInserted` event payload. `AlbumView.OnMediasInserted` writes count and covers in one call.
  - Admin backfill (`BackfillCovers` + `dphotops covers backfill`) runs Refresh and fans the result
    to the view.
- **10 — AlbumCreated cover maintenance** _(pkg/catalog + pkg/pkgfactory)_
  - Reconciliation observer for `AlbumCreated`: fills the new album's covers from transferred-in
    medias; strips covers of pre-existing source albums whose medias were moved out, refilling from
    survivors.
- **11 — AlbumDatesAmended cover maintenance** _(pkg/catalog + pkg/pkgfactory)_
  - Reconciliation observer for `AlbumDatesAmended`: destinations get empties filled from moved-in
    medias; sources get stale covers stripped and refilled from survivors.
- **12 — AlbumDeleted cover maintenance** _(pkg/catalog + pkg/catalogadapters/catalogdynamo + pkg/pkgfactory)_
  - Canonical `#COVERS` record is deleted; destinations that absorbed medias get their covers
    reconciled.
- **16 — AlbumRenamed cover retention** _(pkg/catalog + pkg/catalogadapters/catalogdynamo + pkg/pkgfactory)_
  - Folder-change: the canonical `#COVERS` record is moved to the new SK, covers survive on every
    viewer's row. In-place rename: no cover change.
- **17 — AlbumShared / AlbumUnshared cover propagation** _(pkg/acl/catalogacl + pkg/catalogviews + pkg/pkgfactory)_
  - On share, the visitor's `#COVERS` SK row is written from the owner's current covers so the next
    `ListAlbums` returns the album with its covers in one Query. On unshare, the row is deleted.

## Phase 3 — Owner-triggered re-randomise

- **13 — Randomize use case (owner-triggered)** _(pkg/catalog + pkg/pkgfactory)_
  - A `RandomizeAlbumCovers` use case that calls `CoverMaintenance.Reconcile(albumId, nil, nil)` —
    the fallback-query path redraws every `RANDOM` cover.
- **14 — `POST …/covers/refresh` endpoint** _(api/lambdas + deployments/cdk)_
  - New lambda, owner-edit permission, returns the new covers in the response body.
- **15 — Re-randomise UI on the album page** _(web-nextjs)_
  - Owner-only action on the album page; calls the endpoint and refreshes the displayed covers.

## Phase 4+ — Cherry-pick / unpick / `SetCovers` (deferred)

Out of scope for this feature. A future `PUT /api/v1/owners/{owner}/albums/{folderName}/covers`
endpoint will replace the per-media approach from the earlier draft.

## Dependency graph

```
Phase 2:

   07  (web rendering — independent, already in flight)

   09  (primitive + split view projection + MediasInserted)
      ├─► 10  (AlbumCreated)
      ├─► 11  (AlbumDatesAmended)
      ├─► 12  (AlbumDeleted)
      ├─► 16  (AlbumRenamed)
      └─► 17  (AlbumShared / AlbumUnshared)

Phase 3 (after Phase 2 is green):

   13 ──► 14 ──► 15
```

### What can start when

- **Immediately, in parallel**:
  - **07** (web rendering — depends on 04 which is `done`).
  - **09** (primitive + split view projection + MediasInserted observer).
- **After 09 lands** — five parallel tracks, each a small cover observer plus its wiring; distinct
  files, no shared surface beyond mechanical additions to `pkgfactory/factory_catalog.go`:
  - **10** (AlbumCreated), **11** (AlbumDatesAmended), **12** (AlbumDeleted), **16** (AlbumRenamed),
    **17** (AlbumShared / AlbumUnshared).
- **Phase 3**: **13 → 14 → 15** is a strict sequence (each needs the previous layer).

Numbering gaps: **05** (`wontdo`) and **08** (unused) are left as gaps to avoid renumbering in-flight
branches. **16** and **17** were added for the per-operation split under the new pattern.
