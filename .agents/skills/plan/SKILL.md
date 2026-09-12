---
name: plan
description: "Present an clear implementation plan to the user covering persistence, service, UI, and API layers."
---

The planning phase is an interactive session with the user to settle the design decisions to take before implementing the requirements. **Take advantage of the planning to ask clarification questions.**

Use the following communication tools to present your recommended approach.

## Data model

_Use if a persistence layer is required to implement the ticket._

By data model, you need to describe:

* the entities model: name, properties (with types), primary/secondary indexes.
* the relationships between entities.
* the access patterns (queries and writes), and the indexes required to support them.
* the compatibility strategy: backward compatibility or upfront migration.

## Data flow diagram

_Use if the changes spread across several modules or layers in a non-trivial way._

Explain how the modules are sequenced together, for each use-case affected by the change.

Notes:

* an actual diagram is optional, a well-structured text is accepted.
* a use-case is described by a user, or external system, interaction.
* a module is anything with an interface and an implementation. Deliberately scale-agnostic: a function, a structure, a domain/context, a port/adapter, a UI component, or a tier-spanning slice.

## Modules

_Use if the changes are affecting the interface of one or more modules._

Describe the interface of the new modules, or how the existing modules are affected. Explain the refactoring required if one is necessary.


## Testing strategy

_Use always._

We use the TDD principles: no code should be written if not to satisfy a failing test. Sequence the tests that will be required to push forward the implementation, or how the existing tests will be covering the refactoring.

Extending tests or removing tests should be exceptional, justify each of them. Instead of extending/updating a test, it is preferred to create a new one and flag the existing one as redundant: it demonstrates that the change didn't regress any of the current functionalities.

Add a reference to each test case you describe to help the user to refer back to them in its response.


## REST API

_Use if the operation must be exposed through a REST API, or if it affect an existing endpoint._

Present the new endpoint, or the change to the existing one:

* method, path, path parameters, query parameters
* authorisation rules, and implementation

## Retirement & Simplification

_Use if relevant._

While not required to achieve the requirement, flag any of the following:

* misleading naming - propose to rename elements that are not representing any more what they are (functions, variables, structures, ...).
* stale documentation - propose to update the documentation that became stale due to this change.
* dead branch and redundant code - recommend the deletion of existing code which is not used.
* simplification - flag when some code still uses exiting code while it could use the new one, you can propose to create a ticket instead of creeping this one.
* redundant tests - propose to add a comment, and recommend to delete them in a follow-up changes.
