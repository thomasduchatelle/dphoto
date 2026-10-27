# Neighbour-album rail: limit to 7 albums centred on the displayed one

**Status**: ready
**Depends on**: none — a self-contained change to the rail selection + `AlbumRail`; can start on day 1.

## Description

Today the album rail on the album page (desktop, left side) renders the **entire** filtered list of albums.
Limit it to at most **7 albums**, as a window centred on the currently displayed album: **3 before and 3
after**. At the edges of the list, bias the window so 7 albums stay visible when possible:

- displayed album is **first** → show it + **6 after**;
- displayed album is **last** → **6 before** + it;
- otherwise → **3 before** + displayed + **3 after**;
- when the full list has 7 or fewer albums, show them all.

The currently-displayed album keeps its highlight and must always be within the window.

## Context

- The rail renders in
  `web-nextjs/app/(authenticated)/albums/[ownerId]/[albumId]/_components/AlbumRail/index.tsx`, from the
  `albums` list. It currently does `albums.map(...)` over the full list with no windowing.
- The list order and the "which album is displayed" index are already known to the navigation layer:
  `web-nextjs/domains/catalog/navigation/selector-catalog-viewer-page.ts` already computes `displayedAlbum`
  and the `previousAlbum` / `nextAlbum` neighbours from `state.allAlbums` by index. Prefer computing the
  7-album window **in the selector** (or a small pure helper it calls) and exposing it to the rail, rather
  than slicing inside the React component — keep the component dumb and the windowing logic unit-/story-
  testable. Decide whether the window is derived from `state.albums` (the filtered list the rail shows) or
  `state.allAlbums`; use the same list the rail currently renders so filtering still behaves, and centre on
  the displayed album's index within that list.
- Neither `web/` nor `web-nextjs` had this limit before — this is new behaviour, so there is no reference
  implementation to mirror. Keep it a straightforward pure windowing function.
- Follow the `nextjs` and `ui-components` skills. If the windowing lives in the selector, a focused unit test
  of the pure helper is appropriate (edge cases below). The `AlbumRail` component change is covered by
  Storybook.

## Acceptance Criteria

- When more than 7 albums exist, the rail shows at most 7: the displayed album centred with 3 before and 3
  after.
- When the displayed album is the **first** in the list, the rail shows it + the 6 following albums.
- When the displayed album is the **last**, the rail shows the 6 preceding albums + it.
- When the list has 7 or fewer albums, all are shown (no empty slots, no error).
- The displayed album is always present in the window and keeps its highlight.
- The windowing logic is a pure function covered by unit tests for the edge cases (first, last, middle, fewer
  than 7, exactly 7).

## Out of scope

- Any "show more / expand" affordance to reveal the full list — the window is fixed at 7.
- Changing the mobile neighbour banners (`PreviousAlbumBanner`, `NeighbourAlbumLink`) — they use
  `previousAlbum` / `nextAlbum` and are unaffected.
- Changing the album filtering behaviour.
