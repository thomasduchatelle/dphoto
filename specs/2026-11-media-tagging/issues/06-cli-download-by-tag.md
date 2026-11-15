# 06 — CLI: download originals by tag

Status: ready

## Description

Add a CLI command to download the original files of every media carrying a given tag, across all albums.

`dphoto download --tag <tag> --to <dir>` resolves all media carrying the (normalised) tag for the owner via the `{owner}#{tag}` GSI, then downloads each original from S3 into the target directory. The CLI runs on the client with direct DynamoDB and S3 access (no REST), consistent with the backup command.

Load the `go` skill. Relevant: `cmd/dphoto/cmd` for the command, `pkg/catalog` tag query, and the archive adapters for fetching originals from S3.

## Dependencies

02.1.

## Acceptance criteria

- `dphoto download --tag <tag> --to <dir>` downloads the originals of all medias with that tag, across albums, into the target directory.
- The tag argument is normalised the same way as the rest of the system; `favorite` is accepted.
- Clear feedback on progress and on an empty result; existing files are not re-downloaded needlessly.
- Errors on individual medias do not abort the whole run; a summary reports failures.
- Tests follow the Go table-driven style.

## Out of scope

- Any web/API work (the CLI goes direct to DynamoDB/S3).
- Filtering by multiple tags or by location.
- Re-packaging as a zip or server-side download.
