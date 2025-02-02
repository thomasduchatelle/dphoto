package pkgfactory

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/awsfactory"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogarchiveasync"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogarchivesync"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogdynamo"
	"github.com/thomasduchatelle/dphoto/pkg/singletons"
)

type CatalogFactory interface {
	CreateAlbumCase(ctx context.Context) *catalog.CreateAlbum
	CreateAlbumDeleteCase(ctx context.Context) *catalog.DeleteAlbum
	RenameAlbumCase(ctx context.Context) *catalog.RenameAlbum
	AmendAlbumDatesCase(ctx context.Context) *catalog.AmendAlbumDates
}

type ArchiveAdapterForCatalog interface {
	ArchiveTimelineMutationObserver(ctx context.Context) catalog.TimelineMutationObserver
}

func CatalogRepository(ctx context.Context) *catalogdynamo.Repository {
	return singletons.MustSingleton(func() (*catalogdynamo.Repository, error) {
		return catalogdynamo.NewRepository(AWSFactory(ctx).GetDynamoDBClient(), AWSNames.DynamoDBName()), nil
	})
}

type SyncArchiveAdapterForCatalog struct{}

func (s *SyncArchiveAdapterForCatalog) ArchiveTimelineMutationObserver(ctx context.Context) catalog.TimelineMutationObserver {
	factory.InitArchive(ctx)
	return singletons.MustSingleton(func() (*catalogarchivesync.ArchiveSyncRelocator, error) {
		return new(catalogarchivesync.ArchiveSyncRelocator), nil
	})
}

type ASyncArchiveAdapterForCatalog struct {
	AWSFactory      awsfactory.AWSFactory
	AWSAdapterNames AWSAdapterNames
}

func (s *ASyncArchiveAdapterForCatalog) ArchiveTimelineMutationObserver(ctx context.Context) catalog.TimelineMutationObserver {
	return singletons.MustSingleton(func() (*catalogarchiveasync.ArchiveASyncRelocator, error) {
		return &catalogarchiveasync.ArchiveASyncRelocator{
			SQSClient: s.AWSFactory.GetSQSClient(),
			QueueUrl:  s.AWSAdapterNames.ArchiveRelocateJobsSQSURL(),
		}, nil
	})
}

func AlbumQueries(ctx context.Context) *catalog.AlbumQueries {
	return singletons.MustSingleton(func() (*catalog.AlbumQueries, error) {
		return &catalog.AlbumQueries{
			Repository: CatalogRepository(ctx),
		}, nil
	})
}

func InsertMediasCase(ctx context.Context) *catalog.InsertMedias {
	repository := CatalogRepository(ctx)
	return catalog.NewInsertMedias(
		repository,
		CommandHandlerAlbumSize(ctx),
	)
}

func CatalogMediaQueries(ctx context.Context) *catalog.MediaQueries {
	return singletons.MustSingleton(func() (*catalog.MediaQueries, error) {
		return &catalog.MediaQueries{
			MediaReadRepository: CatalogRepository(ctx),
		}, nil
	})
}

type SimpleCatalogFactory struct {
	ArchiveAdapterForCatalog ArchiveAdapterForCatalog
}

func (s *SimpleCatalogFactory) CreateAlbumCase(ctx context.Context) *catalog.CreateAlbum {
	repository := CatalogRepository(ctx)
	return catalog.NewAlbumCreate(
		repository,
		repository,
		repository,
		s.ArchiveAdapterForCatalog.ArchiveTimelineMutationObserver(ctx),
		CommandHandlerAlbumSize(ctx),
	)
}

func (s *SimpleCatalogFactory) CreateAlbumDeleteCase(ctx context.Context) *catalog.DeleteAlbum {
	repository := CatalogRepository(ctx)
	return catalog.NewDeleteAlbum(
		repository,
		repository,
		repository,
		repository,
		s.ArchiveAdapterForCatalog.ArchiveTimelineMutationObserver(ctx),
		CommandHandlerAlbumSize(ctx),
	)
}

func (s *SimpleCatalogFactory) RenameAlbumCase(ctx context.Context) *catalog.RenameAlbum {
	// TODO ACL Sharing and other resources should be transferred as well when renaming (recreating) an album
	repository := CatalogRepository(ctx)
	return catalog.NewRenameAlbum(
		repository,
		repository,
		repository,
		repository,
		repository,
		repository,
		s.ArchiveAdapterForCatalog.ArchiveTimelineMutationObserver(ctx),
		CommandHandlerAlbumSize(ctx),
	)
}

func (s *SimpleCatalogFactory) AmendAlbumDatesCase(ctx context.Context) *catalog.AmendAlbumDates {
	repository := CatalogRepository(ctx)
	return catalog.NewAmendAlbumDates(
		repository,
		repository,
		repository,
		repository,
		s.ArchiveAdapterForCatalog.ArchiveTimelineMutationObserver(ctx),
		CommandHandlerAlbumSize(ctx),
	)
}
