# ADR-0001: Owner-based album routing with parallel-route modal interception

- Status: accepted
- Date: 2026-01-31
- Scope: `web-nextjs` (NextJS App Router)

## Context

Albums are addressed by both an owner and an album identifier: two different owners can share albums into the same viewer, so an album is not uniquely identified by its own id alone. The UI also needs to open a photo in a full-screen viewer *over* the album grid, while keeping the photo URL shareable and refresh-safe.

## Decision

- Album URLs carry the owner: `/albums/[ownerId]/[albumId]`.
- `/albums` and `/albums/[ownerId]` redirect to `/` (they are not meaningful landing pages).
- The full-screen photo viewer uses a NextJS **parallel + intercepting route**: `@modal/(.)photos/[photoId]/` renders the viewer as a modal over the grid when navigated from within the album.
- A **fallback** route `photos/[photoId]/page.tsx` renders the same viewer as a full page for direct access / refresh / shared links.
- The album layout renders both `{children}` and `{modal}`.

User flow:

1. `/` album list → click album → `/albums/[ownerId]/[albumId]` (grid).
2. Click photo → `/albums/[ownerId]/[albumId]/photos/[photoId]` (modal over grid).
3. ESC / close → back to the grid.
4. Refresh on a photo URL → full-page viewer.

## Consequences

- URLs are shareable and unambiguous across owners.
- Modal behaviour is native to the framework; no bespoke modal/state machinery.
- The viewer must exist in two mount points (intercepted modal + fallback page); keep the viewer a single shared component to avoid drift.
