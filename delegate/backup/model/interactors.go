package model

// Source populates media channel with everything found on the volume.
type Source func(medias chan FoundMedia) (uint, uint, error)

// Filter returns true if file should be backed up
type Filter func(found FoundMedia) bool

// Analyser reads the header of the file to find metadata (EXIF, dimensions, ...)
type Analyser func(found FoundMedia) (*AnalysedMedia, error)

// Downloader downloads locally the file to avoid multi-reads and too high concurrency on slow media
type Downloader func(found FoundMedia) (FoundMedia, error)

// Uploader backups media on an online storage (and update the indexes)
type Uploader func(buffer []*AnalysedMedia, progressChannel chan *ProgressEvent) error

// PreCompletion is called just before the run complete.
type PreCompletion func() error
