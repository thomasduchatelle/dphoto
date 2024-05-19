DynamoDB Schema
=======================================

Each adapter extends core DB schema with extra supported type to consolidate in a single table all elements of the
application.

For readability purpose, the big picture is documented in this file.

Schema
---------------------------------------

### Entries

| PK                       | SK                                         | Description                                           | Module               |
|--------------------------|--------------------------------------------|-------------------------------------------------------|----------------------|
| {OWNER}#ALBUM            | ALBUM#{FOLDER_NAME}                        | Album metadata                                        | catalogdynamo        |
| {OWNER}#MEDIA#{id}       | #METADATA                                  | Media metadata                                        | catalogdynamo        | 
| {OWNER}#MEDIA#{id}       | LOCATION#                                  | Media location if the archive                         | archivedynamo        |
| USER#{EMAIL}             | SCOPE#{TYPE}#{OWNER/VISITOR}#{RESOURCE ID} | Scopes allowed for a user (ownership, shared, ...)    | aclscopedynamodb     |
| USER#{EMAIL}             | IDENTITY#                                  | Details about the user (name, picture, ...)           | aclidentitydynamodb  |
| USER#{EMAIL}#ALBUMS_VIEW | OWNED#{OWNER}#{FOLDER_NAME}#COUNT          | (view) count of media in an album owned by the user   | catalogviewsdynamodb |
| USER#{EMAIL}#ALBUMS_VIEW | VISITOR#{OWNER}#{FOLDER_NAME}#COUNT        | (view) count of medias in an album shared to the user | catalogviewsdynamodb |
| REFRESH#{TOKEN}          | #REFRESH_SPEC                              | Refresh token                                         | aclrefreshdynamodb   |
