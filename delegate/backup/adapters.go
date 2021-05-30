package backup

import (
	"duchatelle.io/dphoto/dphoto/scanner"
	"io"
)

var (
	VolumeRepository VolumeRepositoryAdapter
	OnlineStorage    OnlineStorageAdapter // OnlineStorage creates a new OnlineStorageAdaptor or panic.
	Downloader       DownloaderAdapter    // Downloader creates a new instance of the Downloader
)

type VolumeRepositoryAdapter interface {
	RestoreLastSnapshot(volumeId string) ([]scanner.SimpleMediaSignature, error)
	StoreSnapshot(volumeId string, backupId string, signatures []scanner.SimpleMediaSignature) error
}

type DownloaderAdapter interface {
	DownloadMedia(media scanner.FoundMedia) (scanner.FoundMedia, error)
}

type OnlineStorageAdapter interface {
	// UploadFile uploads the file in the right folder but might change the name to avoid clash with other existing files. Use files name is always returned.
	UploadFile(media ReadableMedia, folderName, filename string) (string, error)
}

type ReadableMedia interface {
	ReadMedia() (io.Reader, error)
	// SimpleSignature is used to get the size to upload
	SimpleSignature() *scanner.SimpleMediaSignature
}
