# 01 — Album-list view becomes a complete read model

Status: wontdo
Phase: 1
Layer: catalog domain — `pkg/catalogviews` + `pkg/catalogviewsadapters/catalogviewsdynamodb`
Depends on: —

> **Superseded** by the sub-tickets `01-01` … `01-09`. Kept here for traceability;
> the acceptance criteria below are still authoritative and are collectively fulfilled by
> the sub-tickets.

## Description

`GET /albums` currently reads media counts from the per-user view record but re-fetches each album's display
fields (Name, Start, End) from the album metadata partition through two divergent paths (owned = Query on
`{OWNER}#ALBUM`, shared = BatchGet by ids), then merges them. Promote the per-user view record to a complete
album projection so the read is served from a single query and the owned/shared split disappears.

This is a refactor with **no change to the `GET /albums` response**. It is the foundation the covers feature
builds on. Record the decision as an ADR.

See `../design.md` (Read model) for the record shape, write-path split, and read collapse.

## Acceptance criteria

- The view record (renamed `AlbumSummaryRecord`) carries `AlbumName`, `AlbumStart`, `AlbumEnd`; the SK is
  reshaped to `{OWNED|VISITOR}#{OWNER}#{FOLDER_NAME}` (the `#COUNT` suffix is dropped as the row is no
  longer count-only).
- Album create, rename-in-place, and amend-dates propagate the display fields to every viewer's record
  (owner and visitors), reusing the existing viewer enumeration. Media count updates keep the atomic
  `ADD Count {diff}` and are unaffected.
- `AlbumView.ListAlbums` is served by a single `ListSummariesForUser` query per user for both owned and
  shared albums; `FindAlbumsByOwner` and `FindAlbumsById` are removed from the read path, and the
  owned/shared providers collapse accordingly.
- The owner's "shared with" grid remains a separate query (unchanged).
- Drift control rebuilds the new display fields from the canonical album records.
- `GET /albums` returns exactly the same fields and values as before (verified by test).
- An ADR under `docs/adr/` records "album-list view is a complete read model, not just counts".
- Tests follow the Go testing strategy; existing view tests updated.

## Out of scope

- Covers (Phase 2+).
- Denormalising the sharing grid.
- Any change to the REST response.

## References

- `../spec.md`, `../design.md`
- Load skills: `go`, `architecture` (and `domain-modeling` for the ADR).
