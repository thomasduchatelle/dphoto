package adapters

import (
	"duchatelle.io/dphoto/dphoto/backup/adapters/exif"
	"duchatelle.io/dphoto/dphoto/backup/adapters/filesystem"
	"duchatelle.io/dphoto/dphoto/backup/adapters/localstorage"
	"duchatelle.io/dphoto/dphoto/backup/adapters/m2ts"
	"duchatelle.io/dphoto/dphoto/backup/adapters/onlinestorage"
	"duchatelle.io/dphoto/dphoto/backup/adapters/volumes"
	"duchatelle.io/dphoto/dphoto/backup/backupmodel"
	"duchatelle.io/dphoto/dphoto/backup/interactors"
	"duchatelle.io/dphoto/dphoto/config"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	interactors.DetailsReaders[interactors.DetailsReaderTypeImage] = new(exif.Parser)
	interactors.DetailsReaders[interactors.DetailsReaderTypeM2TS] = new(m2ts.Parser)
	interactors.SourcePorts[backupmodel.VolumeTypeFileSystem] = new(filesystem.FsHandler)

	config.Listen(func(cfg config.Config) {
		log.Debugln("connecting backup adapters")
		interactors.OnlineStoragePort = onlinestorage.Must(onlinestorage.NewS3OnlineStorage(cfg.GetString("backup.s3.bucket"), cfg.GetAWSSession()))

		var err error
		interactors.DownloaderPort, err = localstorage.NewLocalStorage(os.ExpandEnv(cfg.GetString("backup.buffer.path")), cfg.GetInt("backup.buffer.size"))
		if err != nil {
			panic(err)
		}

		interactors.VolumeRepositoryPort = &volumes.FileSystemRepository{
			Directory: os.ExpandEnv(cfg.GetString("backup.volumes.repository.directory")),
		}
	})
}
