# Delete album modal

**Status**: ready
**Depends on**: `01-header-action-scaffold.md` — needs the Delete button that 01 adds, the stub modal placeholder, and the open/close thunks that 01 creates for the delete slice.

## Description

Replace the delete stub modal on the album page with a working confirmation dialog: an owner can delete the
album they are viewing.

The domain layer is already migrated: see `web-nextjs/domains/catalog/album-delete/`. This story is UI
wiring and a confirmation dialog.

## Context

- The Delete button is added and wired to the open-modal thunk by story 01. This story replaces the stub
  modal component with the real confirmation dialog.
- Thunks and state slices to use (all under `domains/catalog/album-delete/`):
  - Open / close modal thunks.
  - Delete thunk that calls the backend endpoint.
  - Error action/state slice that surfaces backend errors into the dialog.
- Reference the old implementation for behaviour parity:
  `web/src/components/delete-album-dialog/` and `web/src/core/catalog/album-delete/`.
- Deletion in the DPhoto backend removes the album, not its medias — the medias remain in the catalog and
  can be re-indexed under another album. The confirmation message must say this clearly to match the
  current `web/` wording and avoid alarming users.
- The dialog must be visually destructive: Delete button uses `color="error"` or equivalent; focus defaults
  to Cancel.

## Acceptance Criteria

- Clicking Delete in `AlbumHeader` (or `AlbumActionsFab`) opens the confirmation dialog.
- The dialog contains:
  - Warning that the album will be removed (medias remain in the catalog).
  - Album name for context.
  - Cancel and Delete buttons; Delete uses destructive styling; focus defaults to Cancel.
- Clicking Delete:
  - Shows a loading state on Delete.
  - On success: closes the dialog, navigates the user to the home page, and the album disappears from the
    album list.
  - On failure: shows the error inline inside the dialog with a retry; the album remains.
- Clicking Cancel closes the dialog without deleting; pressing ESC does the same.
- Pressing Enter inside the dialog triggers Cancel (safer default), not Delete.
- Dialog is responsive on mobile / tablet / desktop.
- Storybook stories cover: default, loading, submission error.

## Out of scope

- Any "undo" / soft-delete mechanism — the backend does not support it and parity does not require it.
- Deleting medias alongside the album — explicitly not what delete-album does.
