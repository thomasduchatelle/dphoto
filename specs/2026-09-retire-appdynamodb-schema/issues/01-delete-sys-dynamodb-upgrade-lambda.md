# 01 — Delete orphaned `sys-dynamodb-upgrade` lambda

Status: ready

## Description

The lambda handler at `api/lambdas/sys-dynamodb-upgrade/lambda.go` exists in the source tree but is not wired to any CDK/SST stack, Makefile target, or deployment pipeline. It is the only remaining caller of `appdynamodb.CreateTableIfNecessary` with `localDynamodb=false`, and its removal is already flagged in `specs/archived/todo_cognito.md:8` ("Is the sys-dynamodb lambda is used ? should it be deleted ?"). Delete it.

Removing this lambda is a prerequisite to retiring `pkg/awssupport/appdynamodb/`: while it stays, the legacy helper still has a production-facing caller and cannot be treated as test-only.

## Acceptance criteria

- The directory `api/lambdas/sys-dynamodb-upgrade/` is deleted.
- `grep -r sys-dynamodb-upgrade` (case-insensitive, over the whole repo) returns no matches.
- The pre-existing TODO line in `specs/archived/todo_cognito.md:8` referring to this lambda is removed.
- `make build-api` succeeds.
- `go test ./...` and `cd api/lambdas && go test ./...` pass.

## Out of scope

- Any change to `pkg/awssupport/appdynamodb/` — that is Story 03.
- The `tools/dphotoops/migrator/` decision — that is Story 03.

## References

- `specs/2026-09-retire-appdynamodb-schema/spec.md`
- `api/lambdas/sys-dynamodb-upgrade/lambda.go` (to be deleted)
