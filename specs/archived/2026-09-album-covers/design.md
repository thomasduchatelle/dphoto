# Design — Album covers

Cross-issue technical direction. Issues link here rather than repeat it. The "what" is in `spec.md`.

## Read model (Phase 1 — done)

Captured in ADR-0004. The per-user album-list view record (`AlbumSummary`) carries identity, display
fields (`Name`, `Start`, `End`) and `Count`. Writes are split so count updates and display-field
updates never clobber each other. Reads are served by a single Query per user.

Covers are stored as a **sibling row in the same partition** rather than as a field on the summary —
see §View projection below.

## Cover canonical record

The whole cover set of an album is a **single** DynamoDB item:

- PK `{OWNER}#ALBUM`, SK `ALBUM#{FOLDER_NAME}#COVERS`.
- Holds an ordered list of at most 4 covers, each `{MediaId, Filename, Origin}` where `Origin` is
  `RANDOM` or `CHERRY_PICKED`.
- Every write reads the set, enforces the max of 4, and rewrites the whole item.
- Deleted when the album is deleted; migrated (same owner, new folder-name SK) when the album is renamed
  with a folder change.

Documented in `DATA_MODEL.md`.

## Cover-maintenance service

Covers are a **derived read model**. A single service in `pkg/catalog` offers two strategies:

- **`Refresh(albumId, removed)`** — drop every `RANDOM` cover, keep every `CHERRY_PICKED` cover, strip
  any kept cover whose `MediaId` is in `removed`, then fill up to 4 by drawing uniformly at random
  from the album's full image set. Used when the cover set should reflect what is "newest or most
  interesting" — media insertion, admin backfill, Phase 3 owner-triggered re-randomise.
- **`Stabilise(albumId, removed)`** — keep every existing cover (both `RANDOM` and `CHERRY_PICKED`)
  except those in `removed`, then fill empty slots by drawing uniformly at random from the album's
  full image set. Used when stability is preferred over freshness — album created with transferred
  medias, dates amended, album deleted (destination albums).

Both methods return `([]Cover, changed bool, error)`. `changed = true` only when the resulting set
differs from the one loaded from the canonical record. Both load the album's full `IMAGE` medias
unconditionally — the "no query needed" case (empty cover set + fresh medias) is an edge not worth
a dedicated code path.

### Called inline by the catalog use cases

Catalog use cases depend on the cover-maintenance service and call it inline, after their domain
mutation and before firing their lifecycle event. The resulting cover sets travel on the event
payload as `Covers map[AlbumId][]Cover`. Albums whose cover set did not change are **omitted** from
the map (a missing album means "no cover update required").

This mirrors the existing pattern for `TransferMediasService`: an internal service of the catalog
domain, injected into each use case that may need it, called inline, with its result folded into the
emitted event. No internal event bus, no in-domain observers.

### Per-operation strategy

| Catalog operation       | Per affected album                                                                 |
|-------------------------|------------------------------------------------------------------------------------|
| `InsertMedias`          | `Refresh(albumId, nil)` for every album that received medias                       |
| `AlbumCreated`          | `Stabilise(newId, nil)`; `Stabilise(sourceId, movedOutIds)` per source album       |
| `AlbumDatesAmended`     | `Stabilise(destId, nil)` per destination; `Stabilise(sourceId, movedOutIds)` per source |
| `AlbumDeleted`          | Deleted album: canonical record removed. `Stabilise(destId, nil)` per destination  |
| `AlbumRenamed` (folder) | Canonical record moved to the new SK. No cover reconciliation                      |
| `AlbumRenamed` (in-place) | No-op on covers                                                                  |
| `AlbumShared`           | Current covers of the album read inline and attached to the event                  |

## View projection

The `ListAlbums` read path serves everything from a single Query on the user's partition
(`USER#{EMAIL}#ALBUMS_VIEW`). To keep that invariant while making covers a sibling concern, each
album has **two SK rows per viewer**:

- `ALBUM#{owner}#{folder}` — identity, display fields (`Name`, `Start`, `End`), `MediaCount`.
- `ALBUM#{owner}#{folder}#COVERS` — the ordered list of up to 4 covers.

The DynamoDB adapter's `ListSummariesForUser` returns rows merged per album; a missing cover row
means "no covers". Albums with no summary row are skipped (orphan cover rows are discarded
defensively).

### Why two rows, not one field

- Cover writes never touch the main summary row → no attribute-footprint discipline, no clobber risk,
  no need for atomic `ADD Count :d SET Covers = :c` updates.
- Cover writes stay independent from count / display-field writes; the view's event handlers for
  `MediasInserted`, `AlbumCreated`, `AlbumDatesAmended` write cover rows and summary rows
  independently (both idempotent per attribute).
- Share: writing the visitor's cover row is one `PutItem`, decoupled from the main summary row's
  write.
- Still one Query at read time — the user partition remains the single read source.

### How the view handlers consume the event

Each `AlbumView` lifecycle handler receives the catalog event, which carries both the data it
already handled (counts, display fields) and the new `Covers map[AlbumId][]Cover`. For every entry
in the map:

- Non-empty cover list → upsert the viewer cover rows for that album.
- Explicitly empty cover list → delete the viewer cover rows for that album.
- Absent album from the map → no cover write.

Fan-out across viewers uses the existing `AlbumViewByAlbumIndex` GSI already used for count and
display-field updates.

### Drift

Drift reconciliation already rebuilds view rows from the canonical catalog. It rebuilds the
`#COVERS` sibling rows in the same pass, from the canonical cover record. No new behaviour
otherwise.

## REST contract

### Phase 2

`GET /api/v1/albums` — each album in the response has `covers[]` (0–4 entries):

```json
"covers": [
  { "mediaId": "...", "filename": "...", "origin": "RANDOM" }
]
```

`origin` ∈ `RANDOM | CHERRY_PICKED`. The frontend builds image URLs from `owner` + `mediaId` +
`filename` via its existing image loader — no URL is stored or returned.

### Phase 3 — owner-triggered re-randomise

| Phase | Action        | Route                                                              |
|-------|---------------|--------------------------------------------------------------------|
| 3     | Re-randomise  | `POST /api/v1/owners/{owner}/albums/{folderName}/covers/refresh`   |

- Authorisation: owner-edit permission on the album (same rule as rename / amend-dates).
- Body: empty.
- Response: `200 OK` with the new covers list.
- Semantics: `Refresh(albumId, nil)` — redraws every `RANDOM` cover from the full album, keeps
  `CHERRY_PICKED`.

### Phase 4+ (deferred)

A single `PUT /api/v1/owners/{owner}/albums/{folderName}/covers` endpoint will replace the per-media
pick/unpick from the earlier draft, taking the full ordered list of up to 4 covers. Not specified here.
