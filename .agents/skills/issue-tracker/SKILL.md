---
name: issue-tracker
description: Administer specs and issues for this repo. Use when locating the spec/issue for a task, publishing a new feature spec, breaking a feature into issues, recording issue status, or archiving finished work.
---

# Issue Tracker (local markdown)

Specs and issues live as markdown files under `specs/`. There is no external tracker. This skill is the single source of truth for how it works.

## Layout

```
specs/
├── <feature-slug>/           # one directory per feature; dated prefix, e.g. 2026-02-<slug>
│   ├── spec.md               # MANDATORY: the feature spec (the "what")
│   ├── design.md             # recommended: upfront technical direction (see below)
│   ├── stories.md            # recommended: draft of all issues before splitting them out
│   └── issues/               # one file per issue: NN-<slug>.md, numbered from 01
└── archived/                 # finished features and superseded material
```

## The documents

- **`spec.md` (mandatory)** — describes the feature: intent, user journeys, scope, and what's out of scope. The "what", not the "how".
- **`design.md` (recommended)** — technical direction that spans several issues and must be settled upfront: context boundaries (who does what), client/server interface (REST contracts), and data model (decided early for performance and backward-compatibility). Durable, repo-wide decisions graduate to an ADR under `docs/adr/` and are *linked* from the issues, not duplicated.
- **`stories.md` (recommended)** — a draft of every issue (slug + acceptance criteria) so the user and agent can iterate fast before committing to one file per issue. Once agreed, each entry becomes an `issues/NN-<slug>.md`.
- **`issues/NN-<slug>.md`** — one implementable ticket. Self-sufficient on the "what": title, short description, acceptance criteria, out-of-scope. Links to `spec.md` / `design.md` / ADRs for the "how" rather than repeating them. Never a single combined tickets file.

## Status

Each issue carries a `Status:` line near the top. Only three values:

| Status   | Meaning                               |
|----------|---------------------------------------|
| `ready`  | Specified and ready to be implemented |
| `done`   | Implemented and merged                |
| `wontdo` | Will not be actioned                  |

There is no separate status or sprint file — the issue file is the source of truth. Comments and history append at the bottom under a `## Comments` heading.

## Lifecycle

1. **Publish a spec** — create `specs/<feature-slug>/` with `spec.md` (add `design.md` when upfront technical direction is needed).
2. **Draft the issues** — capture them in `stories.md`, iterate with the user, then write one `issues/NN-<slug>.md` per agreed story (numbered from `01`). Don't pre-generate half-baked issues.
3. **Work an issue** — the implementing agent sets `Status: done` when the work is merged (`wontdo` if dropped).
4. **Close the feature** — when all its issues are `done`, move `specs/<feature-slug>/` to `specs/archived/<feature-slug>/`.

## Rules

- Use `git mv` for every move so history follows the file.
- Fix cross-references when you move a file; leave historical implementation reports untouched even if they name paths that no longer exist.
- Don't invent a parallel status tracker — status is the `Status:` line.

For planning an effort too large to hold in one session (a map with decision tickets), use the `wayfinder` skill.
