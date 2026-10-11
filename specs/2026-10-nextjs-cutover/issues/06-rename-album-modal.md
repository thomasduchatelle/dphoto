# Rename album modal (component + stories)

**Status**: ready
**Depends on**: none — the component is unwired, delivered with Storybook only; wiring happens in story 01.

## Description

Deliver a `RenameAlbumModal` React component with Storybook stories, matching the layout and behaviour of
`web/src/components/albums/EditNameDialog/EditNameDialog.tsx` and driven by the props shape returned by the
already-migrated `editNameDialogSelector`.

The component is **not** mounted into the album page by this story. The follow-up wiring story (01) is
responsible for mounting it in `AlbumPageContent` and connecting the handlers to the real thunks.

## Props contract

Props = `EditNameDialogSelection & EditNameDialogHandlers`.

- `EditNameDialogSelection` is imported verbatim from
  `web-nextjs/domains/catalog/album-edit-name/selector-editNameDialogSelector.ts`. Current shape at time of
  writing:
  `{ isOpen, originalName, albumName, customFolderName, isCustomFolderNameEnabled, isLoading,
     technicalError?, isSaveEnabled, nameError?, folderNameError? }` (via `BaseEditNameSelection`). Follow
  the selector export — do not duplicate the type.
- `EditNameDialogHandlers` is a new interface declared next to the component, with one callback per thunk
  the dialog triggers:
  - `onClose(): void` — paired with `closeEditNameDialog`.
  - `onAlbumNameChange(name: string): void` — paired with the name-change action.
  - `onFolderNameEnabledChange(enabled: boolean): void` — paired with the folder-name-enabled action.
  - `onFolderNameChange(folderName: string): void` — paired with the folder-name-change action.
  - `onSave(): void` — paired with `saveAlbumName`.

Place the component at
`web-nextjs/app/(authenticated)/albums/[ownerId]/[albumId]/_components/RenameAlbumModal/`.

## Context

- Reference layout and behaviour: `web/src/components/albums/EditNameDialog/EditNameDialog.tsx`. Keep the
  same structure: linear progress bar on top for loading, title "Edit Name", close icon, name text field,
  folder-name input with the switch toggle, Cancel / Save in actions.
- `FolderNameInput` already exists duplicated inside PR #307
  (`web-nextjs/app/(authenticated)/_components/CreateAlbumModal/FolderNameInput.tsx`). **Copy it into your
  own modal folder** (`RenameAlbumModal/FolderNameInput.tsx`) — do not try to extract to a shared location,
  that cleanup is deferred to story 01.
- Follow the `ui-components` and `nextjs` skills: Storybook for visual testing, no unit test files for pure
  UI, props must be a straight line from the selector output.

## Acceptance Criteria

- `RenameAlbumModal` component accepts `EditNameDialogSelection & EditNameDialogHandlers`.
- When `isOpen` is true, the dialog renders with the name field pre-filled from `albumName`.
- Validation: Save is disabled when `isSaveEnabled` is false; `nameError` / `folderNameError` render inline
  under the matching inputs.
- `technicalError` renders as an `Alert severity="error"` at the top of the content.
- Loading state (`isLoading`) shows the top progress bar and disables inputs + Save.
- The component is a `'use client'` React component, no server-only code.
- Storybook stories cover: short name pre-filled, long name pre-filled, validation error, submission error,
  loading, with custom folder name enabled.
- The dialog is responsive (full-screen on mobile) and keyboard-accessible (ESC closes, Tab moves between
  fields).

## Out of scope

- Mounting the modal in `AlbumPageContent` — story 01.
- Any change to the Rename button in `AlbumHeader` / `AlbumActionsFab` — story 01.
- Extracting `FolderNameInput` to a shared location — deferred to story 01's cleanup.
