package bootstrap

import (
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"github.com/thomasduchatelle/dphoto/domain/archiveadapters/stores3"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
)

func init() {
	config.Listen(func(cfg config.Config) {
		repositoryAdapter := arepositorydynamo.Must(arepositorydynamo.New(cfg.GetAWSSession(), cfg.GetString(config.arepositorydynamodbTable), false))
		storeAdapter := stores3.Must(stores3.New(cfg.GetAWSSession(), cfg.GetString(config.ArchiveMainBucketName)))
		cacheAdapter := stores3.Must(stores3.New(cfg.GetAWSSession(), cfg.GetString(config.ArchiveCacheBucketName)))
		archiveAsyncAdapter := asyncjobsns.New(cfg.GetAWSSession(), cfg.GetString(config.ArchiveJobsSNSARN), cfg.GetString(config.ArchiveJobsSQSURL), asyncjobsns.DefaultImagesPerMessage)
		archive.Init(
			repositoryAdapter,
			storeAdapter,
			cacheAdapter,
			archiveAsyncAdapter,
		)
	})
}
