TODOs
=======================================

**Epics:**

1. [X] reduce API latency by segregating read and write storage (event-sourcing like)
2. [X] merging two React State into one more consistent to improve user experience of Sharing feature
   1. ~~refactoring to break down the massive reducer - needs suggestion and implementation~~
3. [ ] migrating to AWS CDK from Terraform and Serverless
   1. aiming for 2 stacks: long term data stores, and WEB overlay
   2. migration paths to be included and executed
   3. new domain to be used
4. [ ] migrating to Waku from native React: very light React Framework in beta/exploration
   1. progressive migration - decoupling behaviours from framework before re-integrating them to new framework
   2. parallel deployment - /v2 would get to the new UI
   3. Auth0 - moving to public IDP must be considered
   4. NPM - Yarn or PNPM don't look justified for this project
   5. Visual testing to rethink as the tools used seem discontinued and incompatible with new ones
5. [ ] Upload: support a sync mechanism from Android to AWS (existing backup software), post-upload management (inbox: deletion, rotation, datetime
   sliding, ...)
   1. Mobile -> S3 Landing
   2. S3 Landing -> DPhoto backup (support deletion and modification)
6. [ ] Back the CLI with APIs instead of direct accesses to underlying AWS services
   1. user concept is missing in the CLI (takes the owner as the user)
   2. move to Auth0 for authentication on both UI and CLI
7. [ ] Create a monitoring Dashboard with some stats (drive size, cache size, missed cache, popular resolutions, ...)
8. Other features:
   1. [ ] deletion of pictures
   2. [ ] update media timestamps to synchronise a timeline within an album with medias from several capturing devices (camera and phone)

**Small tasks:**

* [ ] Create a landing page for new users (no albums, no medias)
* [ ] Update the `dphoto configure ...` command: since migration from terraform, the proposed option doesn't work
* [ ] Update the readme file and other documentation: let's add some user-friendly screenshots ! ... and some C4 modeling

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
6. ~~remove catalogaclview package~~
7. ~~Ops script to recount the medias retrospectively~~

### Unsorted ...

* ~~the backup should WAIT the end before updating the views.~~ WON'T DO: to keep consistency the views are updated on each batch
* SCOPE should be renamed to PERMISSION
   * PERMISSION should have a generic type ('OWNER' or 'VISITOR') and a RESOURCE {TYPE, OWNER, ID}
* ~~FIX the media listing (reported to miss some files)~~
