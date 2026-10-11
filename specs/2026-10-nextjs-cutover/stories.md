# NextJS Cutover — stories

Draft of the ten issues. Each becomes `issues/NN-<slug>.md`. Parallelism target: issues 02–08 can all be
picked up in parallel from day 1 (modal components are unwired, Storybook-driven); issue 01 is the final
wiring step that mounts them on the album page.

## Parallelism map

```
(day 1, in parallel)   02 fullscreen-viewer-yarl ──── 03 viewer-capabilities
                       04 create-album-modal  (independent, lives on home page)
                       05 share-album-modal           (component + stories only)
                       06 rename-album-modal          (component + stories only)
                       07 edit-dates-modal            (component + stories only)
                       08 delete-album-modal          (component + stories only)
                       09 client-token-refresh        (soft-blocks 05–08; discovered during 04)

(after 05-08 merge)    01 wire-album-header-modals    (mounts them on the album page)

(after everything)     10 cutover
```

Issues 02 and 03 remain sequential (03 depends on 02). Modals 05–08 are **component-only** stories: each
delivers a pure React component plus Storybook stories, with props = `<DialogSelection> & <Handlers>` so the
component matches the exact shape returned by the existing selectors in `web-nextjs/domains/catalog/`. They
are **not** mounted into the album page by these stories — mounting is done by 01 in a single pass once all
four modals exist.

## Design decisions

- **No shared "modal framework"**. Each modal is its own component, its props derived from its own
  domain-layer selector. No cross-story coordination hazard.
- **Each modal story delivers only the component + Storybook**. Props shape: `SelectorFrag & XxxHandlers`
  where `SelectorFrag` comes verbatim from the matching `selector-xxxDialogSelector.ts` in
  `web-nextjs/domains/catalog/<slice>/`, and `XxxHandlers` is a sibling interface declaring one callback
  per thunk. No wiring into `AlbumHeader`, `AlbumActionsFab`, or `AlbumPageContent`.
- **The old `web/` dialog in `web/src/components/albums/<Xxx>Dialog/` is the layout reference** for its
  NextJS equivalent. Keep the same structure, labels, destructive styling, loading/error affordances. Allow
  minor polish only where MUI v7 or NextJS idioms differ.
- **Shared helpers (`DateRangePicker`, `FolderNameInput`)**: PR #307 (04 create-album-modal) places them
  inside its own folder. For stories 06 (rename) and 07 (edit-dates), the subagent copies them locally into
  its own modal folder. The duplication is intentional: it keeps the PRs independent. The follow-up
  wiring story (01) is free to extract them into `web-nextjs/components/` as a cleanup.
- **The viewer stories 02 and 03 are copied verbatim** from the archived feature (issues 09 and 10). The
  internal links to ADRs 0001 and 0003 and the ADR paths are preserved (just one directory level deeper).

## Stories

### 01 — Wire album-header modals into the album page

Done **last** among parity work. Enables the Share / Rename / Dates buttons in `AlbumHeader` (currently
disabled), adds a Delete button next to them, mirrors them in `AlbumActionsFab`, and mounts the four modal
components delivered by stories 05–08 inside `AlbumPageContent`, each wired to its selector and to the
matching thunks. All four buttons are owner-only (hidden for non-owners, not just disabled).

### 02 — Fullscreen media viewer (YARL + route-driven)

Carried verbatim from archived `09-fullscreen-viewer-yarl.md`. YARL as the viewer library, NextJS route as
the source of truth (`photos/[mediaId]`), modal interception per ADR-0001, image rendering via `next/image`
+ custom loader.

### 03 — Viewer capabilities (zoom, download, counter, video)

Carried verbatim from archived `10-viewer-capabilities.md`. YARL bundled plugins; configuration only.

### 04 — Create album modal

Owner-only button on the home page opens a dialog with name (required), start/end date pickers, optional
custom folder name. Wires `createAlbumThunk`. Validation inline. Success closes the dialog and refreshes
the album list. No visual photo preview — deferred.

### 05 — Share album modal (component + stories)

Deliver a `ShareAlbumModal` React component with Storybook stories. Props shape: `SharingDialogFrag &
ShareAlbumModalHandlers`. Behaviour and layout mirror `web/src/components/albums/ShareDialog/index.tsx`.
**Not** mounted into the album page by this story.

### 06 — Rename album modal (component + stories)

Deliver a `RenameAlbumModal` React component with Storybook stories. Props shape: `EditNameDialogSelection &
EditNameDialogHandlers`. Behaviour and layout mirror `web/src/components/albums/EditNameDialog/`.
**Not** mounted into the album page by this story.

### 07 — Edit album dates modal (component + stories)

Deliver an `EditDatesModal` React component with Storybook stories. Props shape: `EditDatesDialogSelection &
EditDatesDialogHandlers`. Behaviour and layout mirror `web/src/components/albums/EditDatesDialog/`.
**Not** mounted into the album page by this story.

### 08 — Delete album modal (component + stories)

Deliver a `DeleteAlbumModal` React component with Storybook stories. Props shape: `DeleteDialogFrag &
DeleteAlbumModalHandlers`. Behaviour and layout mirror `web/src/components/albums/DeleteAlbumDialog/`.
**Not** mounted into the album page by this story.

### 09 — Client-side access-token refresh

Discovered during review of story 04: client adapters rely on the httpOnly cookie being auto-attached and
have no refresh path, so a tab left open past the access-token TTL fails the next write with a 401. Adds a
`POST /auth/refresh` route handler reusing `refreshSession`, and a shared `fetch` wrapper that
retries-on-401 once. Soft-blocks 05–08 so write-stories don't inherit the stale-tab assumption.

### 10 — Cutover: delete `web/` and drop the basepath

Mechanical closure of the feature. Only picked up when 01–08 are all `done`.

- Delete `web/`, update `AGENTS.md`, `README.md`, `Makefile`, CI workflows.
- Remove `basePath: '/nextjs'` from `web-nextjs/next.config.ts` and any URL-building code that assumes it.
- Update `deployments/cdk/` to serve `web-nextjs` at `/`; drop the `web/` lambda.
- Flag any required human action (e.g. updating Google OAuth redirect URIs in the IDP console) in the PR
  description so Thomas can perform it before deploy.
