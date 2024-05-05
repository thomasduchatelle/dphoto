TODOs
=======================================

1. reduce API latency by segregating read and write storage (event-sourcing like)
2. Back the CLI with APIs instead of direct accesses to underlying AWS services
3. Upload: support DAV (backup from phone), post-upload management (deletion, rotation, datetime re-tagging)
4. Migrate to NextJS, and refresh UI build process

API latency
---------------------------------------

Next steps:

* generalise the use of the AWSFactory, and create a Factory for each use case
* re-implement repository to use event streaming on catalog (especially albums)
* adds operations in dphoto-ops to update / upgradce the DB
* create new index for Medias to be found by date (without the albums)