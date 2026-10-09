# 14 — `POST …/covers/refresh` endpoint

Status: done
Phase: 3
Layer: api — `api/lambdas/randomize-album-covers` (new lambda)
Depends on: 13

## Description

Expose the Phase 3 randomize operation over REST. Owner-edit permission; returns the new covers in the
response body so the UI can refresh the album card without re-listing all albums.

See `../design.md` → §Phase 3 REST contract.

## Acceptance criteria

- New lambda folder `api/lambdas/randomize-album-covers/` with `main.go`, `main_test.go`, following
  the layout and conventions of the other lambdas (see `api/lambdas/list-albums/`).
- Route: `POST /api/v1/owners/{owner}/albums/{folderName}/covers/refresh`.
- Path parameters: `owner`, `folderName`. Both URL-decoded.
- Body: ignored (empty expected).
- Authorisation: the caller must have owner-edit permission on the album — same rule as
  `rename-album` / `amend-album-dates`. Reuse the existing authorisation helper from
  `api/lambdas/common/` (follow the pattern used by the rename lambda).
- Behaviour: calls `RandomizeAlbumCoversCase.Randomize(ctx, albumId)`. On success, returns `200 OK`
  with body:
  ```json
  { "covers": [ { "mediaId": "...", "filename": "...", "origin": "RANDOM" } ] }
  ```
  Reuse the `CoverDTO` shape from `api/lambdas/list-albums/main.go` (either import or copy; follow
  the project convention — `expose-rest-api` skill).
- Error mapping:
  - Album not found → `404`.
  - Unauthorised → `403`.
  - Any other error → `500`.
- CDK: register the lambda in `deployments/cdk/` alongside the other catalog lambdas. The resource
  naming follows the CDK conventions (`cdk` skill).
- Tests:
  - `main_test.go`: table-driven cases covering success (owner), forbidden (visitor), not-found,
    randomize port error.
  - CDK synth test extended (if the existing suite covers lambda registration) to assert the new
    lambda is wired.

## Out of scope

- `SetCovers` / per-cover endpoints — Phase 4+.
- UI — issue 15.

## References

- `../spec.md` (Phase 3; User journeys → Re-randomise covers).
- `../design.md` (REST contract).
- Load skills: `go`, `expose-rest-api`, `cdk`.
