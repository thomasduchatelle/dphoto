# Load album medias and render as a simple list

**Status**: ready

## Description

The album page at `/albums/[ownerId]/[albumId]` currently renders only an empty-album placeholder
(`app/(authenticated)/albums/[ownerId]/[albumId]/_components/AlbumPageContent/index.tsx` returns `<NoMedia/>`).

Make it load the album and its medias and render them as a plain, unstyled list. This is a tracer bullet: it
proves the data flow (server load → hydrated state → rendered medias) end to end. A later story turns the list
into a styled, day-grouped grid — so keep the rendering deliberately minimal here.

Medias are already available in the catalog state as `MediaWithinADay[]` (see
`domains/catalog/language/catalog-state.ts`: `Media`, `MediaWithinADay`). Album loading uses the migrated
catalog state and the server-computed-initial-state pattern.

## Context

- Follow **[ADR-0002](../../../docs/adr/0002-server-computed-initial-state-hydration.md)**: the Server
  Component computes the initial state (loading the album + medias through the catalog thunks), a Client
  Component hydrates it, and pure UI components receive state and handlers as props. Do not load with
  `useEffect`.
- Follow **[ADR-0001](../../../docs/adr/0001-owner-based-routing-and-modal-interception.md)** for routing:
  the album URL is `/albums/[ownerId]/[albumId]`; a not-found page is shown when the album is missing.
- Album loading thunks and selectors already exist under `domains/catalog/navigation/` (e.g.
  `thunk-onPageRefresh.ts`, `selector-catalog-viewer-page.ts`, `group-by-day.ts`). Reuse them; do not
  reimplement state management.
- The empty-album component `NoMedia` already exists next to the page.

## Acceptance Criteria

- Navigating to `/albums/[ownerId]/[albumId]` loads that album's detail and medias via the catalog state,
  server-side (no client-side `useEffect` for the initial load).
- The medias are rendered as a plain, unstyled list in the album page component — one line per media showing
  at least its filename and capture time. No images, no grid, no styling.
- A back link/button returns to the home page `/`.
- When the album has no medias, the existing `NoMedia` component is shown.
- When the album does not exist or the user has no access, the not-found page is shown.
- When loading fails, an error is displayed with a recovery action (reuse the existing `ErrorMessage`
  component / error boundary pattern).

## Out of scope

- Grid layout and per-day date headers — next story (day-grouped media grid).
- Thumbnails / `next/image` / the custom image loader — next story.
- Opening a media full-screen — later story (full-screen viewer).
