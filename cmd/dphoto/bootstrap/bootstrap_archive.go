package bootstrap

import (
	"fmt"
	"github.com/thomasduchatelle/dphoto/cmd/dphoto/config"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"github.com/thomasduchatelle/dphoto/pkg/archiveadapters/archivedynamo"
	"github.com/thomasduchatelle/dphoto/pkg/archiveadapters/asyncjobadapter"
	"github.com/thomasduchatelle/dphoto/pkg/archiveadapters/s3store"
)

func init() {
	config.Listen(func(cfg config.Config) {
		repositoryAdapter := archivedynamo.Must(archivedynamo.New(cfg.GetAWSFactory().GetDynamoDBClient(), cfg.GetString(config.ArchiveDynamodbTable)))
		storeAdapter := s3store.Must(s3store.New(cfg.GetAWSV2Config(), cfg.GetString(config.ArchiveMainBucketName)))
		cacheAdapter := s3store.Must(s3store.New(cfg.GetAWSV2Config(), cfg.GetString(config.ArchiveCacheBucketName)))
		archiveAsyncAdapter := asyncjobadapter.New(cfg.GetAWSV2Config(), cfg.GetString(config.ArchiveJobsSNSARN), cfg.GetString(config.ArchiveJobsSQSURL), asyncjobadapter.DefaultImagesPerMessage)
		archive.Init(
			repositoryAdapter,
			storeAdapter,
			cacheAdapter,
			archiveAsyncAdapter,
		)
	})
}

func must[M any](value M, err error) M {
	if err != nil {
		panic(fmt.Sprintf("PANIC - %T couldn't be built: %s", *new(M), err))
	}

	return value
}
