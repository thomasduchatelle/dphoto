# Viewer slideshow: Play starts an auto-advancing full-screen show

**Status**: done
**Depends on**: `02-fullscreen-viewer-yarl.md` (the route-driven YARL viewer), `01-header-action-scaffold.md` (the Play button wiring).

> Post-cutover enhancement. **Not** a launch criterion for the NextJS cutover — the feature ships whether or
> not this is done. Grouped here for convenience.

## Description

Today the **Play** button (desktop `AlbumHeader`, mobile `AlbumActionsFab`) opens the first media in the
full-screen viewer, but the user then has to advance manually. Make Play start an **automatic slideshow**:
the viewer opens on the first media and advances to the next on a timer until the user stops it or reaches the
end of the album.

## Context

- Read **[ADR-0003](../../../docs/adr/0003-photo-viewer-lightbox-driven-by-route.md)**: the viewer is
  `yet-another-react-lightbox` (YARL) and the NextJS route is the source of truth for open/index. YARL ships a
  **Slideshow** bundled plugin — prefer it over bespoke timer code, consistent with issue 03's plugin approach.
- The slideshow must keep the URL contract: as each slide auto-advances, the displayed `mediaId` stays in sync
  with the URL via the existing `onNavigate` (History API update), exactly as manual navigation does. Adding
  the slideshow must not introduce viewer-owned index state that can diverge from the route.
- The Play entry points already exist and route to the first media (`photoUrl(albumId, firstMediaId)`). This
  story adds the "and start playing" intent and the auto-advance behaviour; it must not regress the plain
  open-first-photo behaviour for a direct `photos/[mediaId]` visit.
- Video medias (`MediaType.VIDEO`) interact with auto-advance — decide and document whether the slideshow
  waits for playback to end or skips/advances on the timer. Keep it simple for the MVP.

## Acceptance Criteria

- Clicking **Play** opens the full-screen viewer on the first media and begins auto-advancing through the
  album in display order.
- The viewer exposes a visible control to pause/resume the slideshow (YARL Slideshow plugin control).
- As the slideshow advances, the URL stays in sync with the current media (shareable / refresh-safe), with no
  full-page reload or modal remount.
- Reaching the last media stops the slideshow (finite carousel; no wrap-around unless trivially free).
- Closing the viewer (ESC / close) stops the slideshow and returns to the grid as today.

## Out of scope

- Configurable slideshow interval / per-user preference — a sensible fixed interval is enough for the MVP.
- Transition/animation styling beyond YARL's default.
- Background music or Ken Burns effects.
