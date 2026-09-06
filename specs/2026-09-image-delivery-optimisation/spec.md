# Feature Spec — Image delivery optimisation (4-tier ladder + modern format)

**Status:** draft
**Author:** Arch
**Date:** 2026-09-06

## Intent

Make album and photo pages — especially on mobile — load faster by delivering images that are the right
size for the device and encoded in a modern, lighter format. Today the archive caches only two widths and
serves the same format as the source (typically JPEG), which is both too coarse a size ladder and heavier
than necessary on the wire.

This is a cross-cutting change: it touches the archive backend (`pkg/archive`) which generates and caches
resized images, and the web frontend (`web-nextjs`) which requests them through its custom image loader.

## Background — how it works today

- **Backend (`pkg/archive`)** caches exactly two widths: `MiniatureCachedWidth = 360` and
  `MediumQualityCachedWidth = 2400` (`CacheableWidths`). Any requested width below a cached width is
  downscaled on the fly from the nearest cached size; requests above 2400 are refused.
- **Cache warm-up is lazy:** the first request for a media at a cached width triggers an async job to
  generate and store it; that first request pays the resize cost.
- **Format is preserved from the source:** the resizer (`image_resize/`, using `disintegration/imaging`)
  re-encodes in the original format — a JPEG stays a JPEG. There is no WebP/AVIF path.
- **Frontend** requests images via `GET /api/v1/owners/{owner}/medias/{mediaId}/{filename}?w={width}`. The
  planned NextJS custom loader maps to the legacy widths `360 / 1440 / 2400`.

Two consequences motivate this feature:

1. The size ladder is wrong for the app's real surfaces. `360` is too small for a grid cell on a
   high-DPR phone; the `360 → 2400` jump means phones opening a photo either get an undersized image or a
   needlessly heavy one. There is no tiny tier for a blur placeholder.
2. Serving JPEG leaves 25–50% of bytes on the table versus WebP/AVIF at the same dimensions — the single
   biggest wire-size win available.

## Proposed change

### A 4-tier width ladder

Replace `360 / 2400` with four cached widths, chosen for the two real surfaces (a small grid cell and a
fullscreen viewer) covering device pixel ratios up to 3:

| Tier | Width | Serves |
|------|-------|--------|
| LQIP | 32 | Inline blur placeholder (blur-up), instant grid paint |
| Thumbnail | 480 | Album grid cell (~190 CSS px × DPR 2.5), sharp on all phones |
| Fullscreen (small) | 1280 | Viewer on phones, tablets, and small laptops |
| Fullscreen (large) | 2560 | Viewer on standard, large, and retina desktops |

The frontend loader maps its requested width to the nearest tier at or above it; the viewer requests `1280`
on small screens and `2560` above.

### Modern format

Serve **WebP** for resized tiers (AVIF considered but heavier to encode; WebP is the pragmatic default with
universal support in the app's target browsers). The original-file download path is unchanged and keeps the
source format.

### Blur-up placeholder

The `32` tier exists specifically to be inlined as a blur placeholder so the grid and viewer paint an
immediate blurred image while the real tier loads. How the frontend inlines it (e.g. `next/image`
`blurDataURL`) is an implementation detail for the frontend story.

## User journeys

### Thomas opens a large album on his phone

Thomas opens a 200-photo vacation album over mobile data. The grid paints almost immediately as blurred
placeholders, then thumbnails sharpen in as `480`px WebP images load — only the ones on screen. Scrolling
is smooth; he isn't downloading 200 full-size photos. Tapping a photo opens it fullscreen at `1280`px,
which arrives in well under the 3-second slow-network target.

### Claire browses on a retina laptop

Claire opens the same album on a MacBook. Thumbnails are crisp at `480`, and opening a photo delivers the
`2560` tier so the fullscreen image is sharp on her high-DPR display.

## Scope

**In scope:**

- Change the archive's cached width ladder to `32 / 480 / 1280 / 2560`.
- Add WebP encoding for resized tiers in the resizer.
- Update the frontend custom image loader (`web-nextjs`) to map to the new tiers, and wire the `32` tier as
  the blur placeholder in the grid and viewer.
- A migration/backfill approach for existing media so old albums benefit (see Open Questions).

**Out of scope:**

- AVIF encoding (may be revisited later).
- Changing the API URL contract (`?w=` query parameter stays).
- Changing the original-file download path or its format.
- Serving different tiers based on `Accept` content negotiation — the frontend picks the width explicitly.
- Any catalog / album-metadata changes.

## Open questions

- **Pre-generate vs lazy warm-up.** Today tiers are generated lazily on first request, so the first viewer
  of each media/width pays the resize cost. For a private family library, storage is cheap — should the
  archive **pre-generate all four tiers at backup time** (`pkg/archive` store path) so every read is a plain
  S3 GET? This is the bigger mobile-TTFB win than the exact widths. Recommended: yes.
- **Backfill of existing media.** New tiers only help old albums once regenerated. Do we backfill eagerly
  (batch job over the archive), or keep the lazy warm-up so old media upgrades on first view? Trade-off:
  cost/effort of a one-off batch vs. a slow first view for old photos.
- **Zoom quality.** The YARL Zoom plugin lets users pinch into detail. `2560` will look soft when zoomed on
  a large original. Should the viewer's zoom fetch the **original** on demand, while navigation uses `2560`?
- **Non-resizable media.** PNGs and other formats: encode their tiers as WebP too, or leave as-is? (Source
  download stays original regardless.)
- **AVIF later.** Worth measuring AVIF's extra byte savings against its encode cost as a follow-up, once
  WebP is in place.

## Success criteria

- Grid thumbnails are visually sharp on DPR-3 phones (no upscaling blur) while staying small on the wire.
- A fullscreen photo on a phone loads within the 3-second slow-network target.
- Resized tiers are materially lighter than the current JPEG equivalents (target: ≥25% fewer bytes at the
  same dimensions).
- The album page shows an immediate blurred placeholder before real thumbnails load.
- No regression to the original-file download.
