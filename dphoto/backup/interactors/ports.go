package interactors

import (
	"github.com/thomasduchatelle/dphoto/dphoto/backup/backupmodel"
	"github.com/pkg/errors"
)

var (
	VolumeRepositoryPort backupmodel.VolumeRepositoryAdapter
	DetailsReaders       []backupmodel.DetailsReaderAdapter                                 // DetailsReaders is a list of specific details extractor can auto-register
	SourcePorts          = make(map[backupmodel.VolumeType]backupmodel.MediaScannerAdapter) // SourcePorts maps the type of volume with it's implementation
	OnlineStoragePort    backupmodel.OnlineStorageAdapter                                   // OnlineStoragePort creates a new OnlineStorageAdaptor or panic.
	DownloaderPort       backupmodel.DownloaderAdapter                                      // DownloaderPort creates a new instance of the DownloaderPort
)

func NewSource(volume backupmodel.VolumeToBackup, onCompletion func(uint, uint)) (func(medias chan backupmodel.FoundMedia) (uint, uint, error), error) {
	source, ok := SourcePorts[volume.Type]
	if !ok {
		return nil, errors.Errorf("No scanner implementation provided for volume type %s", volume.Type)
	}

	return func(medias chan backupmodel.FoundMedia) (uint, uint, error) {
		count, size, err := source.FindMediaRecursively(volume, func(media backupmodel.FoundMedia) {
			medias <- media
		})
		if err == nil {
			onCompletion(count, size)
		}
		return count, size, err
	}, nil
}