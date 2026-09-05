# CONTEXT.md Format

## Structure

```md
# {Context Name}

{One or two sentence description of what this context is and why it exists.}

## Language

**Order**:
{A one or two sentence description of the term}
_Avoid_: Purchase, transaction

**Invoice**:
A request for payment sent to a customer after delivery.
_Avoid_: Bill, payment request

**Customer**:
A person or organization that places orders.
_Avoid_: Client, buyer, account
```

## Rules

- **Be opinionated.** When multiple words exist for the same concept, pick the best one and list the others under `_Avoid_`.
- **Keep definitions tight.** One or two sentences max. Define what it IS, not what it does.
- **Only include terms specific to this project's context.** General programming concepts (timeouts, error types, utility patterns) don't belong even if the project uses them extensively. Before adding a term, ask: is this a concept unique to this context, or a general programming concept? Only the former belongs.
- **Group terms under subheadings** when natural clusters emerge. If all terms belong to a single cohesive area, a flat list is fine.

## Multi-context repos

This repo uses multiple contexts. Each context lives in `docs/{context}/CONTEXT.md`. A `CONTEXT-MAP.md` at the root lists all contexts and their relationships:

```md
# Context Map

## Contexts

- [Catalog](./docs/catalog/CONTEXT.md): organises medias into albums
- [Archive](./docs/archive/CONTEXT.md): stores and serves media files
- [Backup](./docs/backup/CONTEXT.md): scans local volumes and uploads new medias

## Relationships

- **Backup → Archive**: Backup uploads new media files to the Archive
- **Backup → Catalog**: Backup indexes new medias into the Catalog
- **Catalog ↔ Archive**: Catalog holds the index; Archive holds the files
```

When the current topic touches a context, update `docs/{context}/CONTEXT.md`. If unclear which context applies, ask.
