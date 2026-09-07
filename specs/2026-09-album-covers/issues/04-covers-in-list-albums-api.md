# 04 — Expose covers on `GET /albums`

Status: ready
Phase: 2
Layer: api — `api/lambdas/list-albums`
Depends on: 03

## Description

Return each album's covers in the album-list response. Breaking changes to the API are accepted.

See `../design.md` (REST contract).

## Acceptance criteria

- `AlbumDTO` gains a `covers` array; each entry is `{ "mediaId": string, "filename": string, "origin":
  "RANDOM" | "CHERRY_PICKED" }`.
- `covers` is populated from `VisibleAlbum.Covers`; it is `[]` when the album has none (0–4 entries).
- No image URL is built or returned server-side — only `mediaId` + `filename` + `origin`.
- Tests follow the Go testing strategy for the handler.

## Out of scope

- Cover mutation endpoints (Phases 3–5).
- Frontend rendering (issue 07).

## References

- `../spec.md`, `../design.md`
- Load skills: `go`, `architecture`.
