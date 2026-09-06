# Viewer capabilities: zoom, download, position counter, video

**Status**: ready

## Description

Enable the remaining full-screen viewer features on top of the working `yet-another-react-lightbox` (YARL)
viewer: content zoom, downloading the original file, a position counter, and video playback. These are all
YARL bundled plugins, so this story is mostly configuration and wiring — no bespoke viewer code.

## Context

- Read **[ADR-0003](../../../docs/adr/0003-photo-viewer-lightbox-driven-by-route.md)**: these features are
  delivered via YARL's bundled plugins (Zoom, Download, Counter, Video), not custom implementations. The route
  remains the source of truth for open/index; adding plugins must not introduce viewer-owned navigation state.
- The viewer is a single shared component used by both the intercepted-modal route and the fallback full-page
  route (see [ADR-0001](../../../docs/adr/0001-owner-based-routing-and-modal-interception.md)); configure the
  plugins once, there.
- Each `Media` (`domains/catalog/language/catalog-state.ts`) has a `type`
  (`MediaType.IMAGE | VIDEO | OTHER`) and a `contentPath`. Video slides are driven by `type === VIDEO`. The
  download must fetch the original media file, not a resized thumbnail.

## Acceptance Criteria

- Photos can be zoomed via pinch, double-tap, and scroll/wheel (YARL Zoom plugin).
- A download control saves the original media file to the user's device.
- A counter shows the current position within the album (e.g. "5 of 47").
- Video medias (`MediaType.VIDEO`) play within the viewer (YARL Video plugin).
- All controls are reachable by keyboard and do not conflict with browser default shortcuts.

## Out of scope

- Restoring the grid scroll position on close — next story.
- The shared-element "zoom from / back to the thumbnail" animation — deferred, not in this epic.
