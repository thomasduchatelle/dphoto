---
name: issue-tracker
description: Administer specs, epics, stories, and issues for this repo. Use when writing or publishing a new feature spec, breaking a feature into issues, or archiving finished work.
---

# Issue Tracker on local Markdown

Specs and issues live as Markdown files under `specs/`. There is no external tracker.

## Layout

```
specs/
├── <feature-slug>/           # one directory per feature; dated prefix, e.g. 2026-02-<slug>
│   ├── spec.md               # MANDATORY: the feature spec (the "what")
│   ├── design.md             # recommended: upfront technical direction (see below)
│   ├── draft.md              # optional, temporary: draft of all issues to converse with the user before splitting them out
│   └── issues/               # one file per issue: NN-<slug>.md, numbered from 01 ; exceptionally, subtasks can be used with are NN.MM-<slug>.md
└── archived/                 # finished features and superseded material
```

## The documents

- **`spec.md` (mandatory)** — describes the feature: intent, user journeys, scope, and what's out of scope. The "what", not the "how".
- **`design.md` (recommended)** — technical direction that spans several issues and must be settled upfront: context boundaries (who does what), client/server interface (REST contracts), and data model (decided early for performance and backward-compatibility). Durable, repo-wide decisions graduate to an ADR under `docs/adr/` and are *linked* from the issues, not duplicated.
- **`draft.md` (optional, and temporary)** — use it when refining the feature to let the user edit it directly, delete it once the `issues/NN-<slug>.md` got generated unless some tasks are missing.
- **`issues/NN-<slug>.md`** — one implementable ticket. Self-sufficient on the "what": title, short description, acceptance criteria, out-of-scope (i.e. implementation doesn't need to read `spec.md`). Contains relevant links to `design.md` and ADRs rather than repeating them.

## Status

Each issue carries a `Status:` line near the top. Only three values:

| Status   | Meaning                                                 |
|----------|---------------------------------------------------------|
| `ready`  | Specified and ready to be implemented (default status). |
| `done`   | Implemented on the current branch.                      |
| `wontdo` | Will not be actioned.                                   |

## Lifecycle

1. **Publish a spec** or **capture an idea** — create `specs/<feature-slug>/` with `spec.md` (add `design.md` only when technical direction is given with the idea).
2. **Draft the issues** — capture them in `draft.md`, iterate with the user, then write one `issues/NN-<slug>.md` per agreed story (numbered from `01`). Don't pre-generate half-baked issues.
3. **Work an issue** — the implementing agent sets `Status: done` when the work is done (`wontdo` if dropped).
4. **Close the feature** — when all its issues are `done`, move `specs/<feature-slug>/` to `specs/archived/<feature-slug>/`. Use `git mv` to archive.

## Definitions

Terms to use with the user are:

* **story**: this is an update, visible by a user of this repository application. A large one can be captured as a `specs/<feature-slug>/spec.md` in this tracker, or as an issue of a larger feature `specs/<feature-slug>/issues/NN-<task-slug>.md`.
* **task**: this is the implementation of a module, a layer, or a class of this application. This is usually captured as an issue of a feature `specs/<feature-slug>/issues/NN-<task-slug>.md`.
* **bug**: this is fixing a behaviour considered incorrect. It will be captured as a subissue if it's part of a feature ; or as a feature itself if it's one out of the blue.
* **epic**: this is a milestone covering several stories. They are not formally represented in this issue tracker: they will likely be initially represented as a feature, but then will be replaced by several stories without formal relationship.
* **issue**: this is a generic term that includes stories, tasks, epics, ...

