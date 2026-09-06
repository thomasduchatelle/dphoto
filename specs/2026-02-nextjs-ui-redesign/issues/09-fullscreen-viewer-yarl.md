# Full-screen media viewer with route-driven modal interception

**Status**: ready

## Description

Let a user open any media from the album grid into a full-screen viewer that overlays the grid, and navigate
between medias with keyboard and touch. The viewer is the `yet-another-react-lightbox` (YARL) library, but the
NextJS route — not the library — is the source of truth for which media is shown and whether the viewer is
open.

Clicking a thumbnail navigates to `/albums/[ownerId]/[albumId]/photos/[mediaId]`, which renders the viewer as
a modal over the still-mounted grid. The URL stays shareable and refresh-safe.

## Context

- Read **[ADR-0003](../../../docs/adr/0003-photo-viewer-lightbox-driven-by-route.md)** first: it mandates YARL
  for the viewer and that the route drives YARL's `open`/`index` (the `mediaId` in the URL selects the media),
  with close = `router.back()`.
- Read **[ADR-0001](../../../docs/adr/0001-owner-based-routing-and-modal-interception.md)** for the routing
  mechanics: a parallel + intercepting route `@modal/(.)photos/[photoId]/` renders the modal over the grid
  when navigated from within the album, and a fallback `photos/[photoId]/page.tsx` renders the same viewer as
  a full page for direct access / refresh / shared links. Keep the viewer a single shared component mounted at
  both places so they cannot drift.
- Install `yet-another-react-lightbox` and import its base styles (`yet-another-react-lightbox/styles.css`).
- The displayed image must render through `next/image` + the custom loader (`libs/image-loader.ts`) using
  YARL's slide render hook, requesting `1440` on small screens and `2400` on larger ones. Do not let YARL load
  raw image URLs directly.
- Medias come from the catalog state (`Media` in `domains/catalog/language/catalog-state.ts`); the album's
  medias, in display order, define the viewer's slide list and each media's index.

## Acceptance Criteria

- Clicking a thumbnail in the grid navigates to `/albums/[ownerId]/[albumId]/photos/[mediaId]` and opens the
  viewer as a modal over the grid; the grid remains mounted behind it.
- YARL's `open` and `index` are derived from the route: the `mediaId` in the URL selects the displayed media.
  The viewer holds no open/index state that can diverge from the URL.
- Closing the viewer (close button or ESC) performs `router.back()` and returns to the grid.
- LEFT / RIGHT arrow keys and touch swipe move to the previous / next media.
- Navigating inside the viewer keeps the URL in sync with the currently displayed media.
- The displayed image uses `next/image` + the custom loader at a screen-appropriate width (`1440` small,
  `2400` larger).
- Refreshing, or opening a `photos/[mediaId]` URL directly, renders the same viewer as a full page.
- A `mediaId` that does not exist in the album shows the not-found page.

## Out of scope

- Zoom, download, position counter, and video playback — next story (viewer capabilities).
- Restoring the grid scroll position on close — later story.
- The shared-element "zoom from / back to the thumbnail" animation — deferred, not in this epic.
