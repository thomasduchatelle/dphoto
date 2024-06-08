TODOs
=======================================

1. reduce API latency by segregating read and write storage (event-sourcing like)
2. Back the CLI with APIs instead of direct accesses to underlying AWS services
   1. user concept is missing in the CLI (takes the owner as the user)
   2. move to Auth0 for authentication on both UI and CLI
3. Upload: support DAV (backup from phone), post-upload management (deletion, rotation, datetime re-tagging)
4. Migrate to NextJS, and refresh UI build process
5. Create a monitoring Dashboard with some stats (drive size, cache size, missed cache, popular resolutions, ...)

API latency
---------------------------------------

Next steps:

* ~~generalise the use of the AWSFactory, and create a Factory for each use case~~
* ~~re-implement repository to use event streaming on catalog (especially albums)~~
* ~~adds commands in `dphoto-ops` to operate the DB (create indexes, migrate the data, ...)~~
* create new index for Medias to be found by date (without the albums)

### Catalog View

1. ~~implements the ports~~
2. ~~integrate it with consumers (API and CLI) -> requires to add 'user' concept to the CLI~~
3. ~~refactor to use a list of providers, each will fetch a list of albums that will be sorted after: OwnedAlbumsProvider, SharedAlbumsProvider,~~
   ~~OfflineViewProviders, etc.~~
4. ~~do the count of media as part of this view ; not the repository~~
5. ~~start building a view that contains for a user:~~
   * owned albums + count + last media date + sharing
   * shared albums [by owner] + count + last media date
6. remove catalogaclview package
7. Ops script to recount the medias retrospectively

### Unsorted ...

* the backup should WAIT the end before updating the views.
* SCOPE should be renamed to PERMISSION
   * PERMISSION should have a generic type ('OWNER' or 'VISITOR') and a RESOURCE {TYPE, OWNER, ID}  
