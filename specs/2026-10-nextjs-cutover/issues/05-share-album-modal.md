# Share album modal (component + stories)

**Status**: ready
**Depends on**: none — the component is unwired, delivered with Storybook only; wiring happens in story 01.

## Description

Deliver a `ShareAlbumModal` React component with Storybook stories, matching the layout and behaviour of
`web/src/components/albums/ShareDialog/index.tsx` and driven by the props shape returned by the already-
migrated `sharingDialogSelector`.

The component is **not** mounted into the album page by this story. The follow-up wiring story (01) is
responsible for mounting it in `AlbumPageContent` and connecting the handlers to the real thunks.

## Props contract

Props = `SharingDialogFrag & ShareAlbumModalHandlers`.

- `SharingDialogFrag` is imported verbatim from
  `web-nextjs/domains/catalog/sharing/selector-sharingDialogSelector.ts`. Current shape at time of writing:
  `{ open: boolean; sharedWith: Sharing[]; suggestions: UserDetails[]; error?: ShareError }`. If the
  selector type evolves, follow it — the dialog consumes whatever the selector exports.
- `ShareAlbumModalHandlers` is a new interface declared next to the component, with one callback per thunk
  the dialog triggers:
  - `onClose(): void` — paired with `closeSharingModal`.
  - `onGrant(email: string): Promise<void>` — paired with `grantAlbumAccess`.
  - `onRevoke(email: string): Promise<void>` — paired with `revokeAlbumAccess`.

Place the component at
`web-nextjs/app/(authenticated)/albums/[ownerId]/[albumId]/_components/ShareAlbumModal/`.

## Context

- Reference layout and behaviour: `web/src/components/albums/ShareDialog/index.tsx` (full file, 230 lines).
  Keep the same structure, labels, chip-based access list, email-input with send icon, suggestions row,
  revoke-with-check animation, destructive styling for errors. The CSS file
  `ShareDialogChipsAnimation.css` can be copied over if the animation is kept.
- Album name as the dialog title / subtitle for context. Pass it as a prop if the selector exposes it; if
  not, leave the title generic ("Sharing album to ...") like the `web/` version — do **not** reach into the
  catalog state inside the component.
- Email validation on blur: inline error if the format is invalid; the send button is disabled until the
  input is a valid, non-empty email.
- Revoking triggers a confirmation — mirror the `web/` UX (chip "delete" icon with success flash on grant,
  error flash on failure). Do not immediately revoke on a single click without any feedback.
- Follow the `ui-components` and `nextjs` skills: Storybook for visual testing, no unit test files for pure
  UI, props must be a straight line from the selector output.

## Acceptance Criteria

- `ShareAlbumModal` component accepts `SharingDialogFrag & ShareAlbumModalHandlers`.
- When `open` is true, the dialog renders with the current access list (chips with avatar + name + email).
- Email input validates format on blur; invalid input disables the Grant action.
- Grant failure (error with matching email) shows the error inline on the input / suggestion chip and leaves
  the input populated for retry.
- Revoke triggers an inline confirmation / animated feedback cycle before actually revoking; errors are
  surfaced inline.
- The component is a `'use client'` React component, no server-only code.
- Storybook stories cover: closed, empty access list, with several users, suggestions visible, email
  validation error, grant error, revoke error.
- The dialog is responsive (full-screen on mobile per MUI `useMediaQuery` pattern) and keyboard-accessible.

## Out of scope

- Mounting the modal in `AlbumPageContent` — story 01.
- Any change to the Share button in `AlbumHeader` / `AlbumActionsFab` — story 01.
- Any change to the sharing thunks or state — the domain layer is already migrated.
