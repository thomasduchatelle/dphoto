# Create album modal

**Status**: ready
**Depends on**: none — lives on the album list page, independent of the album-header scaffolding (01) and the viewer (02/03).

## Description

Replace the current `onCreateAlbum` placeholder on the album list page with a working modal: owners can
create a new album by entering a name, a start date, an end date, and optionally a custom folder name.

The domain layer is already migrated: see `web-nextjs/domains/catalog/album-create/`. This story is UI
wiring and a dialog component.

## Context

- The button lives on the home page via `AlbumGrid` (`onCreateAlbum` prop, already present). This story
  replaces the current handler so clicking the button opens the modal.
- The create thunks dispatch `openCreateAlbumDialog`, `createAlbum`, `closeCreateAlbumDialog`-style actions
  that already exist in `domains/catalog/album-create/`. Follow their contract instead of inventing new
  state.
- The modal must be **owner-only**: non-owners don't see the button at all. Ownership is already enforced
  where the button is rendered — do not duplicate.
- Reference the old implementation for behaviour parity: `web/src/components/create-album-dialog/` and the
  surrounding thunks in `web/src/core/catalog/album-create/`.
- Validation:
  - Name is required (non-empty, trimmed).
  - End date must be strictly after start date.
  - Custom folder name is optional; when the toggle is off, the backend derives the folder name.
- Errors from the backend (e.g. folder-name conflict, forbidden) surface inline in the dialog with a retry
  path.

## Acceptance Criteria

- Clicking the "Create album" button on the home page opens the create-album dialog.
- The dialog contains:
  - Album name text field (required).
  - Start date picker.
  - End date picker.
  - Toggle for custom folder name; when enabled, a text field appears below.
  - Cancel and Create buttons.
- Inline validation errors are shown below the relevant field; the Create button is disabled until the form
  is valid.
- Submitting:
  - Shows a loading state on the Create button.
  - On success: closes the dialog and the new album appears in the list in chronological order.
  - On failure: shows the error inline inside the dialog with a retry (the Create button becomes active
    again) and does not close the dialog.
- Cancel closes the dialog without creating anything.
- Dialog is keyboard-navigable (Tab, Enter to submit, ESC to cancel) and responsive on mobile / tablet /
  desktop.
- Storybook stories cover: default (empty), filled valid, validation error, submission error, loading.

## Out of scope

- Visual photo preview of what the album will contain — deferred (needs a backend "medias by date" endpoint).
- Showing media density indicators on the date pickers — deferred.
