# Design — Album covers

Cross-issue technical direction. Issues link here rather than repeat it. The "what" is in `spec.md`.

## Read model (Phase 1)

Today `list-albums` reads counts from the per-user view record (`AlbumSizeRecord`, SK `…#COUNT`) but
re-fetches album display fields (Name/Start/End) from the album metadata partition through two paths
(owned = Query, shared = BatchGet). We promote the view record to a **complete projection**.

- **Record**: keep PK `USER#{EMAIL}#ALBUMS_VIEW` and SK `{OWNED|VISITOR}#{OWNER}#{FOLDER_NAME}#COUNT`
  (the `#COUNT` suffix stays for backward compatibility). Extend `AlbumSizeRecord` with `AlbumName`,
  `AlbumStart`, `AlbumEnd` (and later `Covers`, see below).
- **Write path is split and must not clobber**: media insert/delete keeps the atomic `ADD Count {diff}`;
  album create/rename/amend-dates `SET`s the display fields on every viewer record (owner + visitors, via
  `ListUsersWhoCanAccessAlbum`). The two update kinds touch disjoint attributes of the same item.
- **Read path collapses**: `ListAlbums` is served by a single `GetAvailabilitiesByUser` Query (owned +
  shared). The owned/shared providers and `FindAlbumsByOwner`/`FindAlbumsById` leave the read path.
- **Sharing grid unchanged**: the owner's "shared with" badges stay a separate query (out of scope to
  denormalise).
- **Drift control** rebuilds the display fields (and covers) from canonical records, so the view stays a
  disposable projection.

This decision is durable and repo-wide → **ADR** under `docs/adr/`.

## Cover canonical record

The whole cover set of an album is a **single** DynamoDB item:

- PK `{OWNER}#ALBUM`, SK `ALBUM#{FOLDER_NAME}#COVERS`.
- Holds an ordered list of at most 4 covers, each `{MediaId, Filename, Origin}` where `Origin` is
  `RANDOM` or `CHERRY_PICKED`. Order is by capture date but not significant.
- Every write reads the set, enforces the max of 4, and rewrites the whole item.

Document it in `DATA_MODEL.md`.

## Completion operation (random fill)

One catalog operation fills empty slots up to 4 with uniform-random eligible medias (`MediaType IMAGE`)
that are not already covers, tagged `RANDOM`. Two entry shapes:

- `CompleteCoversFromCandidates(albumId, candidates)` — fills from a supplied set of medias, no extra query.
  Used by **backup** (passes the images it just inserted).
- `CompleteCovers(albumId)` — queries the album's `IMAGE` medias then completes. Used by **backfill**.

Full sets are left untouched (cheap no-op).

## Cover projection in the view

`AlbumSizeRecord` gains a `Covers` attribute mirroring the canonical set (`{MediaId, Filename, Origin}`).
Any change to an album's covers re-denormalises the set to all viewer records (same fan-out as counts).
`VisibleAlbum` carries `Covers` for the API.

## REST contract

`GET /api/v1/albums` — each album in the response gains:

```json
"covers": [
  { "mediaId": "...", "filename": "...", "origin": "RANDOM" }
]
```

`origin` ∈ `RANDOM | CHERRY_PICKED`. `covers` is `[]` when the album has none (0–4 entries). The frontend
builds the image URL from `owner` + `mediaId` + `filename` via its existing image loader
(`/api/v1/owners/{owner}/medias/{mediaId}/{filename}?w=…`) — no URL is stored or returned.

Mutation endpoints (finalised in their phase, owner-edit permission only):

| Phase | Action        | Route                                                              |
|-------|---------------|--------------------------------------------------------------------|
| 3     | Re-randomise  | `POST /api/v1/owners/{owner}/albums/{folderName}/covers/refresh`   |
| 4     | Pick          | `PUT /api/v1/owners/{owner}/albums/{folderName}/covers/{mediaId}`  |
| 5     | Unpick        | `DELETE /api/v1/owners/{owner}/albums/{folderName}/covers/{mediaId}`|

Pick refuses with an error (→ snackbar) when all 4 covers are already `CHERRY_PICKED`; otherwise it evicts
a `RANDOM` cover when the set is full.
