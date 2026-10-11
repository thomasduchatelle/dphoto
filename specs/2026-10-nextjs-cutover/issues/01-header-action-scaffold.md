# Wire album-header modals into the album page

**Status**: ready
**Depends on**: `05-share-album-modal.md`, `06-rename-album-modal.md`, `07-edit-dates-modal.md`,
`08-delete-album-modal.md` — this story mounts the modal components those stories deliver.

## Description

Final wiring step for the album page's owner-only actions. Enables the Share / Rename / Dates buttons in
`AlbumHeader` (currently disabled), adds a Delete button next to them, mirrors them in `AlbumActionsFab`,
and mounts the four modal components (`ShareAlbumModal`, `RenameAlbumModal`, `EditDatesModal`,
`DeleteAlbumModal`) inside `AlbumPageContent`, each driven by its selector and the matching thunks.

All four actions are **owner-only**: the buttons must not be rendered for non-owners (not just disabled).

## Context

- `AlbumHeader` lives at
  `web-nextjs/app/(authenticated)/albums/[ownerId]/[albumId]/_components/AlbumHeader/index.tsx`. The Share
  / Rename / Dates buttons are currently rendered as `<Button ... disabled>`. Delete is not yet present.
- `AlbumActionsFab` is next to it and mirrors the actions as a mobile speed-dial. Mirror the same owner-only
  logic there, including a new Delete entry.
- Album ownership is derivable from the catalog state / the authenticated user context already used by the
  page.
- The four modal components live under
  `web-nextjs/app/(authenticated)/albums/[ownerId]/[albumId]/_components/` after their respective PRs
  (05–08). Their props are `<SelectorFrag> & <Handlers>` so wiring is: pass the selector output as a
  spread, then one handler per thunk.
- Selectors already exist in `web-nextjs/domains/catalog/`:
  - `sharing/selector-sharingDialogSelector.ts` → `sharingDialogSelector`
  - `album-edit-name/selector-editNameDialogSelector.ts` → `editNameDialogSelector`
  - `album-edit-dates/selector-editDatesDialogSelector.ts` → `editDatesDialogSelector`
  - `album-delete/selector-deleteDialogSelector.ts` → `deleteDialogSelector`
- Thunks already exist in the matching slice folders. Follow the pattern used in `HomeContent` for the
  create-album flow (PR #307) when wiring the client adapter / thunks.
- Shared helpers duplicated across modal folders by stories 06 / 07 (and 04 for create) may be unified now,
  or left duplicated and deferred to a follow-up — the subagent should propose the simpler option in the PR
  description.

## Acceptance Criteria

- Share / Rename / Dates buttons in `AlbumHeader` are no longer `disabled` and are **only rendered when the
  authenticated user owns the album**.
- A **Delete** button is added next to them in `AlbumHeader`, with destructive styling (`color="error"` or
  equivalent), also owner-only.
- `AlbumActionsFab` on mobile mirrors all four owner-only actions.
- Each button dispatches the matching open-modal thunk on click.
- `AlbumPageContent` mounts the four modal components. Each receives: the output of its selector as a
  spread, plus one handler prop per thunk the modal exposes (e.g. `onClose`, `onSave`, `onGrant`, etc.).
- Clicking a button opens its modal with the correct state; closing dismisses it; the album page stays
  interactive underneath.
- Storybook stories for `AlbumHeader` cover both viewer (no buttons) and owner (all four buttons visible)
  states.
- Any shared helpers (`DateRangePicker`, `FolderNameInput`) duplicated across modal folders are either
  consolidated into `web-nextjs/components/` **or** left duplicated with an explicit note in the PR
  description — pick the simpler option, don't force the extraction if it bloats the diff.

## Out of scope

- Any change to the modal components' internal UI — handled by stories 05–08.
- Any change to the create-album button on the home page — handled by story 04.
- Any change to the fullscreen media viewer — handled by stories 02 and 03.
