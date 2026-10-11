# Stories — Remove Mockery

Six stories: five domain conversions (independent, fully parallel) and one cleanup (blocks on all
five).

## Dependency graph

```
   (all five parallel, no deps)
   ┌────────┬────────┬────────┬────────┬────────┐
   │   01   │   02   │   03   │   04   │   05   │
   │  dns   │archive │catalog │aclcore │catalog-│
   │        │        │        │        │  acl   │
   └────┬───┴────┬───┴────┬───┴────┬───┴────┬───┘
        └────────┴────────┴────────┴────────┘
                          │
                          ▼
                     ┌─────────┐
                     │   06    │  cleanup (depends on 01, 02, 03, 04, 05)
                     └─────────┘
```

## Story overview

| ID | Package                | Depends on           | Files touched                                                                                                        | Fakes to introduce                                                                                                                                                                                                     | Failure-injection tests |
|----|------------------------|----------------------|----------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------|
| 01 | `pkg/dns`              | —                    | `certificate_manager_test.go`; new `fakes_test.go`                                                                   | `CertificateManagerInMemory`, `CertificateAuthorityInMemory`                                                                                                                                                             | none                    |
| 02 | `pkg/archive`          | —                    | `store_test.go`, `relocate_test.go`, `content_original_test.go`, `content_resized_test.go`; new `fakes_test.go`      | `ARepositoryInMemory`, `StoreInMemory`, `CacheInMemory`, `AsyncJobInMemory`, `ResizerInMemory`                                                                                                                          | E1, E2 (inline mock)    |
| 03 | `pkg/catalog`          | —                    | `album_queries_test.go`, `album_create_test.go`, `album_delete_test.go`, `album_rename_test.go`, `album_amend_dates_test.go`; new/enlarged `fakes_test.go` | `AlbumRepositoryInMemory` (consolidates all album-store ports), `MediaTransferInMemory`, `TimelineMutationObserverInMemory`, `DeleteAlbumObserverInMemory`, `RenameAlbumObserverInMemory`, `AlbumDatesAmendedObserverInMemory` | E3, E4, E6 (inline mock); DELETE E5a, E5b |
| 04 | `pkg/acl/aclcore`      | —                    | `rules_test.go`, `user_create_test.go`, `identity_queries_test.go`, `refresh_token_generator_test.go`, `authenticate_refresh_token_test.go`, `authenticate_sso_test.go`; new `fakes_test.go` | `ScopeRepositoryInMemory`, `IdentityRepositoryInMemory`, `RefreshTokenRepositoryInMemory`, `AccessTokenGeneratorInMemory`, `RefreshTokenGeneratorInMemory`                                                                | DELETE E8               |
| 05 | `pkg/acl/catalogacl`   | —                    | `case_share_album_test.go`; new `fakes_test.go`                                                                      | `FindAlbumPortInMemory`, `ScopeRepositoryInMemory` (copy-pasted from 04)                                                                                                                                                 | DELETE E7               |
| 06 | cleanup                | 01, 02, 03, 04, 05   | delete `internal/mocks/`, delete `cmd/dphoto/daemon/_udisk_listener_test.go`, remove `Makefile mocks:` target, update `AGENTS.md` and `.agents/skills/go/SKILL.md`, `go mod tidy` | —                                                                                                                                                                                                                       | —                       |

## Success criteria (whole feature)

- No file under `pkg/`, `cmd/`, or `api/` imports `github.com/thomasduchatelle/dphoto/internal/mocks`.
- `internal/mocks/` no longer exists.
- `Makefile` no longer contains a `mocks:` target.
- `grep -rn "\"github.com/stretchr/testify/mock\"" pkg/ cmd/ api/` returns only the 5 KEEP tests
  (E1, E2, E3, E4, E6) — each with a small inline usage.
- `make setup-go && go test ./...` is green.
- `cd api/lambdas && go test ./...` is green.
- `AGENTS.md` and `.agents/skills/go/SKILL.md` explicitly forbid new mocks and describe the Fake
  policy.
