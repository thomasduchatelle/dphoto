# 05 — List the owner's distinct tags

Status: ready

## Description

Provide a way to enumerate the distinct tags an owner has used, grouped by namespace. The Tags page switcher needs this, and it is not answerable by the `{owner}#{tag}` GSI (which finds medias *of* a tag, not the *set* of tags).

Decide and implement the enumeration strategy per `specs/2026-11-media-tagging/design.md` ("Distinct tags (decision owed)" — registry vs. scan): the likely approach is a per-owner tag-registry item maintained on tag add/remove, so enumeration is a single cheap read. Expose it as the `GET /api/v1/owners/{owner}/tags` endpoint pinned in the design doc's "Client/server interface".

Load the `go`, `expose-rest-api`, and `cdk` skills as relevant.

## Dependencies

02.1.

## Acceptance criteria

- An endpoint returns the owner's distinct tags grouped by namespace.
- The result stays consistent as tags are added and removed (no stale or orphaned tags).
- Enumeration does not require a table scan at read time.
- `DATA_MODEL.md` updated if a new item/index is introduced.
- Tests follow the Go table-driven style.

## Out of scope

- Tag counts.
- Any UI (consumed by the switcher story).
