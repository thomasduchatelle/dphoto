# ADR-0003: Photo viewer uses a lightbox library driven by the route

- Status: accepted
- Date: 2026-09-06
- Scope: `web-nextjs` (NextJS App Router)

## Context

The album page must show a media grid and let a user open any media full-screen with zoom, download,
next/previous navigation (keyboard, touch), a position counter, and video playback. Building all of this
by hand is costly and error-prone.

Two constraints shape the choice:

- The day-grouped grid is already backed by state (`MediaWithinADay[]`) and is rendered with Material UI
  and `next/image` + the custom image loader. A gallery library would compete with that grouping and loader.
- ADR-0001 makes the full-screen viewer a NextJS parallel + intercepting route
  (`photos/[photoId]`), so photo URLs are shareable and refresh-safe. Most lightbox libraries own their own
  open/close state and are not URL-driven — a direct conflict.

## Decision

- **The grid uses no gallery library.** It is built with Material UI `sx` and `next/image`.
- **The full-screen viewer is `yet-another-react-lightbox` (YARL)** — MIT-licensed, React 19-ready,
  maintained. Zoom, download, counter, and video are enabled as its bundled plugins (configuration, not
  bespoke code).
- **The NextJS route is the source of truth, not the library.** The `photos/[photoId]` route (ADR-0001)
  drives YARL's `open`/`index`: the `mediaId` in the URL selects the displayed media. Closing the viewer is
  `router.back()`. Navigating inside the viewer keeps the URL in sync. This preserves the shareable,
  refresh-safe URL contract while YARL owns only the visual and gesture layer.
- **The displayed image** is rendered through `next/image` + the custom loader via YARL's slide render hook,
  so the viewer reuses the same width mapping as the rest of the app.

## Consequences

- Zoom / download / counter / video become plugin configuration rather than stories.
- The viewer stays a single shared component mounted at both the intercepted-modal and fallback-page routes
  (ADR-0001); it must not hold its own open/index state that could drift from the URL.
- The shared-element "zoom from / back to the thumbnail" animation is **out of scope**: YARL does not provide
  it. If added later, it must not change the URL contract above.
- If YARL is ever replaced, the route contract is the stable boundary; the swap is confined to the view layer.
