# Display album medias in a day-grouped grid

**Status**: ready

## Description

Replace the plain media list on the album page (`/albums/[ownerId]/[albumId]`) with a responsive grid of
thumbnails, grouped by the day each media was captured, with a date header per group. This is the "proper
rendering" step after the album medias are loaded and listed.

Medias arrive from the catalog state already grouped: `MediaWithinADay[]` where each entry is
`{ day: Date, medias: Media[] }` (see `domains/catalog/language/catalog-state.ts` and
`domains/catalog/navigation/group-by-day.ts`). Each `Media` has an `id`, a `type`
(`MediaType.IMAGE | VIDEO | OTHER`), a `time`, and a `contentPath` used to build its image URL.

## Context

- Thumbnails must use `next/image` with the project's custom image loader (`libs/image-loader.ts`, configured
  in `next.config.ts`), requesting the small width (`360`) for grid display. Do not bypass the loader or do
  client-side image processing.
- Styling is Material UI only, via the `sx` prop with theme breakpoints (`xs`/`sm`/`md`/`lg`). No Tailwind,
  no inline styles. The dark theme and brand color are already configured in `components/theme/`.
- Components must be pure (state + handlers in, no internal data fetching), consistent with
  **[ADR-0002](../../../docs/adr/0002-server-computed-initial-state-hydration.md)**.
- The empty-album component `NoMedia` already exists and must be preserved for albums with no medias.
- Page-specific components live in an `_components/` folder next to the page.

## Acceptance Criteria

- Medias are shown in a responsive grid: 2 columns on mobile (`xs`), 3 on tablet (`sm`), 4 on desktop (`md`),
  5 on large screens (`lg`).
- Medias are grouped by capture day, with a date header above each group, formatted in the user's locale.
- Each thumbnail is a `next/image` using the custom loader at width `360`, sized to fit the grid cell.
- Video medias (`MediaType.VIDEO`) are visually distinguishable from photos (e.g. a play-icon overlay).
- The album name and metadata are shown in a page header, with a back link to the home page `/`.
- Albums with no medias still show the `NoMedia` empty state.
- The grid and the page header have Storybook stories with visual regression coverage.

## Out of scope

- Opening a media full-screen — next story (full-screen viewer).
- Zoom, download, position counter, video playback controls — later story (viewer capabilities).
- Restoring the grid scroll position after the viewer closes — later story.
