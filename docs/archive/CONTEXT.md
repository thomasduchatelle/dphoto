# Archive

The physical file storage layer. Responsible for persisting original media files for long-term safety and for serving resized versions optimised for web rendering.

## Language

**Original**:
The unmodified media file as uploaded. Stored in cool/durable storage and never altered. _Avoid_: source file, raw file

**Miniature**:
A pre-cached resized version of an image, generated to speed up web rendering. Not a long-term safe copy. _Avoid_: thumbnail, resized image, preview

**CachedWidth**:
One of the fixed resolutions pre-generated and stored in hot cache: 360 px (minimum) and 2400 px (maximum). Requests for other widths are served by downscaling a cached version on the fly. _Avoid_: thumbnail size, resize target

**StoreRequest**:
A command to physically archive a media, carrying the catalog identifiers, the owner, the folder name, the signature hash, the original filename, the capture datetime, and a stream factory to read the file. _Avoid_: upload request, ingest request

**Location**:
The physical storage key (path) where a media's original file is kept. Tracked in the location index, keyed by `MediaId`. _Avoid_: file path, storage key, archive key

**CoolStorage**:
Durable, long-term storage for original media files (backed by S3). Safe by design; files are never deleted automatically. _Avoid_: primary storage, permanent storage, S3

**HotCache**:
Transient storage for pre-resized images (also backed by S3). Not durability-guaranteed; files may be evicted and regenerated on demand. _Avoid_: CDN, cache, temporary storage
