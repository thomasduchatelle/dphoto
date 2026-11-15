# Media Tagging — Design

Technical direction spanning several issues. The "how". See [spec.md](./spec.md) for the "what".

## Context boundary

Tagging is a **catalog** concern: a tag is an attribute of a media, which is already owner-scoped in the catalog. No ACL change is needed for v1 (owner-only writes are enforced the same way album edits are; visitor visibility is out of scope).

| Layer   | Component                                       | Responsibility                                                                     |
|---------|-------------------------------------------------|------------------------------------------------------------------------------------|
| Domain  | `pkg/catalog`                                   | Tag value object (parse/normalise), add/remove tag on a media, list medias by tag. |
| Adapter | `pkg/catalogadapters` (dynamo)                  | Persist tag items on the media partition; query medias-by-tag via GSI.             |
| API     | `api/lambdas`                                   | REST endpoints to add/remove a tag and list medias by tag.                         |
| Web     | `web-nextjs/domains/catalog` + album/tags pages | Star gesture, tagging mode, Tags page + header nav.                                |

## Data model

A tag is identified by its **normalised string**: lowercased, spaces removed, URL-safe; `namespace:value` split on the **first** colon; a tag with no colon is a bare tag (reserved for `favorite`). Normalising at the boundary prevents `People:Luke` / `people:luke` duplicates and lets the tag appear directly in a URL on the Tags page.

Tags are stored as **one item per (media, tag)** on the media's partition, so they load alongside the media when an album is read (no `1+P` query). Each tag item carries the `MediaId` and the `AlbumId` (`Owner` + `FolderName`) so the Tags page and the future CLI download can resolve the media and its album without an extra lookup.

A **GSI partitioned by `{Owner}#{Tag}`** answers "all medias of a tag" in `O(matches)`, powering the Tags page and the future `dphoto download --tag`. Adding/removing a tag is one item write/delete; the GSI absorbs the write amplification — an accepted cost for cheap tag browsing.

Exact PK/SK strings, GSI name, and projected attributes are decided during the first issue (the walking skeleton) and recorded in `DATA_MODEL.md`. The decision to store tags as sibling items of the media plus a `{Owner}#{Tag}` GSI is settled here and must not be revisited per-issue.

### Distinct tags (decision owed)

The `{Owner}#{Tag}` GSI finds medias *of* a tag, not the *set* of tags an owner has used. The Tags page switcher needs that set. Decide between a **per-owner tag-registry item** maintained on tag add/remove (one cheap read to enumerate) and an **enumeration query**. The registry is favoured for cost and read simplicity. Settled in the "list the owner's distinct tags" issue, recorded here once chosen; must not require a table scan at read time.

### Future-proofing

- The per-tag item holding `MediaId` + `AlbumId` and the `{Owner}#{Tag}` GSI are exactly what the CLI download-by-tag needs: no model change required later.
- `favorite` is a plain tag in storage; only the UI special-cases it. No schema branch for favourites.

### GSI duplication

Adding the `{Owner}#{Tag}` GSI touches two places until `2026-09-retire-appdynamodb-schema` lands: the CDK construct `deployments/cdk/lib/catalog/catalog-store-construct.ts` (production) **and** the legacy Go helper `pkg/awssupport/appdynamodb/create_update_table.go` (which bootstraps local/test tables). Both must carry it despite the helper's eviction notice.

## Client/server interface (pinned contract)

Under `/api/v1`, matching existing conventions (`owner` + `mediaId` path params; authz handled by the authorizer; owner-only writes). This contract is fixed so `02.2` and `02.3` can be built against a mock in parallel with `02.1`. `{tag}` is the normalised tag string, URL-path-encoded.

### Tag type (shared JSON shape)

```json
{
  "tag": "people:luke",
  "namespace": "people",
  "value": "luke"
}
```

For a bare tag (`favorite`): `{ "tag": "favorite", "namespace": "", "value": "favorite" }`.

### Add a tag to a media

`PUT /api/v1/owners/{owner}/medias/{mediaId}/tags/{tag}` — idempotent. `204 No Content` on success.

### Remove a tag from a media

`DELETE /api/v1/owners/{owner}/medias/{mediaId}/tags/{tag}` — idempotent (removing an absent tag still `204`). `204 No Content` on success.

### List medias carrying a tag (Tags page)

`GET /api/v1/owners/{owner}/tags/{tag}/medias` → `200` with an array of medias. Each media uses the **existing `list-medias` `Media` DTO**, extended with its album so the grid can link back:

```json
[
  {
    "id": "...",
    "type": "IMAGE",
    "filename": "...",
    "time": "2024-01-02T10:00:00Z",
    "source": "...",
    "albumId": {
      "owner": "...",
      "folderName": "/2024-01_Holiday"
    },
    "tags": [
      "favorite",
      "people:luke"
    ]
  }
]
```

### Tags on album medias

The existing `GET /api/v1/owners/{owner}/albums/{folderName}/medias` response gains a `"tags": ["..."]` field (array of normalised tag strings) on each media, so the grid renders the star and active-tag state without a second round-trip. This is an additive, backward-compatible change.

### List the owner's distinct tags (story 05)

`GET /api/v1/owners/{owner}/tags` → `200` with an array of the Tag shape above, covering every distinct tag the owner has used.

```json
{
  "tags": [
    {
      "name": "people:luke"
    }
  ]
}
```

The structure of the tags as an object is to future-proof the model.

## Interaction model (cross-cutting UI decision)

Two distinct gestures, so the common case stays one click and the general case stays fast:

- **Favorite**: an always-available star overlay on each media in the album grid; one click toggles `favorite`. Independent of tagging mode.
- **Tagging mode** (slice 2): entered via `t` or a button, a type-ahead picker selects one active tag (autocomplete existing, create-on-the-fly, recently-used client-side quick-picks). While active, a header banner shows the tag and clicking a media toggles that tag on it ("paint" model). `ESC` or a button exits.

This "pick the tag once, then paint" model is the chosen alternative to server-side multi-select; multi-select is explicitly not implemented.

## Slicing

The interaction model is the risky part, so it is explored first (a throwaway prototype) and the UI is built before/independently of the backend by coding against the pinned REST contract.

1. Interaction spike — pick the tagging interaction (decision recorded below once made).
2. Favourites, end-to-end — one story in three parallel tasks once the contract is pinned: `02.1` backend (domain + persistence + GSI + write & list-by-tag endpoints), `02.2` favourite star on the album grid, `02.3` Tags page + header nav showing favourites. `02.2`/`02.3` code against the mocked contract.
3. Chosen tagging mode wired into the album page.
4. Tag switcher on the Tags page.
5. List the owner's distinct tags (feeds the switcher).
6. CLI download by tag.

## Chosen interaction

_To be filled by the interaction spike (issue 01). Until then, the paint model in the section above is the working assumption, not a decision._
