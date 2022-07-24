package bootstrap

import (
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"github.com/thomasduchatelle/dphoto/domain/archiveadapters/archivedynamo"
	"github.com/thomasduchatelle/dphoto/domain/archiveadapters/jobqueuesns"
	"github.com/thomasduchatelle/dphoto/domain/archiveadapters/s3store"
	"github.com/thomasduchatelle/dphoto/dphoto/config"
)

func init() {
	config.Listen(func(cfg config.Config) {
		repositoryAdapter := archivedynamo.Must(archivedynamo.New(cfg.GetAWSSession(), cfg.GetString(config.ArchiveDynamodbTable), false))
		storeAdapter := s3store.Must(s3store.New(cfg.GetAWSSession(), cfg.GetString(config.ArchiveMainBucketName)))
		cacheAdapter := s3store.Must(s3store.New(cfg.GetAWSSession(), cfg.GetString(config.ArchiveCacheBucketName)))
		archive.Init(repositoryAdapter, storeAdapter, cacheAdapter, jobqueuesns.New(cfg.GetAWSSession(), cfg.GetString(config.ArchiveJobsSNSARN)))
	})
}
