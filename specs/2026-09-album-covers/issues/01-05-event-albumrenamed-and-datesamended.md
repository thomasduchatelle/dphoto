# 01-05 — Events: AlbumRenamedInPlace & AlbumDatesAmended

Status: done
Phase: 1
Layer: catalog domain — `pkg/catalog` + `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 01-01, 01-02

## Description

Fill in the two display-field-mutation event methods and wire them into the corresponding catalog
flows. Both must operate on the display attributes alone (never on `Count`).

Special care on **rename**: `catalog.RenameAlbum` has two branches — an *in-place* rename (name only,
folder unchanged) that today fires **no observer**, and a *folder-changing* rename that reuses the
create+delete path. This ticket introduces a proper observer for the in-place branch so the projection
stays in sync; the folder-changing branch is already covered by `01-03` (create) and `01-06` (delete).

## Acceptance criteria

- `AlbumView.AlbumRenamedInPlace(ctx, albumId catalog.AlbumId, newName string) error` looks up viewers
  via `ListUsersWhoCanAccessAlbumPort` and calls `Repository.SetDisplayFields` with the new name
  (Start/End kept at their current values by passing the existing zero-marker or by loading the
  album; implementation detail — the acceptance is that Start/End are not corrupted).
- `AlbumView.AlbumDatesAmended(ctx, update catalog.DatesUpdate) error` looks up viewers and calls
  `Repository.SetDisplayFields` with the album's Name and the new Start/End from
  `update.UpdatedAlbum`.
- New `catalog.RenameAlbumObserver` firing on the **in-place** branch of `catalog.RenameAlbum` (i.e.
  when `request.RenameFolder == false && request.ForcedFolderName == ""`), with the signature
  `OnAlbumRenamed(ctx, albumId, newName)`. The folder-changing branch stays as-is.
- Two adapters in `pkg/catalogviews`:
  - `AlbumViewRenameObserver` implementing the new `catalog.RenameAlbumObserver` by forwarding to
    `AlbumView.AlbumRenamedInPlace`.
  - `AlbumViewAmendDatesObserver` implementing `catalog.AlbumDatesAmendedObserver` by forwarding to
    `AlbumView.AlbumDatesAmended`.
- Factories: `RenameAlbumCase` and `AmendAlbumDatesCase` in `pkg/pkgfactory/factory_catalog.go`
  updated to pass the new adapters.
- Catalog tests (`pkg/catalog/album_rename_test.go`): add a case
  `it_should_call_the_rename_observer_when_updating_the_name_in_place` — currently the in-place
  branch has no observer coverage.
- `AlbumView` tests (in-memory fake):
  - `TestAlbumView_AlbumRenamedInPlace/it_should_update_the_name_on_all_viewer_rows` — seed owner +
    visitor, rename, assert both rows show the new name and same Start/End/Count.
  - `TestAlbumView_AlbumDatesAmended/it_should_update_start_and_end_on_all_viewer_rows` — seed owner
    + visitor, amend, assert both rows show new dates and same Name/Count.

## Out of scope

- Folder-changing rename branch (already covered by create+delete).
- Media re-transfer triggered by amend-dates (already covered by `01-04`'s `MediasTransferred`).
- Delete / share / unshare / drift.

## References

- `../spec.md`, `../design.md`
- Parent ticket `01-album-list-read-model.md` (interface contract).
- Load skills: `go`, `architecture`.
