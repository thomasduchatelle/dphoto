TODOs
=======================================

1. reduce API latency by segregating read and write storage (event-sourcing like)
2. Back the CLI with APIs instead of direct accesses to underlying AWS services
3. Upload: support DAV (backup from phone), post-upload management (deletion, rotation, datetime re-tagging)
4. Migrate to NextJS, and refresh UI build process
5. Create a monitoring Dashboard with some stats (drive size, cache size, missed cache, popular resolutions, ...)

API latency
---------------------------------------

Next steps:

* generalise the use of the AWSFactory, and create a Factory for each use case
* re-implement repository to use event streaming on catalog (especially albums)
* adds commands in `dphoto-ops` to operate the DB (create indexes, migrate the data, ...)
* create new index for Medias to be found by date (without the albums)