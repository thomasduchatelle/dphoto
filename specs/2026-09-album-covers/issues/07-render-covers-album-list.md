# 07 — Render album covers on the album list

Status: ready
Phase: 2
Layer: web — `web-nextjs`
Depends on: 04

## Description

Show each album's real covers on its card in the album list, replacing the placeholder `thumbnails` field.
Covers are read-only in this phase (no picking or re-randomising yet).

See `../design.md` (REST contract) and `spec.md`.

## Acceptance criteria

- The `Album` type carries the covers from the API (`mediaId`, `filename`, `origin`); the placeholder
  `thumbnails?: string[]` is removed/replaced.
- The album card renders up to 4 covers, building each image URL from the album owner + `mediaId` +
  `filename` through the existing image loader.
- Albums with fewer than 4 (including 0) covers render gracefully (empty slots as today).
- Storybook/visual tests cover 0, partial, and full cover sets.

## Out of scope

- Any cover editing UI (re-randomise/pick/unpick — Phases 3–5).
- Changes to the image resize/format pipeline.

## References

- `../spec.md`, `../design.md`
- Load skills: `nextjs`, `ui-components`.
