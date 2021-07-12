package backupmodel

import (
	"io"
)

// VolumeRepositoryAdapter keeps media found previously on a volume to not re-analyse them or request the remote server if there are already backed-up.
type VolumeRepositoryAdapter interface {
	RestoreLastSnapshot(volumeId string) ([]SimpleMediaSignature, error)
	StoreSnapshot(volumeId string, backupId string, signatures []SimpleMediaSignature) error
}

// ClosableMedia can be implemented alongside FoundMedia if the implementation requires to release resources once the media has been handled.
type ClosableMedia interface {
	Close() error
}

type MediaScannerAdapter interface {
	// FindMediaRecursively scan throw the VolumeToBackup and emit to the channel any media found. Interrupted in case of error.
	// returns number of items found, and size of these items
	FindMediaRecursively(volume VolumeToBackup, callback func(FoundMedia)) (uint, uint, error)
}

type DetailsReaderAdapter interface {
	// Supports returns true if the file can be parsed with this reader. False otherwise.
	Supports(media FoundMedia, mediaType MediaType) bool
	// ReadDetails extracts metadata from the content of the file.
	ReadDetails(reader io.Reader, options DetailsReaderOptions) (*MediaDetails, error)
}

// DownloaderAdapter downloads locally files that are in a slow volume as they will need to be read twice each (analyse + sha256, and upload)
type DownloaderAdapter interface {
	DownloadMedia(media FoundMedia) (FoundMedia, error)
}

type OnlineStorageAdapter interface {
	// UploadFile uploads the file in the right folder but might change the name to avoid clash with other existing files. Use files name is always returned.
	UploadFile(owner string, media ReadableMedia, folderName, filename string) (string, error)
	// MoveFile moves physically a file in remote storage ; returns filename (different in case on conflict)
	MoveFile(owner string, folderName string, filename string, destFolderName string) (string, error)
}

type ReadableMedia interface {
	ReadMedia() (io.Reader, error)
	// SimpleSignature is used to get the size to upload
	SimpleSignature() *SimpleMediaSignature
}
