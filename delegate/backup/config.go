package backup

import "duchatelle.io/dphoto/dphoto/config"

const (
	scanBufferSize = 1024 * 16
)

var (
	imageReaderThreadCount int // imageReaderThreadCount is the number of goroutines that will be used to read metadata from file contents
	uploadThreadCount      int // uploadThreadCount is the number of goroutines that will be used to backup batch of files
	downloadThreadCount    int // downloadThreadCount is the number of concurrent download from volume to local storage allowed
	uploadBatchSize        int // uploadBatchSize number of media to process as a batch (ideally, should be a multiple of 25: the dynamodb write batch size)
)

func init() {
	config.Listen(func(cfg config.Config) {
		imageReaderThreadCount = cfg.GetIntOrDefault("backup.concurrency.imageReader", 4)
		downloadThreadCount = cfg.GetIntOrDefault("backup.concurrency.downloader", 2)
		uploadThreadCount = cfg.GetIntOrDefault("backup.concurrency.uploader", 2)
		uploadThreadCount = cfg.GetIntOrDefault("backup.onlinestorage.batchSize", 25)
	})
}
