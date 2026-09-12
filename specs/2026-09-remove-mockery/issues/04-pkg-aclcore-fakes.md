# 04 — `pkg/acl/aclcore` Fakes

Status: ready
Layer: `pkg/acl/aclcore`
Depends on: —

## Description

Introduce a consolidated `fakes_test.go` for `pkg/acl/aclcore` and convert six test files. One
pure error-passthrough test is deleted. No failure-injection tests are kept in this package.

See `../spec.md` for the principles.

## Files touched

- `pkg/acl/aclcore/fakes_test.go` (new)
- `pkg/acl/aclcore/rules_test.go` (rewrite)
- `pkg/acl/aclcore/user_create_test.go` (rewrite)
- `pkg/acl/aclcore/identity_queries_test.go` (rewrite)
- `pkg/acl/aclcore/refresh_token_generator_test.go` (rewrite)
- `pkg/acl/aclcore/authenticate_refresh_token_test.go` (rewrite; delete E8)
- `pkg/acl/aclcore/authenticate_sso_test.go` (rewrite)

## Fakes to introduce

Create `pkg/acl/aclcore/fakes_test.go` (package `aclcore_test`) with:

- **`ScopeRepositoryInMemory`** — one backing `map[usermodel.UserId][]*aclcore.Scope`. Implements
  `aclcore.ScopesReader`, `aclcore.ScopeWriter`, `aclcore.IdentityQueriesScopeRepository`, and
  `aclcore.ReverseScopesReader`. Assertion is via the reader methods (`ListScopesByUser`,
  `FindScopesById`).
- **`IdentityRepositoryInMemory`** — one backing `map[usermodel.UserId]aclcore.Identity`. Implements
  `aclcore.IdentityDetailsStore` and `aclcore.IdentityQueriesIdentityRepository`. Assertion via
  `FindIdentity` (already exposed).
- **`RefreshTokenRepositoryInMemory`** — one backing `map[string]aclcore.RefreshTokenSpec`.
  Implements `aclcore.RefreshTokenRepository`. `FindRefreshToken` returns
  `aclcore.InvalidRefreshTokenError` on miss. `HouseKeepRefreshToken` uses a settable `NowFunc` so
  expired specs can be removed. Assertion via `FindRefreshToken`.
- **`AccessTokenGeneratorInMemory`** — implements `aclcore.IAccessTokenGenerator`. Returns a
  deterministic `aclcore.Authentication{AccessToken: "at-"+string(email), ExpiresIn: 42, ExpiryTime: <fixed>}`.
  Records `GeneratedFor []usermodel.UserId` for optional assertions.
- **`RefreshTokenGeneratorInMemory`** — implements `aclcore.IRefreshTokenGenerator`. Returns a
  deterministic token (e.g. `"rt-" + spec.Email + "-" + spec.RefreshTokenPurpose`). Records
  `GeneratedFor []aclcore.RefreshTokenSpec` so tests can assert what was requested.

### Important: keep the real `AccessTokenGenerator` and `RefreshTokenGenerator` for their own tests

- `refresh_token_generator_test.go` tests the **real** `aclcore.RefreshTokenGenerator`; only the
  `RefreshTokenRepository` collaborator is Faked. Assertion becomes: after
  `generator.GenerateRefreshToken(spec)`, call `RefreshTokenRepositoryInMemory.FindRefreshToken(got)`
  and assert the spec matches.
- `authenticate_sso_test.go` inspects a real signed JWT; keep the **real** `aclcore.AccessTokenGenerator`
  and only Fake its `PermissionsReader` (= `ScopesReader`, via `ScopeRepositoryInMemory`) and its
  `RefreshTokenGenerator` collaborator (via `RefreshTokenGeneratorInMemory`).

## Test conversion

- **`rules_test.go`**: seed `ScopeRepositoryInMemory` with the scope, call `CoreRules.Owner()`,
  assert the returned owner. Empty state naturally exercises the "no scopes" path.
- **`user_create_test.go`**: seed `ScopeRepositoryInMemory` with the pre-existing scopes; after
  `UserCreate.Create(...)`, verify with `ListScopesByUser` that the expected scopes now exist and
  the ones that should have been deleted are gone.
- **`identity_queries_test.go`**: seed both `ScopeRepositoryInMemory` and `IdentityRepositoryInMemory`;
  call the queries and assert.
- **`refresh_token_generator_test.go`**: use the real `RefreshTokenGenerator` with
  `RefreshTokenRepositoryInMemory`; after generation, verify with `FindRefreshToken`.
- **`authenticate_refresh_token_test.go`**: seed `RefreshTokenRepositoryInMemory` with the
  starting token, seed `IdentityRepositoryInMemory` where needed. After
  `AuthenticateFromRefreshToken`, verify: old token gone (`FindRefreshToken` → not found), new token
  exists (find by `Authentication.RefreshToken`), `RefreshTokenGeneratorInMemory.GeneratedFor`
  contains the expected spec, `AccessTokenGeneratorInMemory.GeneratedFor` contains the email.
- **`authenticate_sso_test.go`**: seed `ScopeRepositoryInMemory` with the user's scopes; keep the
  real `AccessTokenGenerator` so the JWT round-trip still works; verify by parsing the returned JWT
  as today, plus assert `RefreshTokenGeneratorInMemory.GeneratedFor` contains the expected spec and
  `IdentityRepositoryInMemory.FindIdentity` returns the stored identity.

## Failure-injection tests

None kept in this package.

## Tests to DELETE

- `authenticate_refresh_token_test.go:134` — "should pass through error with identity repository"
  (E8). Pure `return err` passthrough.

## Success criteria

- No test file in `pkg/acl/aclcore/` imports `github.com/thomasduchatelle/dphoto/internal/mocks`
  or `github.com/stretchr/testify/mock`.
- `go test ./pkg/acl/aclcore/...` is green.
