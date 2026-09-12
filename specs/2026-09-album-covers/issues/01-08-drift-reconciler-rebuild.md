# 01-08 — Drift reconciler rebuilds display fields

Status: ready
Phase: 1
Layer: catalog domain — `pkg/catalogviews` + `pkg/pkgfactory`
Depends on: 01-01, 01-02

## Description

Teach the drift reconciler to project the full album summary (including `Name/Start/End`) so it can:

1. detect stale display fields against the canonical album records, and
2. backfill display fields on rows written before this feature (legacy rows created when the record
   was count-only).

The reconciliation flow is unchanged in shape (`OwnerDriftReconciler.Reconcile` → `AlbumReCounter`
→ observer/synchronizer). What changes is the payload: `AlbumSummary` now travels with display
fields, so the drift detector compares them and the synchronizer writes them.

The post-deploy operational step for this feature is: run the CLI's owner-drift reconciliation for
every owner once, in non-dry mode.

## Acceptance criteria

- `AlbumReCounter.ReCountMedias` also fetches display fields for its target albums (via
  `FindAlbumsByIdsPort` — a new dependency of `AlbumReCounter`) and populates them into every
  `AlbumSummaryForUsers` it emits.
- `DriftDetector` compares `Name/Start/End` in addition to `MediaCount` and reports an
  `OverrideDrift` when any of them differs.
- `DriftSynchronizerObserver` writes the full record (via `Repository.PutSummaries`) so a rebuild
  covers count and display fields in a single write.
- Legacy behaviour: a row that exists with `Count` and empty display fields is treated as a drift and
  overridden with the canonical values.
- Tests in `pkg/catalogviews/view_drift_control_test.go`:
  - Existing cases updated so seeded `UserAlbumSummary`s carry `Name/Start/End`.
  - `TestOwnerDriftReconciler/it_should_rebuild_display_fields_from_the_canonical_album` — new;
    seeds a row with wrong `Name/Start/End`; asserts the reconciler writes the canonical values.
  - `TestOwnerDriftReconciler/it_should_backfill_missing_display_fields_on_a_legacy_row` — new;
    seeds a row with only `Count`; asserts the reconciler writes `Name/Start/End`.
- Factory `OwnerDriftReconciler` in `pkg/pkgfactory/factory_catalog_view.go` updated to pass
  `AlbumQueries` (or another `FindAlbumsByIdsPort`) into `AlbumReCounter`.
- Ticket description contains the exact command / entry point to run post-deploy (CLI subcommand
  and/or one-liner using the factory) so the operational step is unambiguous.

## Out of scope

- Any change to the read path or the event methods.
- Change of drift-detection strategy (still row-by-row, per user).
- Automatic scheduling of the reconciliation.

## References

- `../spec.md`, `../design.md`
- Parent ticket `01-album-list-read-model.md` (interface contract).
- Load skills: `go`, `architecture`.
