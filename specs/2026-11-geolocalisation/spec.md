# Media geolocalisation — see and find photos by where they were taken

**Author:** Arch
**Date:** 2026-10-26
**Status:** draft — captured as an idea, **to be refined in a later session** before breaking into issues.

## Intent

Let the user explore their library by **place**: see where photos were taken on a map, find all medias taken
in a given area, and jump from a single media to its location. Today the catalog carries no location data at
all (the `Media` model has `time`, `type`, and paths — no GPS), so this feature is as much a **backend /
data-model** effort as a UI one.

This is **not** part of the NextJS cutover (`specs/2026-10-nextjs-cutover/`), which is deliberately parity-only
and excludes backend changes. Geolocalisation is a net-new capability tracked separately here.

## Capabilities (in rough priority order)

1. **Link to Google Maps from a media (MVP).** From the full-screen viewer, when a media has an embedded
   location, offer a link that opens Google Maps at its coordinates. Smallest slice that delivers value and
   forces the end-to-end plumbing (extract GPS → expose it → show it). Good first increment.
2. **Map of an album.** Show where the medias of an album were taken on a map, so the user can see the
   geographic spread of a trip at a glance.
3. **Search medias by location.** Find all medias taken in a certain place / area, across albums.

Media **details** (camera, capture time, lens, etc.) surfaced alongside location in the viewer are a natural
companion and share the same EXIF-extraction dependency; fold them in when refining, or split out if they
distract from the geo story.

## Known hard problems (to refine later — new session)

These are flagged now so they are not forgotten; none are solved yet:

- **Albums spanning very different places.** A single album may contain medias from widely separated
  locations (e.g. a multi-leg trip). "The album's location" is not a single point — the album map and any
  per-album summary must handle clusters / multiple regions gracefully rather than averaging to a meaningless
  centroid.
- **Medias without embedded geolocation.** Many photos (scans, screenshots, stripped EXIF, older cameras)
  carry no GPS. Decide the fallback: omit from the map, infer from album/neighbours, or let the user set a
  location. The map and search must degrade cleanly when location is missing or partial.

## Open questions for refinement

- **Where does location come from?** EXIF GPS at backup/indexing time (`pkg/backup`), re-extracted from the
  archive, or a one-off backfill of existing medias? Which context owns the data (`catalog`? a new one?) and
  how it reaches the API and the `Media` model.
- **Data model & indexing.** How location is stored (DynamoDB single-table — see `DATA_MODEL.md`) and whether
  "search by place" needs a geo-index or can be satisfied by a coarse bucketing scheme kept cheap to operate.
- **Privacy / cost.** GPS is sensitive; map tiles and geocoding have cost and third-party dependencies. The
  MVP (a plain Google Maps *link*) avoids embedding a map SDK — keep later increments cost-conscious per the
  repo's "low operating cost" priority.

## Explicitly out of scope (for now)

- A world/heat map of the entire library.
- Reverse-geocoding to place names, trip auto-naming, or clustering into "visits".
- Editing / correcting a media's location by hand (may become part of the "missing geolocation" answer later).
