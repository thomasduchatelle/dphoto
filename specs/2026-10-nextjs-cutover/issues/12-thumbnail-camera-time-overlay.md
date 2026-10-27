# Media thumbnail overlay: camera and time on hover

**Status**: ready
**Depends on**: none — a self-contained change to `MediaThumbnail`, Storybook-driven; can start on day 1.

## Description

Reach parity with `web/` by showing a small overlay on each media thumbnail that displays the **camera** and
the **time** the media was taken. In `web/` the overlay appears on hover (desktop) at the bottom of the
thumbnail (`ImageListItemBar`). Mirror that here.

Unlike `web/`, the overlay must **not** show the date (year, month, day): the grid already groups medias under
a per-day header, so the date is redundant. Show the **camera** and the **time-of-day only**.

## Context

- Reference layout and behaviour: `web/src/components/albums/MediasList/ImageInList.tsx`. It renders a MUI
  `ImageListItemBar` with `title={dateTimeToString(media.time)}` (full date **and** time) and
  `subtitle={media.source ? '@' + media.source : ''}` (the camera model), hidden by default and shown on
  `:hover` at the `lg` breakpoint. We keep the same hover affordance but drop the date half of the title.
- The data is already on the `Media` type in `web-nextjs/domains/catalog/language/catalog-state.ts`:
  - `time: Date` — render the **time of day only** (e.g. a locale time string like `3:45 PM`), never the
    date. Do not reintroduce `toLocaleDateString()`.
  - `source: string` — the camera model. Prefix with `@` like `web/` (e.g. `@Samsung Galaxy S24`). When
    `source` is empty, omit the camera line rather than showing a bare `@`.
- The component to change is
  `web-nextjs/app/(authenticated)/albums/[ownerId]/[albumId]/_components/AlbumMediaGrid/MediaThumbnail.tsx`.
  The thumbnail is a `Link` wrapping an `img`; add the overlay as an absolutely-positioned element inside it.
- Follow the `ui-components` and `nextjs` skills: Storybook for visual testing, no unit test files for pure
  UI. If `MediaThumbnail` has no story yet, add one; otherwise extend the existing `AlbumMediaGrid` stories to
  cover the overlay states.

## Acceptance Criteria

- Each thumbnail shows, on hover (desktop), an overlay with the camera (`source`, `@`-prefixed) and the time
  of day the media was taken.
- The overlay shows **no date** — only the time-of-day and the camera.
- A media with an empty `source` shows only the time (no stray `@`).
- The overlay does not interfere with clicking the thumbnail to open the viewer, and does not cover the video
  play icon in a way that hides it.
- The time formatting is locale-aware (e.g. `media.time.toLocaleTimeString()` trimmed to hours/minutes),
  consistent with the rest of the app.
- Storybook stories cover: image with camera + time, media with empty `source`, and a video thumbnail with the
  overlay.

## Out of scope

- An overlay inside the full-screen viewer (YARL) — this story is thumbnails only.
- Any richer EXIF metadata (aperture, ISO, dimensions) — not available on the `Media` type.
- Changing the per-day grouping header or the grid layout.
