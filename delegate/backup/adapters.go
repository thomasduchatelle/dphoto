package backup

var (
	VolumeRepository   VolumeRepositoryAdapter
	FileHandler        FileHandlerAdapter
	ImageDetailsReader ImageDetailsReaderAdapter
	OnlineStorage      OnlineStorageAdapter
)

type VolumeRepositoryAdapter interface {
	// return nil when not found
	FindVolumeMetadata(string) (*VolumeMetadata, error)
	CreateNewVolume(volume RemovableVolume) error
	RestoreLastSnapshot(volumeId string) ([]SimpleMediaSignature, error)
	StoreSnapshot(volumeId string, backupId string, signatures []SimpleMediaSignature) error
}

type FileHandlerAdapter interface {
	// close channel once complete, log and skip errors
	FindMediaRecursively(mountPath string, paths chan FoundMedia) error

	// Copy file and compute its SHA256
	CopyToLocal(originPath string, destPath string) (mediaHash string, err error)
}

type ImageDetailsReaderAdapter interface {
	ReadImageDetails(imagePath string) (*MediaDetails, error)
}

type OnlineStorageAdapter interface {
	// start N goroutine to backup files, register errors and clone #completionChannel once finished
	BackupOnline(mediaChannel chan LocalMedia) error
}
