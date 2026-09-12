# 06 — Cleanup: delete `internal/mocks/`, `Makefile mocks:`, dead test; update docs

Status: ready
Layer: repo-wide
Depends on: 01, 02, 03, 04, 05 (all must be merged)

## Description

Final cleanup once every domain story (01-05) has landed. Removes all mockery artifacts, deletes the
dead daemon test file, updates the developer documentation to describe the new policy, and confirms
the tree is fully green.

See `../spec.md` for the principles.

## Preconditions (verify before starting)

Run these greps first and confirm the results before deleting anything:

```
grep -rn "\"github.com/thomasduchatelle/dphoto/internal/mocks\"" .
```

Expected: no results (01-05 all merged and shipped).

```
grep -rn "\"github.com/stretchr/testify/mock\"" pkg/ cmd/ api/
```

Expected: only 5 files — the KEEP failure-injection tests (E1, E2, E3, E4, E6):

- `pkg/archive/relocate_test.go`
- `pkg/archive/store_test.go`
- `pkg/catalog/album_create_test.go` (two cases share one import)
- `pkg/catalog/album_rename_test.go`

If either grep returns anything unexpected, stop and fix in the appropriate 01-05 story.

## Actions

1. **Delete `internal/mocks/`** — the entire directory (112 files).
   ```
   git rm -r internal/mocks
   ```

2. **Delete `cmd/dphoto/daemon/_udisk_listener_test.go`** — dead code (leading `_` makes Go skip
   it). Confirm nothing else references it via
   `grep -rn "_udisk_listener_test\|VolumeManagerPort" .` before deletion.
   ```
   git rm cmd/dphoto/daemon/_udisk_listener_test.go
   ```

3. **Remove the `mocks:` target from `Makefile`** — lines 208 (the `.PHONY` entry containing
   `mocks`) and 210-214 (the target itself). Keep the `.PHONY` line but drop `mocks` from it.

4. **Run `go mod tidy`** in the repo root and in `api/lambdas/` if applicable. Commit any resulting
   `go.mod` / `go.sum` changes. `stretchr/testify/mock` will likely stay as an indirect via
   `stretchr/testify` — that is expected.

5. **Update `AGENTS.md`** — under "How to build and test?" → "Golang - `pkg` and `cmd/dphoto/`",
   add a paragraph:

   > **No mocks.** Test dependencies are in-memory Fakes co-located in a per-package
   > `fakes_test.go`. One Fake per real backing store, implementing every port that touches that
   > data. Verify by state via the read methods the interface already exposes, or via a purpose-built
   > test-only accessor on the Fake. Do not add `mockery` or `github.com/stretchr/testify/mock`
   > imports in new tests. A small number of existing tests use `testify/mock` inline for
   > failure-behaviour that Fakes cannot naturally express (data-loss / rollback / ordering
   > guarantees); do not extend that pattern without a genuine behavioural reason.

6. **Update `.agents/skills/go/SKILL.md`** — in "How to write a test", sharpen the "fields" bullet:

   > **fields**: structure of the fields of the structure under test (ignore if a function is tested).
   > Use an in-memory **Fake** for each dependency — never a generated mock, never
   > `github.com/stretchr/testify/mock`. One Fake per real backing store, implementing every port
   > that reads or writes that data. Fakes live in a per-package `fakes_test.go` and are shared
   > across tests of the same package. Assertion is by state, via a read method the interface
   > already exposes, or via a purpose-built accessor on the Fake (e.g. `PublishedEvents() []Event`).
   > Do not assert "was this method called".

7. **Verification**:
   ```
   make setup-go
   go test ./...
   cd api/lambdas && go test ./...
   cd deployments/cdk && npm install && npm test
   ```

   All three test suites must be green.

8. **Final grep**:
   ```
   grep -rn "mockery\|internal/mocks" . | grep -v "specs/2026-09-remove-mockery"
   ```
   Expected: no results (the spec directory itself will match — that's OK).

## Success criteria

- `internal/mocks/` no longer exists.
- `Makefile` has no `mocks:` target.
- `cmd/dphoto/daemon/_udisk_listener_test.go` no longer exists.
- `AGENTS.md` and `.agents/skills/go/SKILL.md` describe the new policy.
- `go test ./...` is green.
- No non-spec file in the tree mentions `mockery` or imports `internal/mocks`.
