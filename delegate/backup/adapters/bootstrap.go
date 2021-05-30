package adapters

import (
	"duchatelle.io/dphoto/dphoto/backup"
	"duchatelle.io/dphoto/dphoto/backup/adapters/images"
	"duchatelle.io/dphoto/dphoto/backup/adapters/localstorage"
	"duchatelle.io/dphoto/dphoto/backup/adapters/onlinestorage"
	"duchatelle.io/dphoto/dphoto/backup/adapters/volumes"
	"duchatelle.io/dphoto/dphoto/backup/model"
	config2 "duchatelle.io/dphoto/dphoto/internal/config"
	filesystem2 "duchatelle.io/dphoto/dphoto/internal/sources/filesystem"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	backup.ImageDetailsReader = new(images.ExifReader)
	backup.ScannerAdapters[model.VolumeTypeFileSystem] = new(filesystem2.FsHandler)

	config2.Listen(func(cfg config2.Config) {
		log.Debugln("connecting backup adapters")
		backup.OnlineStorage = onlinestorage.Must(onlinestorage.NewS3OnlineStorage(cfg.GetString("backup.s3.bucket"), cfg.GetAWSSession()))

		var err error
		backup.Downloader, err = localstorage.NewLocalStorage(os.ExpandEnv(cfg.GetString("backup.buffer.path")), cfg.GetInt("backup.buffer.size"))
		if err != nil {
			panic(err)
		}

		backup.VolumeRepository = &volumes.FileSystemRepository{
			Directory: os.ExpandEnv(cfg.GetString("backup.volumes.repository.directory")),
		}
	})
}
