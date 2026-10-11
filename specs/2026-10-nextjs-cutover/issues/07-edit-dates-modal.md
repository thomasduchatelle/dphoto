# Edit album dates modal (component + stories)

**Status**: ready
**Depends on**: none — the component is unwired, delivered with Storybook only; wiring happens in story 01.

## Description

Deliver an `EditDatesModal` React component with Storybook stories, matching the layout and behaviour of
`web/src/components/albums/EditDatesDialog/EditDatesDialog.tsx` and driven by the props shape returned by
the already-migrated `editDatesDialogSelector`.

The component is **not** mounted into the album page by this story. The follow-up wiring story (01) is
responsible for mounting it in `AlbumPageContent` and connecting the handlers to the real thunks.

## Props contract

Props = `EditDatesDialogSelection & EditDatesDialogHandlers`.

- `EditDatesDialogSelection` is imported verbatim from
  `web-nextjs/domains/catalog/album-edit-dates/selector-editDatesDialogSelector.ts`. Current shape at time
  of writing:
  `{ isOpen, albumName, startDate, endDate, startAtDayStart, endAtDayEnd, isLoading, errorCode?,
     dateRangeError?, isSaveEnabled }`. Follow the selector export — do not duplicate the type.
- `EditDatesDialogHandlers` is a new interface declared next to the component, with one callback per thunk
  the dialog triggers:
  - `onClose(): void` — paired with `closeEditDatesDialog`.
  - `onStartDateChange(date: Date | null): void` — paired with the start-date action.
  - `onEndDateChange(date: Date | null): void` — paired with the end-date action.
  - `onStartAtDayStartChange(checked: boolean): void` — paired with the start-day-start action.
  - `onEndAtDayEndChange(checked: boolean): void` — paired with the end-day-end action.
  - `onSave(): void` — paired with `updateAlbumDates`.

Place the component at
`web-nextjs/app/(authenticated)/albums/[ownerId]/[albumId]/_components/EditDatesModal/`.

## Context

- Reference layout and behaviour: `web/src/components/albums/EditDatesDialog/EditDatesDialog.tsx`. Keep the
  same structure: linear progress bar on top, title "Edit dates of {albumName}", close icon, date range
  picker, `errorCode` as an inline alert, Cancel / Save in actions.
- `DateRangePicker` already exists duplicated inside PR #307
  (`web-nextjs/app/(authenticated)/_components/CreateAlbumModal/DateRangePicker.tsx`). **Copy it into your
  own modal folder** (`EditDatesModal/DateRangePicker.tsx`) — do not try to extract to a shared location,
  that cleanup is deferred to story 01.
- `@mui/x-date-pickers` and `dayjs` were added to `web-nextjs/package.json` by PR #307. If merge order
  turns out otherwise and they are missing, add them (same versions used by PR #307: `^8.29.4` and
  `^1.11.23`).
- Follow the `ui-components` and `nextjs` skills: Storybook for visual testing, no unit test files for
  pure UI, props must be a straight line from the selector output.

## Acceptance Criteria

- `EditDatesModal` component accepts `EditDatesDialogSelection & EditDatesDialogHandlers`.
- When `isOpen` is true, the dialog renders with both date pickers pre-filled from the selection.
- Validation: Save is disabled when `isSaveEnabled` is false; `dateRangeError` renders inline on the date
  pickers.
- `errorCode` renders as an `Alert severity="error"` at the top of the content.
- Loading state (`isLoading`) shows the top progress bar and disables inputs + Save.
- The component is a `'use client'` React component, no server-only code.
- Storybook stories cover: default pre-filled, invalid range (end before start), submission error, loading,
  "start at day start" toggled off (shows DateTimePicker instead of DatePicker).
- The dialog is responsive (full-screen on mobile) and keyboard-accessible.

## Out of scope

- Mounting the modal in `AlbumPageContent` — story 01.
- Any change to the Dates button in `AlbumHeader` / `AlbumActionsFab` — story 01.
- Orphan-media preview with thumbnails — deferred (needs a backend "medias by date" endpoint).
- Extracting `DateRangePicker` to a shared location — deferred to story 01's cleanup.
