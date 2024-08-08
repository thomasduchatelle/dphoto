DynamoDB Schema
=======================================

Data entries is consolidated into a single table, documented here for readability purpose.

Schema
---------------------------------------

### Entries

| PK                       | SK                                          | Description                                              | Module               |
|--------------------------|---------------------------------------------|----------------------------------------------------------|----------------------|
| {OWNER}#ALBUM            | ALBUM#{FOLDER_NAME}                         | Album metadata                                           | catalogdynamo        |
| {OWNER}#MEDIA#{id}       | #METADATA                                   | Media metadata                                           | catalogdynamo        | 
| {OWNER}#MEDIA#{id}       | LOCATION#                                   | Media location if the archive                            | archivedynamo        |
| USER#{EMAIL}             | SCOPE#{TYPE}#{RESOURCE OWNER}#{RESOURCE ID} | Scopes allowed for a user (ownership, shared, ...)       | aclscopedynamodb     |
| USER#{EMAIL}             | IDENTITY#                                   | Details about the user (name, picture, ...)              | aclidentitydynamodb  |
| USER#{EMAIL}#ALBUMS_VIEW | OWNED#{OWNER}#{FOLDER_NAME}#COUNT           | (view) number of medias in an album owned by the user    | catalogviewsdynamodb |
| USER#{EMAIL}#ALBUMS_VIEW | VISITOR#{OWNER}#{FOLDER_NAME}#COUNT         | (view) number of medias in an album shared with the user | catalogviewsdynamodb |
| REFRESH#{TOKEN}          | #REFRESH_SPEC                               | Refresh token                                            | aclrefreshdynamodb   |

### Global indexes

[create_update_table.go](pkg/awssupport/appdynamodb/create_update_table.go)

| Name                   | PK name / SK name              | PK                           | SK                                          | Description                             |
|------------------------|--------------------------------|------------------------------|---------------------------------------------|-----------------------------------------|
| AlbumIndex             | AlbumIndexPK / AlbumIndexSK    | {OWNER}#{FOLDER_NAME}        | #METADATA                                   | Catalog - Find medias by albums         |
| ReverseLocationIndex   | LocationKeyPrefix / LocationId | {S3 KEY (WITHOUT FILE NAME)} | {MEDIA ID}                                  | Archive - Warmup cache                  |
| ReverseGrantIndex      | ResourceOwner / SK             | {OWNER}                      | SCOPE#{TYPE}#{RESOURCE OWNER}#{RESOURCE ID} | ACL - list to whom resources are shared |
| RefreshTokenExpiration | SK / AbsoluteExpiryTime        | #REFRESH_SPEC                | {DATETIME}                                  | OAuth - housekeeping old refresh token  |
