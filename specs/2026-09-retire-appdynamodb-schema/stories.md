# Stories

Three stories, sequenced so risk grows gradually. Stories 01 and 02 are independent; Story 03 depends on Story 02.

```
Story 01 (delete dead lambda) ──┐  independent
                                 │
Story 02 (extract to schema.json)┘  independent
                                 │
Story 03 (retire appdynamodb) ───► depends on 02
```

## 01 — Delete orphaned `sys-dynamodb-upgrade` lambda

Remove the lambda that is not wired to any CDK/SST stack or Makefile target. It is the only remaining caller of `CreateTableIfNecessary` with `localDynamodb=false` and its deletion is already flagged in `specs/archived/todo_cognito.md:8`.

**Acceptance:**
- `api/lambdas/sys-dynamodb-upgrade/` is deleted.
- `grep -r sys-dynamodb-upgrade` returns no matches across the repo.
- `make build-api` and `go test ./...` pass.
- The corresponding line in `specs/archived/todo_cognito.md` is removed.

**Out of scope:** the `dphotoops` migrator; any change to `appdynamodb`.

## 02 — Extract catalog table schema to shared JSON

Move the schema definition into a single JSON file consumed by both `CreateTableIfNecessary` (Go) and `CatalogStoreConstruct` (CDK). This is the substantive refactor: after this story, drift between the two definitions becomes structurally impossible.

**Acceptance:**
- New file `pkg/awssupport/appdynamodb/schema.json` describes attributes, key schema, and all 5 GSIs in SDK-native shape.
- `CreateTableIfNecessary` builds its `CreateTableInput` from the embedded JSON (via `//go:embed`); no schema literals remain in Go code.
- `CatalogStoreConstruct` builds its `Table` and GSIs from the same JSON file.
- Existing `CatalogTableIndexes` constant in `catalog-store-construct.ts` is either derived from the JSON or removed.
- Rationale for each GSI (currently inline Go comments) is moved to `DATA_MODEL.md`.
- `go test ./...`, `cd api/lambdas && go test ./...`, and `cd deployments/cdk && npm test` all pass.
- `cd deployments/cdk && npm run synth:test` output diff against `main` is either empty or explained.

**Out of scope:** removing `CreateTableIfNecessary`; migrating remaining callers; migrator changes.

## 03 — Retire `pkg/awssupport/appdynamodb/` and resolve the migrator

Delete the legacy Go schema helper once all callers use the JSON-backed replacement, and make an explicit call on the `dphotoops` migrator.

**Acceptance:**
- Decision recorded (in this issue's Comments section) on `tools/dphotoops/migrator/`:
  - If unused: delete `tools/dphotoops/migrator/transformattion_indexes.go` (and the tool if empty).
  - If still used: remove its `CreateTableIfNecessary` call and document `cdk deploy` as a prerequisite for running migrations.
- The remaining Go callers switch to a new helper `dynamotestutils.CreateLocalTable(ctx, client, tableName)` (or equivalent) that reads `schema.json`:
  - `pkg/awssupport/dynamotestutils/context.go`
  - `pkg/catalogadapters/catalogdynamo/repository_album_test.go`
  - `pkg/catalogadapters/catalogdynamo/repository_media_crud_test.go`
  - `pkg/backup_acceptance_test.go`
  - `cmd/dphoto/bootstrap/bootstrap_catalog.go` (Localstack branch)
- `pkg/awssupport/appdynamodb/create_update_table.go` is deleted; empty package directory removed.
- `grep -r CreateTableIfNecessary` returns matches only in `pkg/awssupport/dynamoutils/` (the unrelated generic helper).
- `go test ./...`, `make build-api`, `cd deployments/cdk && npm test` all pass.

**Out of scope:** any schema change; touching `dynamoutils.CreateOrUpdateTable`.

**Depends on:** Story 02.
