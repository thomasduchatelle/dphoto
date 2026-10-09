# 01-01 — AlbumSummary repository (DynamoDB + in-memory)

Status: done
Phase: 1
Layer: catalog domain — `pkg/catalogviews` (model + in-memory fake) + `pkg/catalogviewsadapters/catalogviewsdynamodb`
Depends on: —

## Description

Reshape the persistence layer of the album-list view so it can carry a full album projection (identity +
display fields + count) and expose a set of write operations whose attribute footprints are disjoint —
this is what allows the count path and the display-field path to coexist without clobbering each other in
subsequent tickets.

No production wiring changes in this ticket: the existing `CommandHandlerAlbumSize`, providers, and
factories keep compiling against the new types (adapter method names updated where needed, semantics
preserved). All behaviour changes land in `01-02` and beyond.

See `../design.md` (Read model) and the interface contract agreed for `AlbumView` (parent ticket).

## Naming

Rename across `pkg/catalogviews` and `pkg/catalogviewsadapters/catalogviewsdynamodb`:

| Old                             | New                                    |
|---------------------------------|----------------------------------------|
| `AlbumSize`                     | `AlbumSummary`                         |
| `UserAlbumSize`                 | `UserAlbumSummary`                     |
| `MultiUserAlbumSize`            | `AlbumSummaryForUsers`                 |
| `AlbumSizeDiff`                 | `AlbumMediaCountDiff`                  |
| `AlbumSizeRecord` (DDB)         | `AlbumSummaryRecord`                   |
| `AlbumSizeInMemoryRepository`   | `AlbumSummaryInMemoryRepository`       |

## DynamoDB record

- PK: `USER#{EMAIL}#ALBUMS_VIEW` (unchanged).
- SK: `{OWNED|VISITOR}#{OWNER}#{FOLDER_NAME}` — **`#COUNT` suffix dropped**.
- Attributes: `AlbumOwner`, `AlbumFolderName`, `AvailabilityType`, `UserId`, `Count`,
  and new `AlbumName` (S), `AlbumStart` (S, RFC3339), `AlbumEnd` (S, RFC3339).

## Acceptance criteria

- `AlbumSummary` domain type carries `AlbumId`, `MediaCount`, `Name`, `Start`, `End`.
- New `AlbumSummaryRepository` interface exposing:
  - `ListSummariesForUser(ctx, userId) ([]UserAlbumSummary, error)` — read all rows of a user.
  - `PutSummaries(ctx, []AlbumSummaryForUsers) error` — full-row upsert (Put); used by drift rebuild and
    by the future `AlbumShared` event.
  - `SetDisplayFields(ctx, albumId, users []Availability, name, start, end) error` — UpdateItem
    `SET AlbumName/AlbumStart/AlbumEnd` per user; does **not** touch `Count`.
  - `IncrementCounts(ctx, []AlbumMediaCountDiff) error` — UpdateItem `ADD Count :d`; does **not** touch
    display fields (renamed from `UpdateAlbumSize`, semantics preserved).
  - `SetCounts(ctx, []AlbumMediaCountForUsers) error` — UpdateItem `SET Count = :c`; does **not** touch
    display fields. Used by the future recount-after-transfer path.
  - `DeleteRow(ctx, availability, albumId) error` — one row (renamed from `DeleteAlbumSize`).
  - `DeleteAllRowsForAlbum(ctx, albumId) error` — every viewer row for the album; implemented as a
    fan-out of `DeleteRow` based on `ListUsersWhoCanAccessAlbum` at the call site is acceptable, but the
    method itself lives on the repository so the adapter can scan the projection directly when needed.
- `GetAvailabilitiesByUser` and `GetAlbumSizes` are renamed accordingly (`ListSummariesForUser` /
  `ListSummariesForUserAndOwners`) — no semantic change; drift reconciler switches to the new names.
- `AlbumSummaryInMemoryRepository` implements the same interface, used as the fake in later tickets.
- DDB adapter tests (`view_test.go`) cover every method above. Two new tests are mandatory:
  1. `SetDisplayFields` does not change `Count` on an existing row (seed with `Count=42`; assert
     `Count=42` after).
  2. `IncrementCounts` does not change `AlbumName/AlbumStart/AlbumEnd` on an existing row (seed with
     display fields; assert they survive).
- Existing DDB test helpers (`albumSizeItem`) extended to accept the new attributes; existing cases keep
  their `wantAfter` bodies enriched with display fields where relevant.
- Old provider files (`albums_view_owned_provider.go`, `albums_view_shared_provider.go`,
  `albums_view_provider.go`) still compile: they consume the new type names but keep their present
  behaviour until `01-02` replaces them.

## Out of scope

- Any change to `AlbumView` behaviour or its `ListAlbums`.
- Any change to the event-handling logic in `CommandHandlerAlbumSize` (its methods keep calling the new
  repository methods with equivalent semantics).
- Removal of the read-path providers (happens in `01-02`).
- Drift-reconciler display-field backfill (happens in `01-08`).
- Signature change of `AlbumSharedObserver` (happens in `01-07`).

## References

- `../spec.md`, `../design.md`
- Parent ticket `01-album-list-read-model.md` (interface contract).
- Load skills: `go`, `architecture`.
