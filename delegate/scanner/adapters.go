package scanner

import (
	"io"
	"time"
)

var (
	ImageDetailsReader ImageDetailsReaderAdapter
	SourceAdapters     = make(map[VolumeType]MediaScannerAdapter) // SourceAdapters maps the type of volume with it's implementation
)

// ClosableMedia can be implemented alongside FoundMedia if the implementation requires to release resources once the media has been handled.
type ClosableMedia interface {
	Close() error
}

type MediaScannerAdapter interface {
	// FindMediaRecursively scan throw the VolumeToBackup and emit to the channel any media found. Interrupted in case of error.
	// returns number of items found, and size of these items
	FindMediaRecursively(volume VolumeToBackup, paths chan FoundMedia) (uint, uint, error)
}

type ImageDetailsReaderAdapter interface {
	ReadImageDetails(reader io.Reader, lastModifiedDate time.Time) (*MediaDetails, error)
}
type ReadableMedia interface {
	ReadMedia() (io.Reader, error)
	// SimpleSignature is used to get the size to upload
	SimpleSignature() *SimpleMediaSignature
}
