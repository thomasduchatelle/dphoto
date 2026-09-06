---
name: implement
description: "Implement a piece of work based on a spec or set of tickets."
disable-model-invocation: true
---

Implement the work described by the user in the spec or tickets.

Use /tdd where possible, at pre-agreed seams.

Run typechecking regularly, single test files regularly, and the full test suite once at the end.

Once done, use /code-review to review the work.

Commit your work to the current branch using the pattern: `<context>[/<layer>] - <summary> + <body>`. Never amend or force existing commits.

* **context** is `archive`, `catalog`, `ci`, ... Use `llm` when working on skills or agent documentation, and `proj` when working on the issue-tracker.
* **layer** is added if the change only impct a single layer: `web`, `cli`, `api`, ... 
* add in the body `+next` if you consider your changes can be demoed or tested in a fully deployed environment.
* the rest of the body must help a reviewer to understand the change, or a future developer to understand the implementation decision (something that have been discussed with the user.

