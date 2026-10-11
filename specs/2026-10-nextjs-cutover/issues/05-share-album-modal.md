# Share album modal

**Status**: ready
**Depends on**: `01-header-action-scaffold.md` — needs the enabled Share button and the stub modal placeholder to replace.

## Description

Replace the share stub modal on the album page with a working sharing dialog: an owner can view who has
access to the album, grant access by email, and revoke access with a confirmation step.

The domain layer is already migrated: see `web-nextjs/domains/catalog/sharing/`. This story is UI wiring and
a dialog component.

## Context

- The share button is enabled and wired to `openSharingModal` by story 01. This story replaces the stub
  modal component with the real dialog.
- Thunks and state slices to use (all under `domains/catalog/sharing/`):
  - `openSharingModal`, `closeSharingModal`
  - `grantAlbumAccess` (consumes an email)
  - `revokeAlbumAccess` (consumes a user)
  - `SharingModalErrorOccurred` already surfaces API errors into the state — render them in the dialog.
- Reference the old implementation for behaviour parity: `web/src/components/sharing-dialog/` and
  `web/src/core/catalog/sharing/`.
- Dialog contents:
  - Album name as the dialog title / subtitle for context.
  - Current access list: avatar, name, email, and a Revoke icon per entry.
  - Email input + Grant button.
  - Error area driven by `SharingModalErrorOccurred` state, with a retry affordance.
  - Close button.
- Email validation on blur: inline error if the format is invalid; Grant is disabled until the input is a
  valid, non-empty email.
- Revoking triggers an inline confirmation (dialog or inline popover) — the dialog must not immediately
  revoke on a single click.

## Acceptance Criteria

- Clicking Share in `AlbumHeader` (or `AlbumActionsFab`) opens the sharing dialog, pre-populated with the
  current access list.
- Granting access with a valid email adds the user to the list (avatar + name once the profile is loaded),
  clears the email input, and the album card on the home page updates to show the new avatar on next render.
- Granting access with an invalid email shows an inline validation error and does not call the API.
- Grant failure (e.g. unknown email) shows the error inline and leaves the input populated for retry.
- Clicking the Revoke icon on a user opens a confirmation step (user name + album name + warning + Cancel /
  Revoke). Cancel is the default focus; Revoke uses destructive styling.
- Confirming revoke removes the user from the list on success and surfaces errors inline on failure.
- Share and Revoke are **not visible** to non-owners (enforced by story 01 at the button level — this story
  just has to assume the dialog is only mounted for owners).
- Dialog is keyboard-accessible and responsive on mobile / tablet / desktop.
- Storybook stories cover: empty access list, with several users, email validation error, grant error,
  revoke confirmation open, revoke error.

## Out of scope

- Any change to viewer-side UI (how non-owners see the album) — already handled.
- Showing when the album was shared, or audit history — not a parity feature.
