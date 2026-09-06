# Epic 2 — Photo Viewing & Navigation (refined stories)

Draft of the Epic 2 issues. Supersedes the Epic 2 section of `epics.md`. Once agreed, each story below
becomes `issues/NN-<slug>.md` (numbered from `07`, continuing after Epic 1).

## Why this differs from `epics.md`

- **Dropped** old 2.1 (URL redirects) and 2.5 (loading indicator): already implemented (`not-found.tsx`,
  `NavigationLoadingIndicator`).
- **Dropped** old 2.2 (design PhotoViewer component upfront): we do not design the viewer before the page
  that hosts it, and the viewer is now a library, not a bespoke component.
- **Sequenced as a tracer bullet**: wire the data and render a dumb list first (2.1), then make it a proper
  grid (2.2), mirroring Epic 1's 1.3 → 1.4.

## Design decisions (to record as ADR-0003)

- **Grid: no gallery library.** The day-grouped grid is built with Material UI `sx` + `next/image` and the
  existing custom loader. State already exposes `MediaWithinADay[]`; a gallery lib would fight the grouping.
- **Viewer: adopt `yet-another-react-lightbox` (YARL)** — MIT, React 19-ready, maintained. Zoom, download,
  counter and video are bundled plugins (configuration, not stories).
- **NextJS routing is the source of truth**, not the library. The `photos/[photoId]` intercepting route
  (ADR-0001) drives YARL's `open`/`index`; closing the viewer is `router.back()`. This keeps photo URLs
  shareable and refresh-safe while YARL only owns the visual/gesture layer.
- **Deferred (fast-follow, not in this epic):** the shared-element "zoom from / back to the thumbnail"
  animation. YARL does not provide it; revisiting it later must not change the URL contract above.

References for all stories: `spec.md`, `architecture.md`,
[ADR-0001](../../docs/adr/0001-owner-based-routing-and-modal-interception.md) (routing + modal interception),
[ADR-0002](../../docs/adr/0002-server-computed-initial-state-hydration.md) (server-computed state → pure UI).

---

## Story 2.1: Load album medias and render as a simple list

**Target file:** `issues/07-load-album-medias.md`

As a user,
I want the album page to load the photos and videos of the album I opened,
So that I can confirm the album's content is there before it's styled.

**Acceptance Criteria:**

**Given** I am an authenticated user with access to an album
**When** I navigate to `/albums/[ownerId]/[albumId]`
**Then** the album detail and its medias are loaded via the migrated catalog state (server-computed initial
state per ADR-0002)
**And** the medias are rendered as a plain, unstyled list in the album page component (e.g. filename + capture
time per item) (FR12)
**And** a back link/button returns to the home page `/` (FR17)
**And** when the album has no medias, the existing empty state (`NoMedia`) is shown (FR38)
**And** when the album does not exist or is not accessible, the not-found page is shown (FR45)
**And** when loading fails, an error is shown with a recovery action (FR39, NFR3)

**Out of scope:**

- No grid layout, no day grouping headers — that's Story 2.2
- No image thumbnails / `next/image` — that's Story 2.2
- No full-screen viewer — that's Story 2.3

---

## Story 2.2: Display album medias in a day-grouped grid

**Target file:** `issues/08-media-grid-grouped-by-day.md`

As a user,
I want to view an album's photos in a grid grouped by the day they were captured,
So that I can browse the album chronologically and see the story of that period.

**Acceptance Criteria:**

**Given** I am viewing an album with medias at `/albums/[ownerId]/[albumId]`
**When** the page displays
**Then** medias are shown in a responsive grid: mobile (2 columns), tablet (3), desktop (4), large (5) (NFR8)
**And** medias are grouped by capture day with a date header per group, in the user's locale format
(FR12, FR18)
**And** each thumbnail is rendered with `next/image` + the custom loader at the grid-appropriate width (360),
size-optimised (FR14, NFR2, NFR6)
**And** video medias are visually distinguishable from photos (e.g. a play indicator)
**And** the album name and metadata are shown in the page header, with a back link to `/` (FR17)
**And** the empty album state (`NoMedia`) is preserved
**And** the grid and header have Storybook visual regression tests

**Out of scope:**

- Opening a media full-screen — that's Story 2.3
- Zoom, download, counter, video playback controls — that's Story 2.4
- Restoring scroll position after closing the viewer — that's Story 2.5

---

## Story 2.3: Full-screen viewer via YARL with URL-driven modal interception

**Target file:** `issues/09-fullscreen-viewer-yarl.md`

As a user,
I want to open a media full-screen over the grid and navigate between medias,
So that I can view photos clearly with a focused, shareable experience.

**Acceptance Criteria:**

**Given** I am viewing the album grid
**When** I click a thumbnail
**Then** the URL updates to `/albums/[ownerId]/[albumId]/photos/[mediaId]` and the viewer opens as a modal
over the grid, which stays mounted behind it (ADR-0001) (FR13)
**And** the viewer is `yet-another-react-lightbox`, with its `open`/`index` derived from the route (the
`mediaId` in the URL selects the displayed media)
**And** closing the viewer (X or ESC) performs `router.back()` and returns to the grid
**And** I can move to the next/previous media with LEFT/RIGHT arrow keys and with touch swipe (FR14, FR15,
NFR6, NFR7, NFR9)
**And** navigating within the viewer keeps the URL in sync with the displayed media
**And** the displayed media loads via `next/image` + the custom loader at a screen-appropriate width
(1440 on small screens, 2400 above) (NFR2)
**And** refreshing or opening a `photos/[mediaId]` URL directly renders the same viewer as a full page
(fallback route per ADR-0001)
**And** a non-existent media shows the not-found page (FR45)

**Out of scope:**

- Zoom, download, counter and video plugins — that's Story 2.4
- Restoring the grid scroll position on close — that's Story 2.5
- The shared-element zoom-from/to-thumbnail animation — deferred (fast-follow)

---

## Story 2.4: Viewer capabilities — zoom, download, counter, video

**Target file:** `issues/10-viewer-capabilities.md`

As a user,
I want to zoom, download, and see my position while viewing, and play videos,
So that the full-screen viewer is fully usable for every media type.

**Acceptance Criteria:**

**Given** I have a media open in the viewer
**When** I interact with it
**Then** I can zoom into photo detail via pinch, double-tap, and scroll/wheel (FR16)
**And** a download control saves the original media file to my device
**And** a counter shows my position in the album (e.g. "5 of 47")
**And** video medias play within the viewer using the video plugin
**And** all controls are keyboard-accessible and do not conflict with browser defaults (NFR6)

**Out of scope:**

- Restoring the grid scroll position on close — that's Story 2.5
- The shared-element zoom-from/to-thumbnail animation — deferred (fast-follow)

---

## Story 2.5: Restore scroll to the last viewed media on close

**Target file:** `issues/11-restore-scroll-on-close.md`

As a user,
I want the grid to bring the last media I viewed into view when I close the viewer,
So that I can see where I was and continue browsing from there.

**Acceptance Criteria:**

**Given** I opened a media, navigated to others within the viewer, then closed it
**When** I return to the grid
**Then** the last-viewed media's thumbnail is scrolled into view
**And** the last-viewed media is determined from the route (the `mediaId` last shown), not bespoke state
**And** this works whether I closed via X, ESC, or browser back
**And** if the last-viewed media is already visible, the grid does not jump

**Out of scope:**

- The shared-element zoom-from/to-thumbnail animation — deferred (fast-follow)
</content>
</invoke>
