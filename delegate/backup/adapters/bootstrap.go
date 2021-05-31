package adapters

import (
	"duchatelle.io/dphoto/dphoto/backup/adapters/filesystem"
	"duchatelle.io/dphoto/dphoto/backup/adapters/images"
	"duchatelle.io/dphoto/dphoto/backup/adapters/localstorage"
	"duchatelle.io/dphoto/dphoto/backup/adapters/onlinestorage"
	"duchatelle.io/dphoto/dphoto/backup/adapters/volumes"
	"duchatelle.io/dphoto/dphoto/backup/model"
	"duchatelle.io/dphoto/dphoto/config"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	model.ImageDetailsReaderPort = new(images.ExifReader)
	model.SourcePorts[model.VolumeTypeFileSystem] = new(filesystem.FsHandler)

	config.Listen(func(cfg config.Config) {
		log.Debugln("connecting backup adapters")
		model.OnlineStoragePort = onlinestorage.Must(onlinestorage.NewS3OnlineStorage(cfg.GetString("backup.s3.bucket"), cfg.GetAWSSession()))

		var err error
		model.DownloaderPort, err = localstorage.NewLocalStorage(os.ExpandEnv(cfg.GetString("backup.buffer.path")), cfg.GetInt("backup.buffer.size"))
		if err != nil {
			panic(err)
		}

		model.VolumeRepositoryPort = &volumes.FileSystemRepository{
			Directory: os.ExpandEnv(cfg.GetString("backup.volumes.repository.directory")),
		}
	})
}
