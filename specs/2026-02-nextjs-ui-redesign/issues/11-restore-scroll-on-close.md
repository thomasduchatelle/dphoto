# Restore scroll to the last viewed media when the viewer closes

**Status**: ready

## Description

When a user opens a media, navigates through several inside the full-screen viewer, then closes it, the album
grid should bring the last-viewed media's thumbnail into view. This lets the user see where they ended up and
keep browsing from there, instead of landing back at their original scroll position (or the top).

## Context

- The full-screen viewer is route-driven per
  **[ADR-0003](../../../docs/adr/0003-photo-viewer-lightbox-driven-by-route.md)** and
  **[ADR-0001](../../../docs/adr/0001-owner-based-routing-and-modal-interception.md)**: the currently displayed
  media is the `mediaId` in the URL (`/albums/[ownerId]/[albumId]/photos/[mediaId]`), and closing is
  `router.back()`. Determine the last-viewed media from the route — do not add separate bespoke state to track
  it.
- The grid is the day-grouped media grid on the album page; each media has a stable `id`
  (`domains/catalog/language/catalog-state.ts`).

## Acceptance Criteria

- After opening a media, navigating to others in the viewer, and closing it, the grid scrolls the last-viewed
  media's thumbnail into view.
- The last-viewed media is derived from the route (the last `mediaId` shown), not from extra state added for
  this purpose.
- This works whether the viewer was closed via the close button, ESC, or browser back.
- If the last-viewed media's thumbnail is already visible in the viewport, the grid does not jump or re-scroll.

## Out of scope

- The shared-element "zoom from / back to the thumbnail" animation — deferred, not in this epic.
