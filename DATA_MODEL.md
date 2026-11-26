DynamoDB Schema
=======================================

Data entries is consolidated into a single table, documented here for readability purpose.

Schema
---------------------------------------

### Entries

| PK                       | SK                                          | Description                                              | Module               |
|--------------------------|---------------------------------------------|----------------------------------------------------------|----------------------|
| {OWNER}#ALBUM            | ALBUM#{FOLDER_NAME}                         | Album metadata                                           | catalogdynamo        |
| {OWNER}#ALBUM            | ALBUM#{FOLDER_NAME}#COVERS                  | Album covers (whole ordered set, up to 4) as one item. Moves to the new SK on folder-change rename; unchanged on in-place rename; deleted when the album is deleted. | catalogdynamo        |
| {OWNER}#MEDIA#{id}       | #METADATA                                   | Media metadata, including the `Tags` **String Set** which is the **source of truth** of the media's tags (loaded alongside the media when an album is read). Tags are added/removed with an atomic `ADD`/`DELETE` set update, so concurrent writers merge instead of clobbering. The same update bumps `TagVersion`, a monotonic counter stamped onto the TAG# sibling. | catalogdynamo        | 
| {OWNER}#MEDIA#{id}       | TAG#{TAG}                                   | One derived index item per (media, tag) ; sibling of the media, carries MediaId and the `TagVersion` it was last written with (no album reference). Feeds the MediaByTagIndex GSI. On add/remove the set is updated first (minting a version), then the sibling is reconciled with a conditional write that only wins over an older `TagVersion`, so a stale concurrent add/remove of the same tag can never leave the tag in the set while deleting its sibling. A failure between the two writes can leave an orphan sibling, which `FindMediasByTag` drops (and WARN-logs) by checking each candidate against the `#METADATA` set. See [ADR-0006](docs/adr/0006-media-tag-writes-version-guarded-not-transactional.md). | catalogdynamo        |
| {OWNER}#MEDIA#{id}       | LOCATION#                                   | Media location if the archive                            | archivedynamo        |
| USER#{EMAIL}             | SCOPE#{TYPE}#{RESOURCE OWNER}#{RESOURCE ID} | Scopes allowed for a user (ownership, shared, ...)       | aclscopedynamodb     |
| USER#{EMAIL}             | IDENTITY#                                   | Details about the user (name, picture, ...)              | aclidentitydynamodb  |
| USER#{EMAIL}#ALBUMS_VIEW | OWNED#{OWNER}#{FOLDER_NAME}                 | (view) album summary for an album owned by the user: count + display fields (name/start/end)                   | catalogviewsdynamodb |
| USER#{EMAIL}#ALBUMS_VIEW | OWNED#{OWNER}#{FOLDER_NAME}#COVERS          | (view) covers (0-4) of an album owned by the user, sibling of the summary row                                  | catalogviewsdynamodb |
| USER#{EMAIL}#ALBUMS_VIEW | VISITOR#{OWNER}#{FOLDER_NAME}               | (view) album summary for an album shared with the user: count + display fields (name/start/end)                | catalogviewsdynamodb |
| USER#{EMAIL}#ALBUMS_VIEW | VISITOR#{OWNER}#{FOLDER_NAME}#COVERS        | (view) covers (0-4) of an album shared with the user, sibling of the summary row                               | catalogviewsdynamodb |
| REFRESH#{TOKEN}          | #REFRESH_SPEC                               | Refresh token                                            | aclrefreshdynamodb   |

### Global indexes

[create_update_table.go](pkg/awssupport/appdynamodb/create_update_table.go)

| Name                   | PK name / SK name              | PK                                      | SK                                          | Description                                                                                      |
|------------------------|--------------------------------|-----------------------------------------|---------------------------------------------|--------------------------------------------------------------------------------------------------|
| AlbumIndex             | AlbumIndexPK / AlbumIndexSK    | {OWNER}#{FOLDER_NAME}                   | #METADATA ; MEDIA#{DATETIME}#{ID}           | Catalog - Find medias of an album in one query ; each media carries its tags in `#METADATA` |
| ReverseLocationIndex   | LocationKeyPrefix / LocationId | {S3 KEY (WITHOUT FILE NAME)}            | {MEDIA ID}                                  | Archive - Warmup cache                                                                           |
| ReverseGrantIndex      | ResourceOwner / SK             | {OWNER}                                 | SCOPE#{TYPE}#{RESOURCE OWNER}#{RESOURCE ID} | ACL - list to whom resources are shared                                                          |
| RefreshTokenExpiration | SK / AbsoluteExpiryTime        | #REFRESH_SPEC                           | {DATETIME}                                  | OAuth - housekeeping old refresh token                                                           |
| AlbumViewByAlbumIndex  | AlbumViewIndexPK / SK          | ALBUM#{OWNER}#{FOLDER_NAME}#ALBUMS_VIEW | {AVAILABILITY_TYPE}#{OWNER}#{FOLDER_NAME}   | View - list every viewer row of an album (drives event-driven updates without an ACL round-trip) |
| MediaByTagIndex        | MediaByTagPK / MediaByTagSK    | {OWNER}#{TAG}                           | MEDIA#{MEDIA ID}                            | Catalog - list every media carrying a tag for an owner, across albums, in O(matches). Projection `KEYS_ONLY`: the reader derives the media id from the SK (`MEDIA#{id}`) then batch-reads each `#METADATA` for the real attributes. |
