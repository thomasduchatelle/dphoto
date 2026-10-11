# Delete album modal (component + stories)

**Status**: ready
**Depends on**: none — the component is unwired, delivered with Storybook only; wiring happens in story 01.

## Description

Deliver a `DeleteAlbumModal` React component with Storybook stories, matching the layout and behaviour of
`web/src/components/albums/DeleteAlbumDialog/DeleteAlbumDialog.tsx` and driven by the props shape returned
by the already-migrated `deleteDialogSelector`.

The component is **not** mounted into the album page by this story. The follow-up wiring story (01) is
responsible for mounting it in `AlbumPageContent` and connecting the handlers to the real thunks.

## Props contract

Props = `DeleteDialogFrag & DeleteAlbumModalHandlers`.

- `DeleteDialogFrag` is imported verbatim from
  `web-nextjs/domains/catalog/album-delete/selector-deleteDialogSelector.ts`. Current shape at time of
  writing:
  `{ albums: Album[]; initialSelectedAlbumId?: AlbumId; isOpen: boolean; isLoading: boolean; error?: string }`.
  Follow the selector export — do not duplicate the type.
- `DeleteAlbumModalHandlers` is a new interface declared next to the component, with one callback per
  thunk the dialog triggers:
  - `onClose(): void` — paired with `closeDeleteAlbumDialog`.
  - `onDelete(albumId: AlbumId): void` — paired with `deleteAlbum`.

Place the component at
`web-nextjs/app/(authenticated)/albums/[ownerId]/[albumId]/_components/DeleteAlbumModal/`.

## Context

- Reference layout and behaviour: `web/src/components/albums/DeleteAlbumDialog/DeleteAlbumDialog.tsx` (full
  file, 184 lines). Keep the same structure and two-step UX: a Select of the deletable albums (pre-selected
  from `initialSelectedAlbumId`), then a confirmation step with the formatted album label. Linear progress
  bar on top for loading. `error` renders as an inline alert.
- Destructive styling: the Delete / Confirm buttons use `color="error"` + `variant="contained"`. Cancel is
  the default focus.
- The confirmation step is internal to the component (local `useState` for `confirmMode` and
  `selectedAlbumIndex`) — do **not** push this state into the domain layer. The `web/` implementation does
  the same thing.
- `formatAlbumLabel` can be copied verbatim from the `web/` implementation, including its dependency on
  `toLocaleDateWithDay`. Check whether `web-nextjs` already has a date-utility; if not, inline the small
  helper in the modal folder.
- Follow the `ui-components` and `nextjs` skills: Storybook for visual testing, no unit test files for
  pure UI, props must be a straight line from the selector output.

## Acceptance Criteria

- `DeleteAlbumModal` component accepts `DeleteDialogFrag & DeleteAlbumModalHandlers`.
- When `isOpen` is true, the dialog renders with the album Select pre-populated from `albums` and the
  selection initialized from `initialSelectedAlbumId`.
- Clicking Delete transitions to the confirmation step, showing the formatted label of the selected album;
  Confirm calls `onDelete(selectedAlbum.albumId)`.
- Cancel closes the dialog; ESC does the same.
- Loading state (`isLoading`) shows the top progress bar and disables Cancel + the primary action.
- `error` renders as an `Alert severity="error"` at the top of the content; the album stays visible for
  retry.
- The component is a `'use client'` React component, no server-only code.
- Storybook stories cover: default (selection step), confirmation step, loading, submission error.
- The dialog is responsive and keyboard-accessible.

## Out of scope

- Mounting the modal in `AlbumPageContent` — story 01.
- Any change to the Delete button in `AlbumHeader` / `AlbumActionsFab` — story 01.
- Any "undo" / soft-delete UX — the backend does not support it.
