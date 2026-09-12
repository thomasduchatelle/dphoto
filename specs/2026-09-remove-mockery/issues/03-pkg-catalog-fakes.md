# 03 — `pkg/catalog` Fakes

Status: ready
Layer: `pkg/catalog`
Depends on: —

## Description

Consolidate the album-repository ports into a single `AlbumRepositoryInMemory` Fake, add the
per-observer Fakes, and convert five test files. Three failure-injection tests are KEPT with
inline `testify/mock`; two pure error-passthrough tests are deleted.

See `../spec.md` for the principles.

## Files touched

- `pkg/catalog/fakes_test.go` (new; may absorb existing scattered `RepositoryAdapterFake`,
  `FindAlbumsByOwnerFake`, `InsertAlbumPortFake`, `CreateAlbumObserverFake`, `AlbumDatesAmendedObserverFake`)
- `pkg/catalog/album_queries_test.go` (drop `mockAdapters` helper, use the new Fake)
- `pkg/catalog/album_create_test.go` (rewrite)
- `pkg/catalog/album_delete_test.go` (rewrite; delete E5a, E5b)
- `pkg/catalog/album_rename_test.go` (rewrite)
- `pkg/catalog/album_amend_dates_test.go` (rewrite)

## Fakes to introduce

Create `pkg/catalog/fakes_test.go` (package `catalog_test`) with:

- **`AlbumRepositoryInMemory`** — one backing store `Albums map[catalog.AlbumId]*catalog.Album`.
  Implements every album-store port used across the tests:
  - `catalog.RepositoryAdapter` (already prototyped as `RepositoryAdapterFake`)
  - `catalog.FindAlbumsByOwnerPort` (already prototyped as `FindAlbumsByOwnerFake`)
  - `catalog.FindAlbumByIdFunc` (as a method or via a wrapper)
  - `catalog.InsertAlbumPort`
  - `catalog.DeleteAlbumRepositoryPort`
  - `catalog.UpdateAlbumNamePort`
  - `catalog.AmendAlbumDateRepositoryPort`
  - `catalog.CountMediasBySelectorsPort` (returns 0 by default; medias can be seeded if a test needs a
    positive count)
  - Assertion is entirely through the read methods this Fake exposes; no `Calls`, no `Verify`.
- **`MediaTransferInMemory`** — implements `catalog.MediaTransfer` and
  `catalog.TransferMediasRepositoryPort`. Records `TransferRecords []catalog.MediaTransferRecords`.
  Returns a canned `catalog.TransferredMedias` (settable via a field) so tests can drive the
  downstream observer.
- **`TimelineMutationObserverInMemory`** — implements `catalog.TimelineMutationObserver`. Records
  `Notifications []catalog.TransferredMedias`.
- **`CreateAlbumObserverInMemory`** — extend the existing `CreateAlbumObserverFake`; recorded
  `CreatedAlbums []catalog.Album`.
- **`DeleteAlbumObserverInMemory`** — implements `catalog.DeleteAlbumObserver`. Records
  `Deleted []catalog.AlbumId`.
- **`RenameAlbumObserverInMemory`** — implements `catalog.RenameAlbumObserver`. Records
  `Renamed []catalog.RenameAlbumRequest` (or the appropriate value type).
- **`AlbumDatesAmendedObserverInMemory`** — extend the existing `AlbumDatesAmendedObserverFake`.

## Test conversion

- **`album_queries_test.go`**: drop `mockAdapters`; tests already use `RepositoryAdapterFake` — align
  it with the new consolidated Fake or keep the existing type name inside the same file, but ensure
  no `internal/mocks` import remains.
- **`album_create_test.go`**: the helpers `expectAlbumInserted`, `stubInsertAlbumPortWithError`,
  `stubFindAlbumsByOwnerWith`, `expectTimelineMutationObserverCalled`, `stubTransferMediaPort`,
  `expectMediaTransferCalled` all fold into direct construction of the Fake plus post-hoc state
  assertions. Where a test asserts "not called", replace with `len(observer.Notifications) == 0`
  after the call.
- **`album_delete_test.go`**: similar; **DELETE tests E5a (line 156) and E5b (line 170)** — pure
  error-passthrough.
- **`album_rename_test.go`**: similar.
- **`album_amend_dates_test.go`**: similar; the four `expect…Called` helpers fold in.

## Failure-injection tests (KEEP, use `testify/mock` inline)

### E3 — `album_create_test.go` — "should not call transfer observer if album insert fails (verify the order)"

Use a small local stub or `testify/mock` for the `InsertAlbumPort` that returns the injected error.
Verify by state: `TimelineMutationObserverInMemory.Notifications` is empty after the call.

### E4 — `album_create_test.go` — "should list the existing albums before creating the new one"

Use a stub `FindAlbumsByOwnerPort` that returns the injected error. Verify by state:
`AlbumRepositoryInMemory.Albums` is unchanged (album was not inserted).

### E6 — `album_rename_test.go` — "should interrupt the transfer if the album insertion fails"

Use a stub `InsertAlbumPort` returning the injected error. Verify by state: the old album is still
present in the repository (rollback), `MediaTransferInMemory.TransferRecords` is empty,
`TimelineMutationObserverInMemory.Notifications` is empty.

## Tests to DELETE

- `album_delete_test.go:156` — "should fails if listing albums raises an error" (E5a).
- `album_delete_test.go:170` — "should fails if listing albums raises an error" (E5b).

Remove the local `anExpectedError` variable and helpers if they become unused after the deletion.

## Success criteria

- No test file in `pkg/catalog/` imports `github.com/thomasduchatelle/dphoto/internal/mocks`.
- Only E3, E4, E6 import `github.com/stretchr/testify/mock` (inline).
- `go test ./pkg/catalog/...` is green.
