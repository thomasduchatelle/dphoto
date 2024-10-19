package pkgfactory

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"github.com/thomasduchatelle/dphoto/pkg/archiveadapters/archivedynamo"
	"github.com/thomasduchatelle/dphoto/pkg/archiveadapters/asyncjobadapter"
	"github.com/thomasduchatelle/dphoto/pkg/archiveadapters/s3store"
	"github.com/thomasduchatelle/dphoto/pkg/singletons"
)

func (a *AWSCloud) InitArchive(ctx context.Context) {
	singletons.MustSingletonKey("InitArchive", func() (interface{}, error) {
		repositoryAdapter := archivedynamo.Must(archivedynamo.New(AWSFactory(ctx).GetDynamoDBClient(), AWSNames.DynamoDBName()))
		storeAdapter := s3store.NewWithS3Client(AWSFactory(ctx).GetS3Client(), AWSNames.ArchiveMainBucketName())
		cacheAdapter := s3store.NewWithS3Client(AWSFactory(ctx).GetS3Client(), AWSNames.ArchiveCacheBucketName())
		archiveAsyncAdapter := a.ArchiveFactory.ArchiveAsyncJobAdapter(ctx)
		archive.Init(
			repositoryAdapter,
			storeAdapter,
			cacheAdapter,
			archiveAsyncAdapter,
		)

		return new(interface{}), nil
	})
}

type SyncArchiveFactory struct{}

func (a *SyncArchiveFactory) ArchiveAsyncJobAdapter(ctx context.Context) archive.AsyncJobAdapter {
	return singletons.MustSingletonKey("SyncArchiveFactory.ArchiveAsyncJobAdapter", func() (archive.AsyncJobAdapter, error) {
		log.Info("Using archive.NewSyncJobAdapter() as ArchiveAsyncJobAdapter")
		return archive.NewSyncJobAdapter(), nil
	})
}

type AsyncArchiveFactory struct{}

func (a *AsyncArchiveFactory) ArchiveAsyncJobAdapter(ctx context.Context) archive.AsyncJobAdapter {
	return singletons.MustSingletonKey("AsyncArchiveFactory.ArchiveAsyncJobAdapter", func() (archive.AsyncJobAdapter, error) {
		log.Info("Using asyncjobadapter.NewFromClients() as ArchiveAsyncJobAdapter")
		return asyncjobadapter.NewFromClients(AWSFactory(ctx).GetSNSClient(), AWSFactory(ctx).GetSQSClient(), AWSNames.ArchiveJobsSNSARN(), AWSNames.ArchiveJobsSQSURL(), asyncjobadapter.DefaultImagesPerMessage), nil
	})
}
