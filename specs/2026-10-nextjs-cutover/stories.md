# NextJS Cutover — stories

Draft of the ten issues. Each becomes `issues/NN-<slug>.md`. Parallelism target: after issue 01 lands,
issues 02–09 can all be picked up in parallel; issue 10 is the final merge/deploy.

## Parallelism map

```
01 header-action-scaffold ──┬── 02 fullscreen-viewer-yarl ──── 03 viewer-capabilities
                            │
                            ├── 05 modal-share-album
                            ├── 06 modal-rename-album
                            ├── 07 modal-edit-dates
                            └── 08 modal-delete-album

(independent)   04 create-album-modal
(independent)   09 client-token-refresh    (soft-blocks 05–08; discovered during 04)

(after all)     10 cutover
```

Issues 02 and 03 are sequential (03 depends on 02). All five modals (04–08) are independent of each other and
of the viewer, provided 01 has landed and placed the stub modal components. 10 is the final cutover.

## Design decisions

- **No shared "modal framework"**. Each modal is its own component, owning its own open/close state via the
  existing thunks (`openXxxModal` / `closeXxxModal`). This avoids a cross-story coordination hazard.
- **The owner-action buttons live in `AlbumHeader`** (desktop) **and `AlbumActionsFab`** (mobile). Issue 01
  enables them and wires each to a `useXxxModal` hook + stub component. The modal stories replace the stub.
- **Viewer stories 02 and 03 are copied verbatim** from the archived feature (issues 09 and 10). The
  internal links to ADRs 0001 and 0003 and the ADR paths are preserved (just one directory level deeper).

## Stories

### 01 — Header owner-action scaffolding

Enable the Share / Rename / Dates buttons in `AlbumHeader`, add a Delete button, and create stub modal
components so the modal stories can land in parallel without conflicts. Visibility is driven by album
ownership (owner-only). Buttons dispatch the corresponding `openXxxModal` thunk; the stub component renders
"Not implemented yet" and a Close button wired to `closeXxxModal`. Same wiring for `AlbumActionsFab` on
mobile.

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

### 05 — Share album modal

Owner-only, from `AlbumHeader`. Lists users with access (avatars + names + emails), grant access by email
with validation, revoke with confirmation. Wires existing `openSharingModal` / `grantAlbumAccessThunk` /
`revokeAlbumAccessThunk` / `closeSharingModal` thunks and state in `domains/catalog/sharing/`.

### 06 — Rename album modal

Owner-only, from `AlbumHeader`. Dialog with name field pre-filled. Wires `album-edit-name` thunks. Updated
name reflected immediately in header and album list.

### 07 — Edit album dates modal

Owner-only, from `AlbumHeader`. Dialog with start/end date pickers pre-filled. Backend-driven orphan warning
surfaces as a dialog message. Wires `album-edit-dates` thunks. On success, album page refreshes to show the
updated media set. No visual photo preview — deferred.

### 08 — Delete album modal

Owner-only, from `AlbumHeader`. Confirmation dialog with destructive styling, Cancel as default focus. Wires
`album-delete` thunks. On success, user is navigated back to the home page and the album disappears from the
list.

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
