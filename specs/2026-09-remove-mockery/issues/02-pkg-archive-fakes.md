# 02 — `pkg/archive` Fakes

Status: ready
Layer: `pkg/archive`
Depends on: —

## Description

Replace the mockery-generated mocks for the five `pkg/archive` adapters (`ARepositoryAdapter`,
`StoreAdapter`, `CacheAdapter`, `AsyncJobAdapter`, `ResizerAdapter`) with a single `fakes_test.go`
of in-memory Fakes, and convert the four archive test files.

Two failure-injection tests (E1, E2) are KEPT and use `testify/mock` inline (per the feature spec).

See `../spec.md` for the principles.

## Files touched

- `pkg/archive/fakes_test.go` (new)
- `pkg/archive/store_test.go` (rewrite)
- `pkg/archive/relocate_test.go` (rewrite)
- `pkg/archive/content_original_test.go` (rewrite)
- `pkg/archive/content_resized_test.go` (rewrite)

## Fakes to introduce

Create `pkg/archive/fakes_test.go` (package `archive_test`) with:

- **`ARepositoryInMemory`** — implements `archive.ARepositoryAdapter`. Backing state
  `Media map[string]map[string]string` (`owner -> id -> key`). `FindById` returns
  `archive.NotFoundError` on miss. `AddLocation`, `UpdateLocations`, `FindByIds`,
  `FindIdsFromKeyPrefix` all behave like the real one.
- **`StoreInMemory`** — implements `archive.StoreAdapter`. Backing state `Content map[string][]byte`.
  - `Upload(values DestructuredKey, content io.Reader)` derives the key deterministically
    (`values.Prefix + values.Suffix`) and stores the bytes.
  - `Copy(origin, destination)` copies bytes to a deterministic destination key.
  - `Delete(locations)` removes them.
  - `SignedURL(key, duration)` returns `"signed://" + key` if the key exists, `NotFoundError` otherwise.
  - `Download(key)` returns the stored bytes or `NotFoundError`.
  - Test-only accessor: `Has(key) bool` for convenience.
- **`CacheInMemory`** — implements `archive.CacheAdapter`. Similar backing map; `Get` returns
  `archive.NotFoundError` on miss.
- **`AsyncJobInMemory`** — implements `archive.AsyncJobAdapter`. Records calls in
  `LoadedImages [][]*archive.ImageToResize` and `WarmUpCalls []struct{Owner, MissedKey string; Width int}`.
- **`ResizerInMemory`** — implements `archive.ResizerAdapter`. Returns deterministic bytes
  (e.g. `[]byte("resized")`) and a fixed media type (e.g. `"image/jpeg"`).

## Test conversion

Assert by state:

- After `archive.Store(...)`: `ARepositoryInMemory.FindById(owner, id)` returns the expected key,
  `StoreInMemory.Download(key)` returns the expected bytes, `AsyncJobInMemory.LoadedImages` contains
  the expected `ImageToResize` request.
- After `archive.Relocate(...)`: repository points to the new key, old content is gone
  (`StoreInMemory.Has(oldKey) == false`), new content is present.
- After `archive.GetMediaOriginalURL(...)`: the returned URL matches the deterministic `signed://...`.
- After content-resized paths: cache holds the expected bytes.

## Failure-injection tests (KEEP, use `testify/mock` inline)

Two cases must retain error-injection. Use a small inline stub or `testify/mock` for the single
collaborator that fails. Do NOT reintroduce mockery generation.

### E1 — `pkg/archive/relocate_test.go` — "should not delete anything if the index cannot be updated"

The test verifies that when `UpdateLocations` fails, the store's `Delete` is NOT called. Since the
Fake would perform the delete unconditionally, this test needs an `ARepositoryAdapter` whose
`UpdateLocations` returns an error while all other methods keep working. Options:

- Wrap `ARepositoryInMemory` in a thin decorator whose `UpdateLocations` returns the injected error.
- Or use `testify/mock` locally for just this test case's repository.

Verify by state: after the call fails, `StoreInMemory.Has(oldKey)` is still `true` (no delete happened).

### E2 — `pkg/archive/store_test.go` — "should not index the new location if the upload failed"

The test verifies that when `Upload` fails, the repository is NOT updated. Use the same technique
(decorator or inline mock) for `StoreAdapter.Upload`. Verify by state:
`ARepositoryInMemory.FindById(owner, id)` returns `NotFoundError` after the call.

## Success criteria

- No test file in `pkg/archive/` imports `github.com/thomasduchatelle/dphoto/internal/mocks`.
- Only the two failure-injection tests (E1, E2) may import `github.com/stretchr/testify/mock`.
- `go test ./pkg/archive/...` is green.
