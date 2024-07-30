package pkgfactory

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"github.com/thomasduchatelle/dphoto/pkg/archiveadapters/archivedynamo"
	"github.com/thomasduchatelle/dphoto/pkg/archiveadapters/asyncjobadapter"
	"github.com/thomasduchatelle/dphoto/pkg/archiveadapters/s3store"
	"github.com/thomasduchatelle/dphoto/pkg/singletons"
)

type SyncArchiveFactory struct{}
type AsyncArchiveFactory struct {
	SyncArchiveFactory
}

type initBackupMarker struct{}

func (a *SyncArchiveFactory) InitArchive(ctx context.Context) {
	singletons.MustSingleton(func() (initBackupMarker, error) {
		repositoryAdapter := archivedynamo.Must(archivedynamo.New(AWSFactory(ctx).GetDynamoDBClient(), AWSNames.DynamoDBName()))
		storeAdapter := s3store.NewWithS3Client(AWSFactory(ctx).GetS3Client(), AWSNames.ArchiveMainBucketName())
		cacheAdapter := s3store.NewWithS3Client(AWSFactory(ctx).GetS3Client(), AWSNames.ArchiveCacheBucketName())
		archiveAsyncAdapter := a.ArchiveAsyncJobAdapter(ctx)
		archive.Init(
			repositoryAdapter,
			storeAdapter,
			cacheAdapter,
			archiveAsyncAdapter,
		)

		return initBackupMarker{}, nil
	})
}

func (a *SyncArchiveFactory) ArchiveAsyncJobAdapter(ctx context.Context) archive.AsyncJobAdapter {
	return singletons.MustSingletonKey("AsyncJobAdapter", func() (archive.AsyncJobAdapter, error) {
		return archive.NewSyncJobAdapter(), nil
	})
}

func (a *AsyncArchiveFactory) ArchiveAsyncJobAdapter(ctx context.Context) archive.AsyncJobAdapter {
	return singletons.MustSingletonKey("AsyncJobAdapter", func() (archive.AsyncJobAdapter, error) {
		return asyncjobadapter.NewFromClients(AWSFactory(ctx).GetSNSClient(), AWSFactory(ctx).GetSQSClient(), AWSNames.ArchiveJobsSNSARN(), AWSNames.ArchiveJobsSQSURL(), asyncjobadapter.DefaultImagesPerMessage), nil
	})
}
