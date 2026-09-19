# 01-09 — Cleanup, ADR, DATA_MODEL, acceptance tests

Status: ready
Phase: 1
Layer: catalog domain — `pkg/catalogviews` + `pkg/pkgfactory` + `docs/adr` + `DATA_MODEL.md` + pkg acceptance suite
Depends on: 01-03, 01-04, 01-05, 01-06, 01-07, 01-08

## Description

Retire the now-dead code, publish the ADR, update `DATA_MODEL.md`, and extend the pkg acceptance
tests to assert `ListAlbums` end-to-end.

## Acceptance criteria

- **Deletion / retirement**:
  - `pkg/catalogviews/command_handler_counts.go` (`CommandHandlerAlbumSize` and its methods) is
    deleted; the `TODO` on line 196 is also gone.
  - Any port left orphaned by `01-02` (e.g. `FindAlbumByOwnerPort` on the read side, the old
    `providers.go` file if it still exists) is removed.
  - Factory `pkg/pkgfactory/factory_catalog_view.go::CommandHandlerAlbumSize` is deleted; all its
    call sites in `factory_catalog.go`, `factory_backup.go`, `factory_acl.go` now reference the
    per-event adapters delivered by `01-03` … `01-07`. The `AlbumCreatedAsTimelineMutation` /
    `AlbumDeletedAsTimelineMutation` / `AlbumRenamedAsTimelineMutation` /
    `AlbumDatesAmendedAsTimelineMutation` adapters used to bridge the old
    `TimelineMutationObserver` into the new per-use-case events become dead once the count logic
    has moved into the AlbumView per-use-case handlers — remove them too.
- **DATA_MODEL.md**: the two `USER#{EMAIL}#ALBUMS_VIEW / …#COUNT` rows are updated:
  - SK column: `OWNED#{OWNER}#{FOLDER_NAME}` / `VISITOR#{OWNER}#{FOLDER_NAME}` (no `#COUNT`).
  - Description column: "(view) album summary for an album owned by the user: count + display
    fields (name/start/end)" and the visitor equivalent.
- **ADR**: `docs/adr/0004-album-list-view-is-a-complete-read-model.md` created, following the shape
  of ADR-0002/0003:
  - **Context**: today's `GET /albums` merges the per-user count row with a per-album metadata
    fetch through two divergent paths (Query for owned, BatchGet for shared).
  - **Decision**: the per-user view record is a complete projection of the album (identity + display
    + count); the write path is split into disjoint attribute sets so count updates and
    display-field updates never clobber each other; the read path collapses to a single Query per
    user (plus one Query for the owner's sharing grid, unchanged); drift reconciliation rebuilds
    display fields from the canonical records.
  - **Consequences**: viewer records must be fanned out on album create / rename-in-place /
    amend-dates / share; `list-albums` no longer needs the album metadata partition; owned/shared
    read providers are gone; a one-off post-deploy drift reconciliation backfills legacy rows.
- **Pkg acceptance tests** (existing suite under `pkg/`): extend to call `pkgfactory.AlbumView(...)`
  `.ListAlbums` after each of the following real operations (against Docker DynamoDB / whatever the
  suite already uses) and assert the returned `VisibleAlbum` slice reflects the change:
  - After `CreateAlbum` → the owner sees the album with `Name/Start/End` and `MediaCount=0`.
  - After media insert → the owner's `MediaCount` reflects the insert.
  - After `ShareAlbum` → the visitor sees the album with `Name/Start/End` and the current count.
  - After `AmendAlbumDates` → both owner and visitor see the new dates.
  - After `RenameAlbum` (in place) → both see the new name.
  - After `DeleteAlbum` → neither sees the album.
- The `list-albums` lambda in `api/lambdas/list-albums` is unchanged: DTO fields identical, values
  identical (a smoke test may be added if the project has a convention; if not, mention the
  manual/curl verification in the PR description).

## Out of scope

- Anything related to covers (Phase 2+).
- Denormalising the sharing grid.
- CLI/admin tooling changes beyond what `01-08` already documents.

## References

- `../spec.md`, `../design.md`
- Parent ticket `01-album-list-read-model.md` (interface contract).
- Load skills: `go`, `architecture`, `domain-modeling` (for the ADR).
