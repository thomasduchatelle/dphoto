package interactors

import (
	"duchatelle.io/dphoto/dphoto/backup/model"
	"github.com/pkg/errors"
)

var (
	VolumeRepositoryPort   model.VolumeRepositoryAdapter
	ImageDetailsReaderPort model.ImageDetailsReaderAdapter
	SourcePorts            = make(map[model.VolumeType]model.MediaScannerAdapter) // SourcePorts maps the type of volume with it's implementation
	OnlineStoragePort      model.OnlineStorageAdapter                             // OnlineStoragePort creates a new OnlineStorageAdaptor or panic.
	DownloaderPort         model.DownloaderAdapter                                // DownloaderPort creates a new instance of the DownloaderPort
)

func NewSource(volume model.VolumeToBackup, onCompletion func(uint, uint)) (model.Source, error) {
	source, ok := SourcePorts[volume.Type]
	if !ok {
		return nil, errors.Errorf("No scanner implementation provided for volume type %s", volume.Type)
	}

	return func(medias chan model.FoundMedia) (uint, uint, error) {
		count, size, err := source.FindMediaRecursively(volume, func(media model.FoundMedia) {
			medias <- media
		})
		if err == nil {
			onCompletion(count, size)
		}
		return count, size, err
	}, nil
}
