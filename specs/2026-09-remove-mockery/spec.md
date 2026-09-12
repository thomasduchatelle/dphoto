# Remove Mockery — replace with in-memory Fakes

## Intent

Stop using the generated `mockery` mocks (under `internal/mocks/`) and, more broadly, stop writing tests
against mock libraries. Every dependency in a test becomes a **Fake**: a small in-memory implementation
that behaves like the real adapter.

The goal is tests that read as small integration tests within their package — exercising real behaviour
of the collaborators, verifying outcomes by state, not by "was this method called".

## Principles

1. **No `mockery`.** The `internal/mocks/` directory is deleted; the `Makefile mocks:` target is removed.
2. **Dependencies are Fakes.** In-memory implementations that behave like the real one. **One Fake per
   real backing store**, implementing every port that reads or writes that data (e.g. a single
   `AlbumRepositoryInMemory` implements `InsertAlbumPort`, `FindAlbumsByOwnerPort`,
   `AmendAlbumDateRepositoryPort`, `DeleteAlbumRepositoryPort`, `UpdateAlbumNamePort`,
   `TransferMediasRepositoryPort`, `CountMediasBySelectorsPort`).
3. **Fakes are per package.** They live in a single `pkg/<pkg>/fakes_test.go` (test-only compilation)
   and are shared across all tests of that package. Cross-package reuse is via copy-paste, not
   inter-package imports.
4. **No "was-method-called" assertions.** Verify by state:
   - Preferred: assert through a **read method the interface already exposes** (after `SaveScope`,
     assert with `FindScopesById`).
   - Fallback: add a **test-only accessor** on the Fake (e.g. `PublishedEvents() []Event`).
5. **"Not found" errors are natural** — the Fake returns `NotFoundError`, `AlbumNotFoundErr`,
   `IdentityDetailsNotFoundError`, `InvalidRefreshTokenError`, `CertificateNotFoundError`, etc. from an
   empty state. No injection needed.
6. **`testify/mock` may be used inline in a handful of failure-behaviour tests** where the code under
   test must react to a genuine error from a collaborator (data-loss safety, ordering, rollback). Those
   cases are enumerated below. `testify/mock` is used locally, no generation, no `internal/mocks/`.
7. **Pure error-passthrough tests are deleted** — they only assert Go's `return err` operator.

## Scope

In scope:

- `pkg/dns`
- `pkg/archive`
- `pkg/catalog`
- `pkg/acl/aclcore`
- `pkg/acl/catalogacl`
- `internal/mocks/` (deleted at the end)
- `Makefile` (`mocks:` target removed)
- `cmd/dphoto/daemon/_udisk_listener_test.go` (dead code — leading `_` makes Go skip it; deleted at the end)
- `AGENTS.md` and `.agents/skills/go/SKILL.md` (documented policy update)

Out of scope:

- `api/lambdas/` — already free of `testify/mock`.
- `cmd/dphoto/` — no live tests use mocks (only the ignored `_udisk_listener_test.go`).
- All TypeScript projects.
- Refactors to production code (interfaces, packages) — only test code changes here.

## Inventory of failure-injection tests

Genuine failure-behaviour tests that KEEP `testify/mock` inline (5 cases):

| # | Test                                                                                            | Why it's kept                                                                              |
|---|-------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------|
| E1 | `pkg/archive/relocate_test.go` — "should not delete anything if the index cannot be updated"   | Data-loss safety: aborts delete after failed index update.                                 |
| E2 | `pkg/archive/store_test.go` — "should not index the new location if the upload failed"         | Data-loss safety: no orphaned index entries when upload fails.                             |
| E3 | `pkg/catalog/album_create_test.go` — "should not call transfer observer if album insert fails" | Event-ordering guarantee: no events for phantom albums.                                    |
| E4 | `pkg/catalog/album_create_test.go` — "should list existing albums before creating"             | Ordering guarantee: list before insert.                                                    |
| E6 | `pkg/catalog/album_rename_test.go` — "should interrupt the transfer if the album insertion fails" | Rollback semantics.                                                                     |

Pure error-passthrough tests that are DELETED (4 cases):

| #  | Test                                                                                             | Why it's deleted                                          |
|----|--------------------------------------------------------------------------------------------------|-----------------------------------------------------------|
| E5a | `pkg/catalog/album_delete_test.go:156` — "should fails if listing albums raises an error"       | Pure `return err` passthrough — no behaviour to lock in. |
| E5b | `pkg/catalog/album_delete_test.go:170` — "should fails if listing albums raises an error"       | Same.                                                     |
| E7  | `pkg/acl/catalogacl/case_share_album_test.go:82` — "should passthroughs an other error"         | Same, redundant with the `AlbumNotFoundErr` case above.  |
| E8  | `pkg/acl/aclcore/authenticate_refresh_token_test.go:134` — "should pass through error with identity repository" | Same. |

## Structure

Six stories: five domain stories (01-05, one per package, fully parallel) and one cleanup story
(06, depends on 01, 02, 03, 04, 05).

See `stories.md` for the overview and `issues/` for the individual issue files.
