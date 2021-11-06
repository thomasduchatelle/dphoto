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

    go install github.com/thomasduchatelle/dphoto/delegate@master
    dphoto configure

    dphoto backup /x/y/z

Contribute
------------------------------------

Components:

* [Delegate](./delegate/README.md): installed on the end-user computer, backup photos and videos when a media is plugged in.
