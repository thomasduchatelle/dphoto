# Enable owner-action buttons on the album page with stub modals

**Status**: done
**Depends on**: none

## Description

The album page's `AlbumHeader` (desktop) and `AlbumActionsFab` (mobile) already present Share / Rename /
Dates buttons as **disabled**. This story enables them for album owners, adds a **Delete** button next to
them, and wires each to a stub modal component so that the modal stories (05, 06, 07, 08) can land
independently without colliding on the same files.

The thunks these actions dispatch already exist in `web-nextjs/domains/catalog/`:

- Sharing: `openSharingModal` / `closeSharingModal` in `domains/catalog/sharing/`
- Rename: thunks in `domains/catalog/album-edit-name/`
- Edit dates: thunks in `domains/catalog/album-edit-dates/`
- Delete: thunks in `domains/catalog/album-delete/`

If the open/close thunks for rename / edit-dates / delete do not yet exist, create them following the same
pattern as the sharing modal thunks (`thunk-openSharingModal.ts`, `thunk-closeSharingModal.ts`,
`action-sharingModalOpened.ts`, `action-sharingModalClosed.ts`) and extend the catalog state in
`domains/catalog/language/` with the matching open/error fields.

## Context

- `AlbumHeader` lives at
  `web-nextjs/app/(authenticated)/albums/[ownerId]/[albumId]/_components/AlbumHeader/index.tsx`. The four
  action buttons are currently rendered as `<Button ... disabled>`.
- `AlbumActionsFab` lives next to it and mirrors the actions as a mobile speed-dial.
- Album ownership is derivable from the catalog state / the authenticated user context already used by the
  page.
- Share / rename / dates / delete must be **owner-only**: viewers must not see the buttons at all (not just
  disabled).
- A stub modal is a component that is mounted when the corresponding `modalOpen` state is true, renders a
  dialog with a title ("Share / Rename / Edit dates / Delete — not implemented yet"), a Close button wired
  to the matching `closeXxxModal` thunk, and nothing else. This stub will be replaced by the full modal in
  stories 05 / 06 / 07 / 08.

## Acceptance Criteria

- Share / Rename / Dates buttons in `AlbumHeader` are no longer `disabled` and are **only rendered when the
  authenticated user owns the album**.
- A **Delete** button is added next to them in `AlbumHeader`, with destructive styling (`color="error"` or
  equivalent), also owner-only.
- `AlbumActionsFab` on mobile mirrors all four owner-only actions.
- Each button dispatches the matching `openXxxModal` thunk when clicked.
- Four stub modal components are created, one per action, mounted on the album page and controlled by their
  respective open-state slice. Each stub renders a dialog with a title, a short "Not implemented yet"
  message, and a Close button wired to `closeXxxModal`.
- Clicking a button opens the corresponding stub modal; closing it dismisses it; the album page remains
  interactive underneath.
- Storybook stories for `AlbumHeader` cover both viewer and owner states (viewer: no buttons; owner: all
  four buttons visible).
- Open/close thunks and state slices for rename / edit-dates / delete exist and are tested (following the
  sharing pattern), if they did not already.

## Out of scope

- Any real modal content — stories 05, 06, 07, 08 replace the stubs.
- Changes to the create-album button on the home page — handled by story 04.
- Any change to the fullscreen media viewer — handled by stories 02 and 03.
