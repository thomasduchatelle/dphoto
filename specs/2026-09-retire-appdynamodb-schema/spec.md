# Retire `appdynamodb.CreateTableIfNecessary`

## Intent

Remove the legacy Go-side DynamoDB schema management (`pkg/awssupport/appdynamodb/create_update_table.go`) whose eviction was announced in-code but never completed. Establish a single source of truth for the catalog table schema so tests, Localstack bootstrap, and CDK stay aligned without duplicating attribute and index definitions.

The file itself carries a "notice of eviction" comment: production schema management belongs in CDK, but six callers still depend on the Go helper (Localstack CLI bootstrap, four test files, one lambda, one ops tool). The lambda is orphaned; the rest need a replacement path before the helper can be deleted.

## Current situation

- `CreateTableIfNecessary` defines all attributes and the 5 GSIs (`AlbumIndex`, `ReverseLocationIndex`, `ReverseGrantIndex`, `RefreshTokenExpiration`, `AlbumViewByAlbumIndex`) in Go.
- `deployments/cdk/lib/catalog/catalog-store-construct.ts` defines the same 5 GSIs in TypeScript.
- The two definitions are maintained by hand and can drift.
- One caller (`api/lambdas/sys-dynamodb-upgrade/`) is not wired to any CDK/Make target — it is dead code.

## Outcome

- The catalog table schema (attributes, key schema, GSIs) is defined once, in a shared JSON file.
- Both Go (tests + Localstack bootstrap) and CDK read that file to build their respective DynamoDB objects.
- `pkg/awssupport/appdynamodb/` is deleted; the orphaned lambda is deleted.
- The rationale for each GSI (currently inline Go comments) moves to `DATA_MODEL.md`.

## Scope

**In scope:**
- Deleting `api/lambdas/sys-dynamodb-upgrade/`.
- Extracting the schema to a shared JSON file consumed by both Go and CDK.
- Migrating the remaining Go callers to a test-scoped helper that reads the JSON.
- Deciding the fate of `tools/dphotoops/migrator/transformattion_indexes.go` (delete or drop its schema-creation step).
- Moving GSI rationale into `DATA_MODEL.md`.

**Out of scope:**
- Any change to the schema itself (no new indexes, no attribute changes).
- Changes to `pkg/awssupport/dynamoutils/CreateOrUpdateTable` (the generic helper — kept as-is).
- Refactoring `CatalogStoreConstruct` beyond what's needed to read the JSON.

## Decision record

Considered options for eliminating the duplication:

- **A. Test-only helper, keep duplication.** Cheapest but does not solve the drift risk.
- **B. Go tests read `cdk synth` output.** Eliminates duplication but adds a Node toolchain dependency and 3–8s cold-start cost to every Go test run. Rejected: too heavy for a pet project.
- **C. Tests against a pre-built DynamoDB container.** Too much operational overhead for the value.
- **D. Single source of truth (JSON), both sides read it.** Zero test-runtime cost, zero new toolchain dependencies (Go `encoding/json` + native `require()` in TypeScript), eliminates drift. **Chosen.**

Schema shape: **SDK-native**, mirroring `dynamodb.CreateTableInput` in Go. Translating to the CDK `Table` construct is the smaller of the two translation directions.

Comments/rationale: kept in `DATA_MODEL.md` rather than JSONC or `_comment` fields, to keep the machine-oriented schema file boring and to consolidate the "why" in the existing data model doc.

## References

- `pkg/awssupport/appdynamodb/create_update_table.go` — the file being retired.
- `deployments/cdk/lib/catalog/catalog-store-construct.ts` — CDK construct that will consume the JSON.
- `DATA_MODEL.md` — destination for GSI rationale.
- `specs/archived/todo_cognito.md:8` — pre-existing TODO flagging the `sys-dynamodb-upgrade` lambda for deletion.
