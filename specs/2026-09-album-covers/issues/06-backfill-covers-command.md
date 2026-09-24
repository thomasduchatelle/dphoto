# 06 — Backfill covers for all owners/albums

Status: ready
Phase: 2
Layer: CLI — `cmd/dphotops`
Depends on: 02

## Description

A one-off, administrator-run command that gives existing albums covers, since backup only completes albums
that receive new medias. Add a subcommand under `cmd/dphotops` (the admin binary that reuses `dphoto`'s
config).

See `../design.md` (Completion operation) and `spec.md` (Backfill existing albums).

## Acceptance criteria

- A CLI entry point iterates every album of every owner and completes empty cover sets using the querying
  variant (`CompleteCovers`, which lists the album's `IMAGE` medias).
- Idempotent and safe to re-run: albums with a full set are untouched; `CHERRY_PICKED` covers are never
  altered.
- Runs across all owners in one invocation.
- Tests follow the Go testing strategy.

## Out of scope

- Any web/API surface (this is CLI only).
- Replacing or re-randomising existing covers.

## References

- `../spec.md`, `../design.md`
- Load skills: `go`, `architecture`.
