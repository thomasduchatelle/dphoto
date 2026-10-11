# Edit album dates modal

**Status**: ready
**Depends on**: `01-header-action-scaffold.md` — needs the enabled Dates button, the stub modal placeholder, and the open/close thunks that 01 creates for the edit-dates slice.

## Description

Replace the edit-dates stub modal on the album page with a working date-range editor: an owner can change
the album's start and end dates, with the backend re-indexing the medias accordingly.

The domain layer is already migrated: see `web-nextjs/domains/catalog/album-edit-dates/`. This story is UI
wiring and a dialog component.

## Context

- The Dates button is enabled and wired to the open-modal thunk by story 01. This story replaces the stub
  modal component with the real dialog.
- Thunks and state slices to use (all under `domains/catalog/album-edit-dates/`):
  - Open / close modal thunks.
  - Save thunk that calls the backend endpoint and triggers re-indexing.
  - Error action/state slice that surfaces backend errors (including orphan warnings) into the dialog.
- Reference the old implementation for behaviour parity:
  `web/src/components/edit-album-dates-dialog/` and `web/src/core/catalog/album-edit-dates/`.
- Dialog contents:
  - Context line with the current album name.
  - Start date picker, pre-filled.
  - End date picker, pre-filled.
  - Cancel and Save buttons.
- Validation: end date must be strictly after start date; Save disabled otherwise.
- The backend may return an orphan-media warning (the new range would exclude medias currently in the
  album). The dialog surfaces this as an inline warning with details and a confirmation to proceed (either
  a second click on Save, or an explicit "Proceed anyway" button — mirror the old UX).

## Acceptance Criteria

- Clicking Dates in `AlbumHeader` (or `AlbumActionsFab`) opens the dialog with both dates pre-filled.
- Adjusting dates to an invalid range (end ≤ start) disables Save and shows an inline validation message.
- Submitting a valid range:
  - Shows a loading state on Save.
  - On success: closes the dialog; the album page refreshes to show the updated media set and the header
    reflects the new date range.
  - On failure: shows the error inline inside the dialog with a retry.
  - If the backend returns an orphan-media warning: show the warning with details inside the dialog; the
    user must explicitly acknowledge before the change is applied.
- Cancel closes the dialog without saving; pressing ESC does the same.
- Dialog is keyboard-accessible (Tab to move, arrow keys inside the date picker, Enter to submit, ESC to
  cancel) and responsive on mobile / tablet / desktop.
- Storybook stories cover: default pre-filled, invalid range, orphan warning, submission error, loading.

## Out of scope

- Visual photo preview of what the album would contain with the new range — deferred (needs a backend
  "medias by date" endpoint).
- Live preview of the orphaned medias as thumbnails — deferred.
