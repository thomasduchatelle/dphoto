# 05 — Backup completes covers for touched albums

Status: ready
Phase: 2
Layer: backup domain — `pkg/backup` + `cmd/dphoto`
Depends on: 02

## Description

When a backup run adds new medias, complete the covers of the affected albums so newly-populated albums show
covers without waiting for the admin backfill. Keep the backup hot path cheap.

See `../design.md` (Completion operation) and `spec.md` (Backup fills missing covers).

## Acceptance criteria

- At the **end of a backup batch**, for each album that received at least one new media, call the completion
  operation using the **just-inserted `IMAGE` medias as candidates** (`CompleteCoversFromCandidates`) — no
  extra full-album query per batch.
- Albums whose cover set is already full are left untouched (no-op).
- Albums that received no additions are not processed.
- If a touched album's additions contain no images, its empty slots stay empty (filled later by
  re-randomise or backfill).
- Tests follow the Go testing strategy for the backup flow.

## Out of scope

- Querying the whole album to complete (that is the backfill path, issue 06).
- Re-randomising or replacing existing covers.

## References

- `../spec.md`, `../design.md`
- Load skills: `go`, `architecture`.
