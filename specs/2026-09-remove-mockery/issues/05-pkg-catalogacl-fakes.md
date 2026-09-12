# 05 — `pkg/acl/catalogacl` Fakes

Status: ready
Layer: `pkg/acl/catalogacl`
Depends on: —

## Description

Convert `pkg/acl/catalogacl/case_share_album_test.go` to Fakes. One pure error-passthrough test is
deleted. The existing `AlbumSharedObserverFake` stays.

`ScopeRepositoryInMemory` is copy-pasted from 04's spec (Fakes live in `*_test.go` and cannot be
imported across packages; the copy is small and independently maintainable).

See `../spec.md` for the principles.

## Files touched

- `pkg/acl/catalogacl/fakes_test.go` (new)
- `pkg/acl/catalogacl/case_share_album_test.go` (rewrite; delete E7)

## Fakes to introduce

Create `pkg/acl/catalogacl/fakes_test.go` (package `catalogacl_test`) with:

- **`FindAlbumPortInMemory`** — one backing `map[catalog.AlbumId]*catalog.Album`. Implements
  `catalogacl.FindAlbumPort`. Returns `catalog.AlbumNotFoundErr` on miss.
- **`ScopeRepositoryInMemory`** — copy-paste from 04's design. Backing
  `map[usermodel.UserId][]*aclcore.Scope`. Implements `aclcore.ScopeWriter` (the only interface used
  here) plus enough of `aclcore.ScopesReader` to verify state via `FindScopesById` /
  `ListScopesByUser`.
- Keep the existing `AlbumSharedObserverFake`.

## Test conversion

- Case "should create the ACL rule when the album exists": seed `FindAlbumPortInMemory` with the
  album. After `ShareAlbumWith(...)`, assert:
  - the observer received the sharing (`AlbumSharedObserverFake.Shared`),
  - the scope now exists in `ScopeRepositoryInMemory` via `FindScopesById(expectedScopeId)`.
- Case "should return an error if the album doesn't exist": leave `FindAlbumPortInMemory` empty;
  natural `AlbumNotFoundErr`. Assert the error and that no scope was saved.

## Failure-injection tests

None kept in this package.

## Tests to DELETE

- `case_share_album_test.go:82` — "should passthroughs an other error" (E7). Pure `return err`
  passthrough, redundant with the `AlbumNotFoundErr` case just above it.

## Success criteria

- `pkg/acl/catalogacl/case_share_album_test.go` no longer imports
  `github.com/thomasduchatelle/dphoto/internal/mocks` or `github.com/stretchr/testify/mock`.
- `go test ./pkg/acl/catalogacl/...` is green.
