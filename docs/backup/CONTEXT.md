# Backup

The pipeline that scans a local volume for media files, cross-references them against the catalog to identify new ones, and uploads them to the archive while indexing them into the catalog.

## Language

**SourceVolume**:
An abstraction for any scannable media source — local disk, USB drive, Android device, or remote storage. Enumerates `FoundMedia` items without knowledge of the pipeline stages that follow. _Avoid_: disk, drive, source, folder

**FoundMedia**:
A file discovered on a `SourceVolume`, providing its path, size, last-modification time, and a stream reader. It has not yet been analysed. _Avoid_: file, raw file, discovered file

**AnalysedMedia**:
A `FoundMedia` enriched with its computed SHA-256 hash, `MediaType`, and `MediaDetails`. The output of the analysis stage; ready for catalog cross-referencing. _Avoid_: processed file, hashed file

**CatalogReference**:
The result of cross-referencing an `AnalysedMedia` with the catalog: resolves whether the media already exists, which album it belongs to, and its `MediaId`. _Avoid_: catalog lookup, deduplication result

**BackingUpMediaRequest**:
A pairing of an `AnalysedMedia` with its `CatalogReference`, representing a media confirmed for upload. Passed to the upload stage. _Avoid_: upload item, queued media

**Scan**:
A read-only run of the backup pipeline that analyses files and produces `ScannedFolder` summaries without uploading or modifying the catalog. _Avoid_: dry run, preview, analysis run

**ScannedFolder**:
A summary of one directory found during a scan: date bounds, per-day media distribution, and reject count. _Avoid_: folder summary, directory report

**Report**:
The summary result of a completed backup run: counts of skipped medias and per-album breakdown of uploaded medias. _Avoid_: backup result, upload summary

**RejectedMedia**:
A media that could not be analysed or catalogued — unreadable file, missing capture date, or unsupported format. Counted in the `Report` but not uploaded. _Avoid_: failed media, skipped file, error
