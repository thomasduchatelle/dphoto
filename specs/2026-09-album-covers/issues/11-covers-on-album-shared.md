# 11 — Carry covers on `AlbumShared`

Status: ready
Phase: 2
Layer: ACL — `pkg/acl/catalogacl` + `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: —

## Description

When an album is shared with a new visitor, the visitor's view row must be written with the current
cover set in one shot so the next `GET /albums` returns the album with covers (no second denormalisation
pass, no reliance on drift). The share use case already reads the album; it reads the covers too and
passes them on the `AlbumSharedObserver` call.

Touches `pkg/acl/catalogacl` + `AlbumView.OnAlbumShared`; different files than issues 09 and 10, so
these three can be worked in parallel. Unlike 10, this ticket does not need the upgraded
`RandomizeCoversPort` — it only reads the current covers via `FindCoversByAlbum` (already shipped).

See `../design.md` → Event-payload matrix (row `AlbumShared`), §Cover projection in the view.

## Acceptance criteria

- `pkg/acl/catalogacl`:
  - `AlbumSharedObserver` interface signature changes:
    - Old (from issue 01-07): `AlbumShared(ctx, album catalog.Album, userEmail usermodel.UserId) error`.
    - New: `AlbumShared(ctx, album catalog.Album, covers []catalog.Cover, userEmail usermodel.UserId) error`.
  - `ShareAlbumCase.ShareAlbumWith` reads the current cover set via
    `CoverRepository.FindCoversByAlbum(ctx, album.AlbumId)` (new dependency, port on the case) and
    passes the result as `covers` to every registered `AlbumSharedObserver`.
- `pkg/catalogviews`:
  - `AlbumView.OnAlbumShared(ctx, album, covers, userId)` writes one full visitor row via
    `PutSummaries` — `Name/Start/End` from `album`, `Count` from `MediaCounterPort.CountMedia`,
    `Covers` from the new argument, `Availability = VisitorAvailability(userId)`.
  - `AlbumViewAlbumSharedObserver` adapter signature updated to match the new interface.
- `pkg/pkgfactory`:
  - `ShareAlbumCase` wires in the `CoverRepository` (reuse `pkgfactory.CatalogRepository`).
  - The AlbumView adapter registration is updated for the new signature.
- `AlbumUnSharedObserver` is unchanged (unshare does not need covers).
- Tests:
  - `pkg/acl/catalogacl/case_share_album_test.go`: `AlbumSharedObserverFake` signature updated; new
    assertion that the observer receives the covers read from the `CoverRepository`.
  - `pkg/catalogviews/albums_view_events_test.go`:
    - `TestAlbumView_AlbumShared/it_should_add_a_visitor_row_with_covers_and_count` — share an album
      with 2 covers; assert `ListAlbums` for the visitor returns the album with those covers and the
      count from `MediaCounterPort`.
    - `TestAlbumView_AlbumShared/it_should_write_an_empty_covers_slice_when_the_album_has_none` —
      ensure empty is handled correctly.
  - Acceptance suite (`01-09`) extended: after `ShareAlbum`, the visitor's `ListAlbums` returns the
    album with the owner's current covers.
- **Dead-code cleanup — conditional**: `AlbumView.OnAlbumCoversChanged` is dead once issues 09, 10
  and 11 have all landed. Whichever of the three lands **last** must verify
  (`grep OnAlbumCoversChanged`) that no production caller remains and, if so, delete the method and
  its tests from `pkg/catalogviews/albums_view.go` and `pkg/catalogviews/albums_view_events_test.go`
  as part of the same PR. The repository method `SetCoversForAllViewers` is kept (used by 10 and by
  drift).

## Out of scope

- `AlbumUnshared` — unchanged.
- Any other lifecycle event — issues 09, 10, 12.

## References

- `../spec.md` (Automatic cover maintenance → album shared).
- `../design.md` (Event-payload matrix).
- Load skills: `go`, `architecture`.
