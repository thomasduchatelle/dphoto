# Catalog

The index of all medias, organised into albums owned by an owner. Responsible for album lifecycle (create, rename, delete, merge date ranges) and for answering queries about which medias belong to which album.

## Language

**Owner**:
The tenant under which all albums and medias are stored. Every resource in the catalog belongs to exactly one owner. _Avoid_: user, account, tenant

**Album**:
A named, time-bounded container for medias. An album spans a contiguous date range and holds all medias captured within that range for its owner. _Avoid_: collection, gallery, folder

**AlbumId**:
The composite key that uniquely identifies an album: `Owner` + `FolderName`. _Avoid_: album key, album reference

**FolderName**:
A normalised, URL-safe string (e.g. `/2024-01_Holiday`) that uniquely identifies an album within an owner. It is both the storage path and the business key. _Avoid_: album path, album slug

**Media**:
A single photo or video stored in the catalog. Identified by its `MediaId` and described by its `MediaSignature` and `MediaDetails`. _Avoid_: file, image, photo, asset

**MediaId**:
A stable, opaque string identifier for a media item, derived from its `MediaSignature`. _Avoid_: file ID, asset ID

**MediaSignature**:
The content-based business key of a media: SHA-256 hash + file size. Two files with the same signature are considered the same media. _Avoid_: hash, checksum, fingerprint

**MediaDetails**:
Technical metadata extracted from the file: dimensions, capture datetime, GPS coordinates, camera make/model, orientation, and video codec. _Avoid_: EXIF, metadata, file info

**MediaType**:
Classifies a media as `IMAGE`, `VIDEO`, or `OTHER`. _Avoid_: file type, mime type

**TimeRange**:
A half-open datetime interval `[Start, End)` used to define album boundaries and to filter medias. _Avoid_: date range, period, interval
