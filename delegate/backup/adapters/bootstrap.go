package adapters

import (
	"duchatelle.io/dphoto/dphoto/backup"
	"duchatelle.io/dphoto/dphoto/backup/adapters/filesystem"
	"duchatelle.io/dphoto/dphoto/backup/adapters/images"
	"duchatelle.io/dphoto/dphoto/backup/adapters/localstorage"
	"duchatelle.io/dphoto/dphoto/backup/adapters/onlinestorage"
	"duchatelle.io/dphoto/dphoto/backup/adapters/volumes"
	"duchatelle.io/dphoto/dphoto/backup/model"
	"github.com/aws/aws-sdk-go/aws/session"
)

func init() {
	backup.VolumeRepository = nil
	backup.ImageDetailsReader = new(images.ExifReader)
	backup.ScannerAdapters[model.VolumeTypeFileSystem] = new(filesystem.FsHandler)
	backup.OnlineStorageFactory = func() backup.OnlineStorageAdapter {
		storage, err := onlinestorage.NewS3OnlineStorage(backup.OnlineBackupLocation, session.Must(session.NewSession()))
		if err != nil {
			panic(err)
		}
		return storage
	}
	backup.DownloaderFactory = func() backup.DownloaderAdapter {
		downloader, err := localstorage.NewLocalStorage(backup.LocalMediaPath, backup.LocalBufferAreaSizeInOctet)
		if err != nil {
			panic(err)
		}
		return downloader
	}
	backup.VolumeRepository = &volumes.FileSystemRepository{
		Directory: backup.LocalMediaPath,
	}
}
