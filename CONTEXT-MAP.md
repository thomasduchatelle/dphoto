# Context Map

DPhoto is organised around four bounded contexts. Each manages its own data, language, and responsibilities.

## Contexts

- [Catalog](./docs/catalog/CONTEXT.md): organises medias into albums owned by an owner
- [Archive](./docs/archive/CONTEXT.md): stores original media files and serves resized versions
- [Backup](./docs/backup/CONTEXT.md): scans local volumes and uploads new medias into the archive and catalog
- [ACL](./docs/acl/CONTEXT.md): controls authentication and what each user is permitted to access

## Relationships

- **Backup → Archive**: Backup uploads new media files to the Archive for long-term storage
- **Backup → Catalog**: Backup indexes new medias into the Catalog, creating albums as needed
- **Catalog ↔ Archive**: Catalog holds the logical index of medias; Archive holds the physical files; both are keyed by `MediaId`
- **ACL → Catalog**: ACL enforces which owners and albums a user may read or modify
