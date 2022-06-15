package bootstrap

import (
	"crypto/sha1"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/adapters/avi"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/adapters/exif"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/adapters/filesystem"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/adapters/localstorage"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/adapters/m2ts"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/adapters/mp4"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/adapters/onlinestorage"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/adapters/s3source"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/adapters/volumes"
	"github.com/thomasduchatelle/dphoto/dphoto/backup/interactors"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
	"os"
	"path"
)

func init() {
	interactors.DetailsReaders = append(interactors.DetailsReaders,
		new(exif.Parser),
		new(m2ts.Parser),
		new(mp4.Parser),
		new(avi.Parser),
	)
	interactors.SourcePorts[backup.VolumeTypeFileSystem] = new(filesystem.FsHandler)

	config.Listen(func(cfg config.Config) {
		log.Debugln("connecting backup adapters")
		interactors.OnlineStoragePort = onlinestorage.Must(onlinestorage.NewS3OnlineStorage(cfg.GetString("backup.s3.bucket"), cfg.GetAWSSession()))

		var err error
		interactors.DownloaderPort, err = localstorage.NewLocalStorage(os.ExpandEnv(cfg.GetString("backup.buffer.path")), cfg.GetInt("backup.buffer.size"))
		if err != nil {
			panic(err)
		}

		owner := cfg.GetString("owner")
		interactors.VolumeRepositoryPort = &volumes.FileSystemRepository{
			Directory: path.Join(os.ExpandEnv(cfg.GetString("backup.volumes.repository.directory")), fmt.Sprintf("%x", sha1.Sum([]byte(owner)))),
		}

		// note - use contextual credential to access S3 volume, not the one used by DPhoto.
		interactors.SourcePorts[backup.VolumeTypeS3] = s3source.NewS3Source(session.Must(session.NewSession()))
	})
}