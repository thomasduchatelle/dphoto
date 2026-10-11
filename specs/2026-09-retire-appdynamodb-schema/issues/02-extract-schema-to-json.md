# 02 — Extract catalog table schema to shared JSON

Status: ready

## Description

The catalog DynamoDB table schema (attributes, key schema, and 5 GSIs) is currently defined twice: once in Go at `pkg/awssupport/appdynamodb/create_update_table.go` and once in TypeScript at `deployments/cdk/lib/catalog/catalog-store-construct.ts`. Both definitions are maintained by hand and can drift.

Extract the schema into a single JSON file at `pkg/awssupport/appdynamodb/schema.json` and rewrite both sides to read from it. The JSON shape should be SDK-native (mirroring `dynamodb.CreateTableInput`) — this keeps the Go side a direct unmarshal and pushes the small translation work to the TypeScript/CDK side.

Rationale currently living as inline Go comments (`// from 'archivedynamo' extension`, `// from 'acl' extension`, `// from 'catalogviews' extension`) must be moved to `DATA_MODEL.md` so it is not lost when the Go file is deleted in Story 03.

After this story, drift between Go and CDK schema definitions is structurally impossible.

## Acceptance criteria

- New file `pkg/awssupport/appdynamodb/schema.json` describes:
  - All 9 attribute definitions currently in the Go code.
  - The primary key schema (`PK` hash, `SK` range).
  - All 5 GSIs: `AlbumIndex`, `ReverseLocationIndex`, `ReverseGrantIndex`, `RefreshTokenExpiration`, `AlbumViewByAlbumIndex`, each with its key schema and projection settings.
- `CreateTableIfNecessary` in `pkg/awssupport/appdynamodb/create_update_table.go` builds its `CreateTableInput` from the embedded JSON via `//go:embed`. No attribute names, index names, or key schema literals remain hard-coded in Go.
- The `tableVersion` constant and `localDynamodb` provisioned-throughput behaviour are preserved.
- `CatalogStoreConstruct` in `deployments/cdk/lib/catalog/catalog-store-construct.ts` builds the `Table` and its GSIs from the same JSON file (via `require()` or `readFileSync` + `JSON.parse`).
- The `CatalogTableIndexes` constant currently exported from `catalog-store-construct.ts` is either derived from the JSON or removed if no external consumer depends on the exact array.
- Each GSI has a documented rationale in `DATA_MODEL.md` (why it exists, what queries it serves). Existing inline Go comments are the starting point.
- `go test ./...` passes.
- `cd api/lambdas && go test ./...` passes.
- `cd deployments/cdk && npm test` passes.
- `cd deployments/cdk && npm run synth:test` produces output that is either byte-identical to `main` or whose diff is documented in the PR description.

## Out of scope

- Deleting `pkg/awssupport/appdynamodb/create_update_table.go` — Story 03.
- Migrating the other Go callers to a new helper — Story 03.
- Any change to the schema itself (no added indexes, no attribute changes, no projection changes).
- Changes to `pkg/awssupport/dynamoutils/CreateOrUpdateTable`.

## References

- `specs/2026-09-retire-appdynamodb-schema/spec.md`
- `pkg/awssupport/appdynamodb/create_update_table.go`
- `deployments/cdk/lib/catalog/catalog-store-construct.ts`
- `DATA_MODEL.md`
