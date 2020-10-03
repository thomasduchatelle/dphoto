DPhoto
====================================

Backup photo to your private AWS Cloud, and share them with family and friends. Core features:

| Feature | Version | Description |
| ------- | ------- | ----------- |
Backup medias | **in dev** | Backup photos and videos from USB Sticks, Flash drives, and Camera (USB) when they're plugged
Organise by album | **in dev** | Photos and videos are organised by album based on their creation date. Each album is a directory in S3.
Web viewer | - | See photos by albums, tag them, search by tag or dates
Migration script | - | Medias already uploaded are re-ordered by albums and de-duplicated
Media sharing | - | Albums can be shared and contributed by several users

Contribute
------------------------------------

Components:

* [Delegate](./delegate/README.md): installed on the end-user computer, backup photos and videos when a media is plug-in.