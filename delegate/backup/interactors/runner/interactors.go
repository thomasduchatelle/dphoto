package runner

import "duchatelle.io/dphoto/dphoto/backup/backupmodel"

// Source populates media channel with everything found on the volume.
type Source func(medias chan backupmodel.FoundMedia) (uint, uint, error)

// Filter returns true if file should be backed up
type Filter func(found backupmodel.FoundMedia) bool

// Analyser reads the header of the file to find metadata (EXIF, dimensions, ...)
type Analyser func(found backupmodel.FoundMedia) (*backupmodel.AnalysedMedia, error)

// Downloader downloads locally the file to avoid multi-reads and too high concurrency on slow media
type Downloader func(found backupmodel.FoundMedia) (backupmodel.FoundMedia, error)

// Uploader backups media on an online storage (and update the indexes)
type Uploader func(buffer []*backupmodel.AnalysedMedia, progressChannel chan *backupmodel.ProgressEvent) error

// PreCompletion is called just before the run complete.
type PreCompletion func() error
