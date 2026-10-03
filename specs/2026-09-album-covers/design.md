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

## Cover-maintenance pattern

Covers are a **derived read model**: they are kept in sync reactively, after the fact. The catalog use
cases (create, delete, rename, amend-dates, insert-medias) do not know covers exist and their event
payloads do not carry covers.

A single primitive, `CoverMaintenance.Reconcile(albumId, added, removed)`, applies the invariant:

1. Load the current canonical cover set.
2. Drop every `RANDOM` cover; keep every `CHERRY_PICKED` cover.
3. Strip any kept cover whose `MediaId` is in `removed`.
4. Fill empty slots up to 4, drawing uniformly at random from `added` first (eligible = `IMAGE`, not
   already a cover); if `added` is insufficient, fall back to querying the album's remaining `IMAGE`
   medias.
5. Persist the resulting set.
6. If the set actually changed, notify every `CoversChangedObserver` with the new covers (empty slice
   means "no covers").

One **per-event observer** per catalog operation extracts `added` / `removed` from the lifecycle event
and calls `Reconcile`. These observers are the sole bridge between catalog mutations and cover
maintenance. They live in `pkg/catalog/cover_lifecycle.go` next to `CoverMaintenance`.

### Per-operation reconciliation inputs

| Event                 | Per affected album                                                                                                 |
|-----------------------|--------------------------------------------------------------------------------------------------------------------|
| `MediasInserted`      | `added = inserted medias for this album`, `removed = ∅`                                                            |
| `AlbumCreated`        | New album: `added = TransferredMedias.Transfers[newId]`, `removed = ∅`. Sources: `added = ∅`, `removed = moved-out`|
| `AlbumDatesAmended`   | Destinations: `added = moved-in`, `removed = ∅`. Sources: `added = ∅`, `removed = moved-out`                       |
| `AlbumDeleted`        | Deleted album: canonical record **deleted**, `OnCoversChanged(albumId, nil)` fired. Destinations: `added = moved-in`|
| `AlbumRenamed` (folder-change) | Canonical record **moved** from old SK to new SK; `OnCoversChanged(oldId, nil)` + `OnCoversChanged(newId, migrated)` fired |
| `AlbumRenamed` (in-place)      | No-op                                                                                                      |

## View projection

The `ListAlbums` read path serves everything from a single Query on the user's partition
(`USER#{EMAIL}#ALBUMS_VIEW`). To keep that invariant while making covers a sibling concern, each album
has **two SK rows per viewer**:

- `ALBUM#{owner}#{folder}` — identity, display fields (`Name`, `Start`, `End`), `MediaCount`.
- `ALBUM#{owner}#{folder}#COVERS` — the ordered list of up to 4 covers.

The DynamoDB adapter's `ListSummariesForUser` returns rows merged per album; a missing cover row means
"no covers". Albums with no summary row are skipped (orphan cover rows are discarded defensively).

### Why two rows, not one field

- Cover writes never touch the main summary row → no attribute-footprint discipline, no clobber risk,
  no need for atomic `ADD Count :d SET Covers = :c` updates.
- Cover maintenance stays independent from count / display-field maintenance; the view's event
  handlers for `MediasInserted`, `AlbumCreated`, `AlbumDatesAmended`, `AlbumRenamed` stop touching
  covers entirely.
- Share: writing the visitor's cover row is one `PutItem`, decoupled from the main summary row's
  write.
- Still one Query at read time — the user partition remains the single read source.

### Cover projection observer

A single observer in `pkg/catalogviews` listens on `CoversChangedObserver` and, for each notification,
fans out to every viewer who can access the album, writing (or deleting, if the new set is empty) the
`#COVERS` SK row per viewer. The fan-out uses the existing `ListUsersWhoCanAccessAlbumPort`.

The `AlbumView`'s existing per-event handlers (`OnMediasInserted`, `OnAlbumCreated`,
`OnAlbumDatesAmended`, `OnAlbumRenamed`, `OnAlbumDeleted`, `AlbumShared`, `AlbumUnShared`) are
**unchanged** by this feature. Covers propagate through their own observer chain, not through them.
`AlbumShared` and `AlbumUnShared` do write/delete the visitor's `#COVERS` SK row — see the sharing
story.

### Drift

Drift reconciliation already rebuilds view rows from the canonical catalog. It gains the
`#COVERS` SK row in its rebuild path. No new behaviour otherwise.

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
- Response: `200 OK` with the new covers list (so the UI can refresh without re-querying the album list).
- Side effect: `CoverMaintenance.Reconcile` runs with `added = ∅`, `removed = ∅` (triggers the fallback
  query path, which redraws every `RANDOM` cover).

### Phase 4+ (deferred)

A single `PUT /api/v1/owners/{owner}/albums/{folderName}/covers` endpoint will replace the per-media
pick/unpick from the earlier draft, taking the full ordered list of up to 4 covers. Not specified here.
