# ADR-0004: Album-list view is a complete read model

- Status: accepted
- Date: 2026-09-23
- Scope: `pkg/catalogviews` + `pkg/catalogviewsadapters/catalogviewsdynamodb`

## Context

Before this change, `GET /albums` served two data needs from two partitions. The per-user view
partition (`USER#{EMAIL}#ALBUMS_VIEW`) carried a count-only row per accessible album, and the
album's display fields (Name, Start, End) had to be fetched from the album metadata partition
(`{OWNER}#ALBUM`). Owned albums were retrieved by a Query on the user's `{OWNER}#ALBUM` partition;
shared albums went through a BatchGet by album ids built from the user's ACL scopes. The two
result sets were then merged with the count rows into a single response.

This split was fine when the view carried only counts, but the covers feature (Phase 2+) needs the
view to carry richer per-album fields. Adding more denormalised fields to a count-only row would
have kept the awkward two-partition merge and the divergent owned/shared read paths.

## Decision

- The per-user view record becomes a **complete projection** of the album: identity (owner + folder
  name), display fields (name, start, end), and the media count. The SK loses its `#COUNT` suffix
  (`{OWNED|VISITOR}#{OWNER}#{FOLDER_NAME}`): the row is no longer count-only.
- The write path is split into **disjoint attribute sets** so count updates and display-field
  updates never clobber each other:
  - `PutSummaries` writes the full projection when a viewer's row is first materialised (create,
    share).
  - `SetDisplayFieldsForAllViewers` updates the display fields on every viewer's row (rename in
    place, amend dates).
  - `IncrementCountForAllViewers` performs the atomic `ADD Count {diff}` per viewer (medias
    inserted, transferred).
  - `SetCountForAllViewers` overwrites the count when a full recount is needed (media transfers
    changing the totals of neighbouring albums).
- The read path collapses to **a single Query per user** (`ListSummariesForUser`) for both owned
  and shared albums. The owner's "shared with" grid is still a separate Query — unchanged, its
  denormalisation is out of scope for this ADR.
- Viewer discovery — "who has a row for this album?" — is answered by a projection GSI
  (`AlbumViewByAlbumIndex`) whose PK is derived from the album, so an album mutation can fan out
  to every viewer's row without an ACL round-trip.
- Drift reconciliation rebuilds the display fields on every viewer row from the canonical album
  metadata partition (`OwnerDriftReconciler`), so a legacy row that still has a stale name or
  dates can be repaired without replaying events.

## Consequences

- Viewer records must be **fanned out** on album create, rename-in-place, amend-dates, and share,
  not just when a media count changes. Every use case in `pkg/catalog` fires an event that the
  view observes.
- `list-albums` no longer touches the album metadata partition. The DTO returned by the lambda is
  unchanged.
- The owned/shared read providers (`albums_view_owned_provider.go`,
  `albums_view_shared_provider.go`, and the `providers.go` shell) are gone, replaced by the single
  `Repository.ListSummariesForUser` call inside `AlbumView.ListAlbums`.
- A one-off post-deploy `dphotops drift` reconciliation backfills legacy `#COUNT`-suffixed rows
  with the new display fields (see `cmd/dphotops/cmd/drift.go`).
- The view is now on the critical path of every album write, not just media inserts. The disjoint
  attribute sets and the projection GSI are what keep this safe: partial failures do not corrupt
  the count, and drift reconciliation is the safety net for anything that slips through.
