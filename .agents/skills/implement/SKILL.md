---
name: implement
description: "Implement a piece of work based on a spec or set of tickets. Describes the commit message pattern and the pull request template."
---

Implement the work described by the user in the spec or tickets.

## Implementation

Load the skills related to the language your working with.

Run single test files regularly, and the full test suite once at the end. 

Before handing over your work:

1. run typechecking and linting,
2. progress the issue(s) you're working on as `done`,
3. commit if you are not working interactively with a user,
4. create a pull request if explicitly requested.


## Commit messages

Use the pattern: `<context>[/<layer>] - <summary> + <body>`. Never amend or force existing commits.

* **context** is `archive`, `catalog`, `ci`, ... Use `llm` when working on skills or agent documentation, and `proj` when working on the issue-tracker.
* **layer** is added if the change only impct a single layer: `web`, `cli`, `api`, ...
* the rest is a **short** description, bullet points, of what has been changed. Explain the intention like "rename variable X", not a list of updated files.

## Pull request

Use the main (or first) commit headline as the pull request title.

In the body, you must explain how it works, and why it works, to the reviewer, step by step. 

* Start where the data comes in and follow the dataflow. 
* Describe the responsibility of the different modules, use the interface to explain how they interact and were their boundaries lie.
* Explain your testing strategy: what is tested at unit level, why a wider test is required, ...
* Keep the explanation on the intentions - why was that module changed or created - do not fall in an enumeration of changes. The reviewer will read soon enough each line changed, you need to explain him what he will be looking at.
* Refresh the reviewer of what was already present and used to deliver the requirements. A short recap is a time saver, but keep it extremely short, just enough to revive the reviewer memory.
