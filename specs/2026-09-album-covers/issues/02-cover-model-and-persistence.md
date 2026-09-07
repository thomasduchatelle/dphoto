# 02 — Cover model & single-record persistence + random completion

Status: ready
Phase: 2
Layer: catalog domain — `pkg/catalog` + `pkg/catalogadapters/catalogdynamo`
Depends on: —

## Description

Introduce the **cover** concept in the catalog and persist an album's cover set as a single DynamoDB record.
Provide the operation that fills empty slots with uniform-random photos. This phase covers automatic
completion only — pick, unpick, and re-randomise come in later phases.

See `../design.md` (Cover canonical record, Completion operation) and the catalog language for `Cover` /
`CoverOrigin`.

## Acceptance criteria

- Domain model: `Cover` = `MediaId` + `Filename` + `CoverOrigin`; `CoverOrigin` ∈ `RANDOM`,
  `CHERRY_PICKED`. An album has 0–4 covers.
- Persistence: a single item per album, PK `{OWNER}#ALBUM`, SK `ALBUM#{FOLDER_NAME}#COVERS`, holding the
  whole ordered set. Reads and writes operate on the full set; the max of 4 is enforced on write.
- Completion operation with two entry shapes:
  - `CompleteCoversFromCandidates(albumId, candidates)` — fills empty slots (up to 4) from the supplied
    medias, no extra query.
  - `CompleteCovers(albumId)` — queries the album's `IMAGE` medias, then completes.
  - Both select **uniformly at random** among eligible `IMAGE` medias not already covers, tag them
    `RANDOM`, and are a no-op when the set is already full.
- Only `MediaType IMAGE` is eligible; videos/other are never selected.
- `DATA_MODEL.md` documents the new record.
- Tests follow the Go testing strategy (table-driven, in-memory adapter for domain tests).

## Out of scope

- Cherry-pick, unpick, re-randomise operations (Phases 3–5).
- Denormalising covers into the view (issue 03).
- API and UI.
- Triggering completion from backup (issue 05) or backfill (issue 06).

## References

- `../spec.md`, `../design.md`
- Load skills: `go`, `architecture`.
