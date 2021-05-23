package backup

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"io"
	"time"
)

var (
	VolumeRepository   VolumeRepositoryAdapter
	ImageDetailsReader ImageDetailsReaderAdapter
	ScannerAdapters    = make(map[model.VolumeType]MediaScannerAdapter) // ScannerAdapters maps the type of volume with it's implementation
	OnlineStorage      OnlineStorageAdapter                             // OnlineStorage creates a new OnlineStorageAdaptor or panic.
	Downloader         DownloaderAdapter                                // Downloader creates a new instance of the Downloader
)

type VolumeRepositoryAdapter interface {
	RestoreLastSnapshot(volumeId string) ([]model.SimpleMediaSignature, error)
	StoreSnapshot(volumeId string, backupId string, signatures []model.SimpleMediaSignature) error
}

// ClosableMedia can be implemented alongside FoundMedia if the implementation requires to release resources once the media has been handled.
type ClosableMedia interface {
	Close() error
}

type MediaScannerAdapter interface {
	// FindMediaRecursively scan throw the VolumeToBackup and emit to the channel any media found. Interrupted in case of error.
	// returns number of items found, and size of these items
	FindMediaRecursively(volume model.VolumeToBackup, paths chan model.FoundMedia) (uint, uint, error)
}

type ImageDetailsReaderAdapter interface {
	ReadImageDetails(reader io.Reader, lastModifiedDate time.Time) (*model.MediaDetails, error)
}

type DownloaderAdapter interface {
	DownloadMedia(media model.FoundMedia) (model.FoundMedia, error)
}

type OnlineStorageAdapter interface {
	// UploadFile uploads the file in the right folder but might change the name to avoid clash with other existing files. Use files name is always returned.
	UploadFile(media ReadableMedia, folderName, filename string) (string, error)
}

type ReadableMedia interface {
	ReadMedia() (io.Reader, error)
	// SimpleSignature is used to get the size to upload
	SimpleSignature() *model.SimpleMediaSignature
}
