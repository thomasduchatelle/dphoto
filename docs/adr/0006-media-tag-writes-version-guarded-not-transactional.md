# ADR-0006: Media tag writes are version-guarded, not transactional

- Status: accepted
- Date: 2026-11-26
- Scope: `pkg/catalogadapters/catalogdynamo` (tag persistence) + `DATA_MODEL.md`

## Context

A media's tags are modelled as two representations that must agree:

- the **source of truth**: a `Tags` String Set on the media's `#METADATA` row, loaded alongside
  the media when an album is read;
- a **derived index**: one `TAG#{tag}` sibling item per (media, tag) feeding the `MediaByTagIndex`
  GSI, because a GSI cannot be built on a set attribute (DynamoDB forbids set-typed index keys) —
  the sibling is the only way to answer "all medias of a tag".

`AddTag`/`RemoveTag` therefore mutate two items. The hazard is concurrency: the favourite star,
the paint-tagging click, and API retries all issue independent, genuinely concurrent writes on the
same (media, tag). The forbidden end-state is **the tag present in the source of truth but its
sibling absent** — a real tag the user applied becomes invisible on the Tags page, with no way for
the read path to rediscover it (the GSI simply never returns the media).

## Decision

Tag writes are kept **lock-free and non-transactional**, made safe with a monotonic version used as
a compare-and-set guard on the sibling:

- Every set mutation on `#METADATA` is a single atomic `UpdateItem` that both merges the tag
  (`ADD`/`DELETE Tags :tag`) **and** bumps a counter (`ADD TagVersion :1`), returning the new value
  via `ReturnValues: UPDATED_NEW`. Because DynamoDB serialises updates to a single item, `TagVersion`
  is monotonic: the set-mutation that logically happened later always carries the higher version.
- That version is stamped onto the `TAG#` sibling. The sibling reconcile is conditional —
  `putTagSibling` and `deleteTagSibling` only win when
  `attribute_not_exists(TagVersion) OR TagVersion < :v`. A stale concurrent operation can therefore
  never clobber a sibling a newer operation owns; a lost race is a swallowed
  `ConditionalCheckFailedException` (a no-op success).
- Ownership is enforced without a prior read: the set update carries `attribute_exists(PK)`, so a
  media the owner does not have is rejected as `MediaNotFoundError`.
- `FindMediasByTag` treats the `#METADATA` set as the final authority: it drops (and WARN-logs) any
  GSI hit whose `#METADATA` no longer carries the tag, so a residual orphan sibling is never shown.

## Considered options

1. **Read-modify-write (`GetItem` then `SET Tags = <whole list>`).** The first implementation.
   Rejected: two concurrent writers both read the same base list and the last `SET` wins, silently
   dropping the other's tag — a lost update on the source of truth itself.
2. **Atomic set `ADD`/`DELETE` + fixed write ordering ("index is always a superset of the truth").**
   Fixed concurrent *adds* (the set merges), but a concurrent **add + remove of the same tag** could
   still interleave so the set ended with the tag while the sibling was deleted — the forbidden
   state. The "superset" invariant did not actually hold across mixed verbs.
3. **`TransactWriteItems` grouping the set update and the sibling write.** Fully correct for both
   concurrency *and* a mid-operation crash. Rejected as the default: it doubles the write cost of
   every tag toggle (2 WRU per item), introduces the repo's first transaction and a
   `TransactionCanceledException` retry path, and buys protection only for a failure mode we judged
   rare (see consequences). At the project's scale the dollar difference is negligible, so the
   decision is driven by simplicity, not cost.
4. **Version-guarded CAS (chosen).** Closes the concurrency hole — including add/remove of the same
   tag — with no transaction, one extra number attribute, and single-item serialisation providing the
   monotonicity for free.

## Consequences

- **Accepted residual: a mid-operation crash can strand a tag.** `AddTag` updates the set first (to
  mint the version), then writes the sibling. If the process dies *between* those two calls, the tag
  is in the source of truth with no sibling until the client retries the add (`putTagSibling`
  overwrites the missing/older sibling on the next attempt). This is a crash, not a concurrency race;
  it is far rarer than the concurrency path and self-corrects on the client's retry. We ship it over
  a transaction on purpose.
- **Residual orphan siblings are benign and self-healing.** A partial failure in the other direction
  (sibling present, truth without the tag) is dropped on read by `FindMediasByTag` and overwritten by
  the next real write of that (media, tag). The read-path truth-check is load-bearing: removing it
  would surface orphans.
- `TagVersion` lives on `#METADATA`, shared across all of a media's tags. This is deliberate — a
  single per-media serialisation point is what makes the version monotonic — and is harmless to the
  album read path (`MediaRecord` has no `TagVersion` field, so the attribute is ignored on unmarshal).
- The `MediaByTagIndex` projects `KEYS_ONLY`: the reader derives the media id from the SK and
  batch-reads `#METADATA` for the real attributes, so the GSI carries no attributes nothing reads.
