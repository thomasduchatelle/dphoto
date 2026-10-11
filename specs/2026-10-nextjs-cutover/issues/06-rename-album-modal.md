# Rename album modal

**Status**: ready
**Depends on**: `01-header-action-scaffold.md` — needs the enabled Rename button, the stub modal placeholder, and the open/close thunks that 01 creates for the rename slice.

## Description

Replace the rename stub modal on the album page with a working rename dialog: an owner can change the
album's display name.

The domain layer is already migrated: see `web-nextjs/domains/catalog/album-edit-name/`. This story is UI
wiring and a dialog component.

## Context

- The Rename button is enabled and wired to the open-modal thunk by story 01. This story replaces the stub
  modal component with the real dialog.
- Thunks and state slices to use (all under `domains/catalog/album-edit-name/`):
  - Open / close modal thunks.
  - Save thunk that calls the backend rename endpoint.
  - Error action/state slice that surfaces backend errors into the dialog.
- Reference the old implementation for behaviour parity:
  `web/src/components/edit-album-name-dialog/` and `web/src/core/catalog/album-edit-name/`.
- Dialog contents:
  - Album name text field, pre-filled with the current name.
  - Context line showing album metadata (date range + photo count).
  - Cancel and Save buttons.
- Validation: name must be non-empty (trimmed); Save disabled until valid.

## Acceptance Criteria

- Clicking Rename in `AlbumHeader` (or `AlbumActionsFab`) opens the dialog with the name pre-filled and
  focus set to the name field.
- Changing the name and clicking Save:
  - Shows a loading state on Save.
  - On success: closes the dialog; the new name appears immediately in the album page header, in the
    browser tab title (if applicable), and on the home page album card.
  - On failure: shows the error inline inside the dialog with a retry.
- Clearing the name disables Save and shows a validation error below the field.
- Cancel closes the dialog without saving; pressing ESC does the same.
- Dialog is keyboard-accessible (Tab, Enter to submit, ESC to cancel) and responsive on mobile / tablet /
  desktop.
- Storybook stories cover: pre-filled with a short name, pre-filled with a long name, empty / validation
  error, submission error, loading.

## Out of scope

- Changing the album's underlying folder name — not a parity feature; the backend endpoint for rename only
  changes the display name.
