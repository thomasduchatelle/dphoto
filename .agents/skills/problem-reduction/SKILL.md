---
name: problem-reduction
description: Break down a problem that cannot fit into a single development session (or review) into smaller chunks. Use it when refining features and stories with the user.
---

# Problem reduction

Problem reduction is the most important task of the software development: don't be a scribe, **be a peer by bringing your ideas and pushing back if something doesn't sound right**.

The objective of the problem reduction is to:

* explore and clarify an idea into unambiguous requirements.
* split a work that wouldn't fit into a single coding session into several that can each be implemented, reviewed, and merged individually.

## Distinction of the WHAT and the HOW

**The requirement and the technical solution must never be mixed together.** They must be documented independently:

* **WHAT** are the requirement and are presented from a consumer point of view. Use the `issue-tracker` to capture them.
* **HOW** are the design decisions that are made during the refinement (DB indexing, responsibility between layers, UI components structure, ...). Capture them in dedicated design document, or ADR if they are hard to revert or if they affect future work.

**Focus on the WHAT and only engage on the HOW when a decision spans several issues. Leave implementation details to the coding agents, do not plan the work, do not list the files to update.**


## Stories or Tasks?

Reducing is done through two axes:

* **stories - vertical slices** - this is the primary axis to use. Find how an idea can be divided into thin slices: stories. Each of them defines an end-to-end behaviour visible and beneficial for the consumer of the application in this repository. It crosses any layer necessary: UI, CLI, API, DB, ... It is preferable when stories can be implemented independently, but it is expected they have dependency between each others or that they need to be implemented in sequence to avoid conflicts.
* **tasks - horizontal layers** - this is only used when the stories are already as small as possible, but the work required wouldn't fit into a single coding session. A task will cover a single component of the application (i.e. the DB repository, or the exposure of the API, or the domain logic, or the UI).

Find the right balance between vertical slices (better for coding agents working end-to-end) and small, independently reviewable units (better for the user). Patterns and examples:

* **walking skeleton** - an initial story with the bare minimum to wire the journey end to end; functional requirements are added in subsequent stories (likely in parallel, one per use case). Use it when the main concern is connecting different domains or layers.
* **domain driven** - an initial story focused on domain logic to resolve uncertainty about the required interface; subsequent PRs implement the persistence adapter, expose the feature, wire observers, etc. Use it when the domain logic is complex and likely to require several iterations with the user.
* **contract-first** - the interface is defined first in the design document, then each story implements one side of the contract (consumer/supplier). The contract can be within the same language (Golang, Typescript, or Java interfaces) or over the network (REST, JSON). Use when the implementation is straightforward and can proceed independently per technology (web / backend / infra) once the interface is defined.

**Prioritise parallelisation of coding tasks**. A single issue is preferable to a sequence of several ones as long as it fits into a coding session.

## Issue template

Each issue - whether they are an epic, a story, or a task - follows the same template:

* **a title**
* **a short description** - self-contained but restricted to what's absolutely required for the task to be implemented, link relevant design or ADR documents, never refer to another issue or link the parent spec.
* **dependencies** - the numbers of the issues that must be implemented before.
* **acceptance criteria** - make the objective unambiguous, use bullet points.
* **an out-of-scope list** - prevent any overlap with other issues, close debates already discussed.

Keep every part concise and short, do not extrapolate.

## Common pitfalls

* **no history in the issues** - when amending an issue, the resulting issue must be as it would be if it was written properly on the first time: do not mention the reason of the change or what it was before, it will confuse the reader. The user might explain its decision and idea so you can adapt or push back ; these explanations are adressed to you, they are not to be recorded. 
* **keep it concise** - if the user can describe a requirement in two sentences, you shall keep it 2 sentences as well. Only extend it if it brings important information that were discussed before with the user. Do not add things that you're thinking about: if you can come up with the solution, so can the agent. If that's non-obvious, you should debate it with the user first.
* **fragmented issues** - there is no benefit to splitting an issue that already fits into a coding session. If issues become too small for no benefit (parallelisation), raise it to the user.
* **sequential tasks** - several tasks that need to be implemented in a certain order and without possible parallelisation is a sign of fragmented issues.
* **no planning** - any file listing or code sample is likely a plan which shouldn't have been done at refinement time. This is a sign of fragmented issues.
