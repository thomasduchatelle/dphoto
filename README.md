DPhoto
====================================

Backup photo to your private AWS Cloud, and share them with family and friends. Core features:

| Feature | Version | Description |
| ------- | ------- | ----------- |
Backup medias | v1.0 | Backup photos and videos from USB Sticks, Flash drives, and Camera (USB) when they're plugged
Organise by album | v1.0 | Photos and videos are organised by album based on their creation date. Each album is a directory in S3.
Migration script | v1.0 | Medias already uploaded in S3 are re-ordered by albums and de-duplicated (with interactive command line interface)
Web viewer | - | See photos by albums, tag them, search by tag or dates
Media sharing | - | Albums can be shared and contributed by several users

Getting Started
------------------------------------

Install 'dphoto' command line interface and configure it using the following:

    go install github.com/thomasduchatelle/dphoto/dphoto@latest
    dphoto configure

Then use command `backup` to upload media in a directory, or scan to interactively organise the albums.

    dphoto scan /x/y/z
    dphoto backup /x/y/z

Contribute
------------------------------------

Components:

* [DPhoto CLI](./dphoto/README.md): installed on the end-user computer, backup photos and videos using command line interface
* `infra-data`: terraform project to create required infrastructure on AWS for the CLI to work. Project won't be re-usable in a different context without overriding backend and some other defaults.

### Releasing process

To release a new version, without CI:

1. verify everything is checked-in (`git status`)
2. push on `master`
3. **infra-data** - verify and approve the plan on [https://app.terraform.io/app/dphoto/workspaces](https://app.terraform.io/app/dphoto/workspaces)
4. **dphoto CLI** - create a git tag: `git tag v1.x.y && git push --tags`
