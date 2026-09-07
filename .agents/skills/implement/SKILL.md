---
name: implement
description: "Implement a piece of work based on a spec or set of tickets."
disable-model-invocation: true
---

Implement the work described by the user in the spec or tickets.

## Planning

During planning you must:

- take advantage of the planning to ask clarification questions
- propose different - non-trivial - design solutions and ask the user to choose one
- present your plan which must start by:
  a. the main design decisions
  b. your test cases (names and headline, not code)
  c. followed by the usual plan (no code)

## Implementation

The tests must be robust against refactoring: they demonstrate the story is implemented (the acceptance criteria), they should not aim for a 100% coverage by testing the implementation.

Run typechecking regularly, single test files regularly, and the full test suite once at the end.

Progress the issue(s) you're working on as `done`.

If working on a issue, commit your work to the current branch using the pattern: `<context>[/<layer>] - <summary> + <body>`. Never amend or force existing commits.

* **context** is `archive`, `catalog`, `ci`, ... Use `llm` when working on skills or agent documentation, and `proj` when working on the issue-tracker.
* **layer** is added if the change only impct a single layer: `web`, `cli`, `api`, ... 
* add in the body `+next` if you consider your changes can be demoed or tested in a fully deployed environment, and +pr otherwise.
* the rest of the body must help a reviewer to understand the change, or a future developer to understand the implementation decisions.

It's never too late to ask questions or confirmation of the design if it doesn't work as expected or something else make what was requested impossible to implement as asked.

