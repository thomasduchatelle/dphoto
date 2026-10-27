# NextJS Cutover — complete feature parity and retire `web/`

**Author:** Arch
**Date:** 2026-10-11
**Supersedes:** `specs/archived/2026-02-nextjs-ui-redesign/` (issues 01–08 shipped, 09–11 moved or deferred).

## Intent

Finish the migration from `web/` (deprecated Waku app) to `web-nextjs` (NextJS + MUI) by reaching **feature
parity** and performing the cutover: delete `web/`, drop the `/nextjs` basepath from `web-nextjs`, and update
the deployment so the NextJS app serves at the root path.

This is a **parity feature**, not a redesign. Everything beyond parity (random photos, timeline navigation,
visual date-picker, progressive image loading, hover previews) is explicitly out of scope and tracked in
`TODO.md` for separate features once the cutover is complete.

## What's already done

Issues 01–08 of the archived redesign feature already landed:

- MUI dark theme, state migrated from `web/src/core/catalog/`, custom image loader
- Album list page (home) with cards, density, filtering, and sharing status
- Album page with day-grouped media grid, back navigation, neighbour-album rail
- Sharing thunks and state (`domains/catalog/sharing/`) already migrated from `web/`; also album-create,
  album-edit-dates, album-edit-name, album-delete thunks

The missing surface is the **full-screen media viewer** and the **modals** that let an owner manage an album.
The domain layer for those modals already exists; this feature is UI wiring plus the viewer.

## In scope

Six capabilities, each delivered as a modal / viewer wired to existing thunks:

1. **Full-screen media viewer** — click a thumbnail to open the media over the grid with keyboard + touch
   navigation, zoom, download, counter, video playback. Covered by the three viewer stories carried from the
   archived feature (see `docs/adr/0001`, `docs/adr/0003`).
2. **Create album** — on the album list page: button opens a dialog with name + start/end dates + optional
   custom folder name. Basic version, no visual photo preview.
3. **Share album** — on the album page: owner-only, lists users with access, grant by email, revoke with
   confirmation. The sharing thunks (`domains/catalog/sharing/`) are already migrated.
4. **Rename album** — on the album page: owner-only, dialog with pre-filled name field.
5. **Edit dates** — on the album page: owner-only, dialog with pre-filled start/end pickers, orphan-media
   warning. Basic version, no visual photo preview.
6. **Delete album** — on the album page: owner-only, confirmation dialog with destructive styling.

Plus the mechanical pieces:

- **Album header action scaffolding** — enable the (currently disabled) owner-action buttons in
  `AlbumHeader`, add a Delete button, and wire each button to a stub modal component so modal stories can land
  in parallel without conflicts.
- **Media thumbnail overlay** — reach parity with `web/`'s hover overlay by showing the camera and the time
  on each grid thumbnail on hover. Unlike `web/`, the overlay does **not** show the date (year/month/day) — it
  is redundant with the per-day grouping header. Camera + time only.
- **Neighbour-album rail windowing** — limit the album rail on the album page to at most 7 albums centred on
  the displayed one (3 before + 3 after), biasing the window to keep 7 visible at the list edges (6 after when
  the displayed album is first, 6 before when it is last).
- **Cutover** — delete `web/`, remove `basePath: '/nextjs'` from `web-nextjs/next.config.ts`, update CDK /
  deployment config to route the root path to `web-nextjs`, update docs.

## Out of scope (deferred, see `TODO.md`)

- Random photo discovery on the home page
- Timeline navigation between neighbouring albums (beyond the existing rail)
- Visual date-picker with photo preview on create / edit-dates (needs a backend endpoint to find medias by
  date)
- Progressive image loading (blur-up)
- Album preview on hover
- Scroll-restore to last-viewed media when the viewer closes (nice-to-have; the grid stays mounted behind the
  modal so the user returns to their original scroll position regardless)
- Any backend change — the APIs consumed are what `web/` already consumes

## Launch criteria

- All ten issues of this feature are `done`
- Thomas can create, share, rename, re-date, and delete albums using `web-nextjs` as his primary interface
- `web/` has been deleted from the repository
- `web-nextjs` serves at the root path (no `/nextjs` basepath)
- Deployment is green

## Entry points (UI contract)

To let the modal stories land in parallel, the entry points for each action are fixed upfront:

| Action      | Entry point                                                                     | Visibility  |
|-------------|---------------------------------------------------------------------------------|-------------|
| Create      | Button on the album list page (home)                                            | Owners only |
| Share       | Button in `AlbumHeader` on the album page (desktop) + `AlbumActionsFab` (mobile)| Owners only |
| Rename      | Button in `AlbumHeader` on the album page (desktop) + `AlbumActionsFab` (mobile)| Owners only |
| Edit dates  | Button in `AlbumHeader` on the album page (desktop) + `AlbumActionsFab` (mobile)| Owners only |
| Delete      | Button in `AlbumHeader` on the album page (desktop) + `AlbumActionsFab` (mobile)| Owners only |

The disabled Share / Rename / Dates buttons already exist in `AlbumHeader`; Delete is added, and all four are
wired by the first issue of this feature. The create button on the home page is independent and already
partially present (`AlbumGrid` exposes an `onCreateAlbum` prop).

## Cutover act (final issue)

Mechanical, done last once all parity work is `done`:

- Delete the `web/` directory and all references to it (`AGENTS.md`, `README.md`, `Makefile` targets, CI).
- Remove `basePath: '/nextjs'` from `web-nextjs/next.config.ts` and any derived URL assumptions.
- Update `deployments/cdk/` so `web-nextjs` serves at `/`; drop the `web/` lambda / CloudFront behaviour.
- Verify Google OAuth redirect URIs still match (any `/nextjs` suffix in the configured callback must be
  removed in the IDP console; a wizard step may be needed if a human action is required).
- Run full test suites, synth the CDK template, deploy.
