# Stories — Album covers

Forward-looking backlog, grouped by phase. Each issue is a vertical slice that can be owned end-to-end
by one agent. See `spec.md` (what) and `design.md` (cross-issue technical direction).

## Phase 2 — Covers propagated, lifecycle-maintained, displayed

- **07 — Render album covers on the album list** _(web-nextjs)_
  - Album card shows real covers (0–4) via the image loader. Replaces the placeholder `thumbnails`
    field. Storybook / visual tests for 0, partial, full.
- **09 — Covers on `MediasInserted` + Randomize primitive upgrade** _(pkg/catalog + pkg/catalogadapters + pkg/catalogviews + pkg/pkgfactory + cmd/dphotops)_
  - Upgrade the cover-completion primitive: rename `CompleteCovers` → `RandomizeCovers`, drop every
    `RANDOM`, keep `CHERRY_PICKED`, fill empties up to 4; return the resulting `[]Cover`.
  - Media-insert use case calls the primitive per affected album; the `MediasInserted` event carries
    the resulting covers per album.
  - `AlbumView.OnMediasInserted` denormalises covers alongside the count update in a single per-row
    write (new atomic repository method `ADD Count :d SET Covers = :c`).
  - `pkg/backup` and `cmd/dphoto` are **not** touched.
- **10 — Covers on `AlbumCreated` and `AlbumDatesAmended`** _(pkg/catalog + pkg/catalogviews)_
  - Both events gain `Covers map[AlbumId][]Cover` covering the amended album and every
    source/destination touched by `TransferredMedias`.
  - Destinations: `Randomize` with the moved-in medias. Sources: strip stale covers whose media
    moved out (no re-randomisation).
  - `AlbumView` handlers denormalise via `PutSummaries` (create) or `SetCoversForAllViewers`
    (amend-dates) — no atomic ADD+SET needed.
- **11 — Covers on `AlbumShared`** _(pkg/acl/catalogacl + pkg/catalogviews)_
  - `ShareAlbumCase` reads the current covers and passes them to `AlbumSharedObserver`; the view
    adapter writes a complete visitor row in one `PutSummaries` call.
- **12 — Canonical cover record lifecycle (delete + folder-rename migration)** _(pkg/catalog + pkg/catalogadapters + pkg/catalogviews)_
  - `AlbumDeleted` → delete the canonical `#COVERS` record.
  - `AlbumRenamed` folder-change → move the canonical record to the new SK; AlbumView's folder-rename
    path preserves covers on the new rows. In-place rename is a no-op.

### Dead-code cleanup

`AlbumView.OnAlbumCoversChanged` is dead once 09, 10, 11 have all landed. The last of the three to
merge grep-checks the method has no production caller and deletes it in the same PR. No separate
ticket.

## Phase 3 — Owner-triggered re-randomise

- **13 — Randomize use case (owner-triggered)** _(pkg/catalog + pkg/catalogviews + pkg/pkgfactory)_
  - A `RandomizeAlbumCovers` use case that calls `Randomize(albumId)` and emits a dedicated
    `AlbumCoversRandomised` event; `AlbumView` consumes it to update every viewer row.
- **14 — `POST …/covers/refresh` endpoint** _(api/lambdas + deployments/cdk)_
  - New lambda, owner-edit permission, returns the new covers in the response body.
- **15 — Re-randomise UI on the album page** _(web-nextjs)_
  - Owner-only action on the album page; calls the endpoint and refreshes the displayed covers.

## Phase 4+ — Cherry-pick / unpick / `SetCovers` (deferred)

Out of scope for this feature. A future `PUT /api/v1/owners/{owner}/albums/{folderName}/covers`
endpoint will replace the per-media approach from the earlier draft.

## Dependency graph

```
Phase 2 — four parallel tracks:

   ┌──────────────────────────────────────────────────────────┐
   │  09 ──► 10                                                │  catalog event propagation
   │     └─► 11                                                │  (10 and 11 parallel after 09)
   └──────────────────────────────────────────────────────────┘

   ┌──────────────────────────────────────────────────────────┐
   │  12                                                        │  canonical record lifecycle
   └──────────────────────────────────────────────────────────┘  (independent)

   ┌──────────────────────────────────────────────────────────┐
   │  07                                                        │  web rendering
   └──────────────────────────────────────────────────────────┘  (independent, already in flight)

Phase 3 — one sequence (after Phase 2 is green):

   13 ──► 14 ──► 15
```

### What can start when

- **Immediately, in parallel** — four agents, no in-flight dependency between them:
  - **07** (web rendering — already in flight; depends on 04 which is `done`)
  - **09** (MediasInserted + Randomize primitive upgrade)
  - **11** (AlbumShared — only reads covers, doesn't need the upgraded primitive)
  - **12** (canonical record lifecycle)
- **After 09 lands** (upgraded `RandomizeCoversPort`):
  - **10** starts (needs the primitive that returns `[]Cover` to attach to events).
- **Phase 3**:
  - **13 → 14 → 15** is a strict sequence (each needs the previous layer).

Numbering gaps: **08** was folded into 09 (so the primitive upgrade ships with its first consumer
rather than as a schema-only PR). Issue numbers 05 (`wontdo`) and 08 (unused) are left as gaps to
avoid renumbering in-flight branches.
