# 03 — Retire `pkg/awssupport/appdynamodb/` and resolve the migrator

Status: ready

## Description

With `schema.json` in place (Story 02) and the orphaned lambda removed (Story 01), the remaining callers of `appdynamodb.CreateTableIfNecessary` are all dev/test scoped: Localstack CLI bootstrap, four test files, and the `dphotoops` migrator tool. Replace them with a test-scoped helper and delete the `appdynamodb` package.

Introduce `CreateLocalTable(ctx context.Context, client *dynamodb.Client, tableName string) error` (name negotiable) inside `pkg/awssupport/dynamotestutils/`. It reads the same `schema.json` produced in Story 02 and calls `dynamoutils.CreateOrUpdateTable`. It always sets local-DynamoDB-friendly provisioned throughput (the current `localDynamodb=true` branch); there is no production caller left.

The `dphotoops` migrator (`tools/dphotoops/migrator/transformattion_indexes.go`) is the one non-test caller. Decide its fate as part of this story:
- If the migrator tool is no longer used, delete `transformattion_indexes.go`. If the migrator becomes empty as a result, delete the tool.
- If the migrator is still used, drop the `CreateTableIfNecessary` call from it and document in the migrator's usage notes (or its README) that `cdk deploy` must run first so the table schema is up to date.

Record the decision (and evidence supporting it) in this issue's Comments section before merging.

## Acceptance criteria

- New helper `dynamotestutils.CreateLocalTable` (or equivalent name) exists and reads `pkg/awssupport/appdynamodb/schema.json`.
- The following callers use the new helper instead of `appdynamodb.CreateTableIfNecessary`:
  - `pkg/awssupport/dynamotestutils/context.go`
  - `pkg/catalogadapters/catalogdynamo/repository_album_test.go`
  - `pkg/catalogadapters/catalogdynamo/repository_media_crud_test.go`
  - `pkg/backup_acceptance_test.go`
  - `cmd/dphoto/bootstrap/bootstrap_catalog.go` (the `cfg.GetBool(config.Localstack)` branch)
- Explicit decision on `tools/dphotoops/migrator/` recorded in this issue's Comments section, and executed:
  - If deleted: `tools/dphotoops/migrator/transformattion_indexes.go` is gone; if the tool becomes empty, the tool is gone too.
  - If kept: the `CreateTableIfNecessary` call is removed and the prerequisite (`cdk deploy` before running migrations) is documented.
- `pkg/awssupport/appdynamodb/create_update_table.go` is deleted. If the `appdynamodb` package directory becomes empty, it is removed (the `schema.json` moves to its new home under `dynamotestutils/` or `deployments/cdk/`, whichever the author judges best given who reads it — CDK still needs access).
- `grep -r CreateTableIfNecessary` returns matches only inside `pkg/awssupport/dynamoutils/` (the unrelated generic helper).
- `go test ./...` passes.
- `cd api/lambdas && go test ./...` passes.
- `make build-api` succeeds.
- `cd deployments/cdk && npm test` passes.

## Out of scope

- Any schema change.
- Touching `pkg/awssupport/dynamoutils/CreateOrUpdateTable` (the generic helper that is kept).
- Introducing new tests beyond what is necessary to cover the new helper.

## Depends on

- Story 02 (needs `schema.json` to exist).

## References

- `specs/2026-09-retire-appdynamodb-schema/spec.md`
- `specs/2026-09-retire-appdynamodb-schema/issues/02-extract-schema-to-json.md`
- `pkg/awssupport/appdynamodb/create_update_table.go` (to be deleted)
- `pkg/awssupport/dynamotestutils/context.go`
- `tools/dphotoops/migrator/transformattion_indexes.go`

## Comments

<!-- Record here the decision on tools/dphotoops/migrator/ (delete vs keep-and-adapt) with evidence, before merging. -->
