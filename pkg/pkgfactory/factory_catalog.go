package pkgfactory

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogarchivesync"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogdynamo"
	"github.com/thomasduchatelle/dphoto/pkg/singletons"
)

func CatalogRepository(ctx context.Context) *catalogdynamo.Repository {
	return singletons.MustSingleton(func() (*catalogdynamo.Repository, error) {
		return catalogdynamo.NewRepository(AWSFactory(ctx).GetDynamoDBClient(), AWSNames.DynamoDBName()), nil
	})
}

func ArchiveTimelineMutationObserver(ctx context.Context) *catalogarchivesync.ArchiveSyncRelocator {
	factory.InitArchive(ctx)
	return singletons.MustSingleton(func() (*catalogarchivesync.ArchiveSyncRelocator, error) {
		return new(catalogarchivesync.ArchiveSyncRelocator), nil
	})
}

func AlbumQueries(ctx context.Context) *catalog.AlbumQueries {
	return singletons.MustSingleton(func() (*catalog.AlbumQueries, error) {
		return &catalog.AlbumQueries{
			Repository: CatalogRepository(ctx),
		}, nil
	})
}

func CreateAlbumCase(ctx context.Context) *catalog.CreateAlbum {
	repository := CatalogRepository(ctx)
	return catalog.NewAlbumCreate(
		repository,
		repository,
		repository,
		ArchiveTimelineMutationObserver(ctx),
		CommandHandlerAlbumSize(ctx),
	)
}

func CreateAlbumDeleteCase(ctx context.Context) *catalog.DeleteAlbum {
	repository := CatalogRepository(ctx)
	return catalog.NewDeleteAlbum(
		repository,
		repository,
		repository,
		repository,
		ArchiveTimelineMutationObserver(ctx),
		CommandHandlerAlbumSize(ctx),
	)
}

func RenameAlbumCase(ctx context.Context) *catalog.RenameAlbum {
	// TODO ACL Sharing and other resources should be transferred as well when renaming (recreating) an album
	repository := CatalogRepository(ctx)
	return catalog.NewRenameAlbum(
		repository,
		repository,
		repository,
		repository,
		repository,
		repository,
		ArchiveTimelineMutationObserver(ctx),
		CommandHandlerAlbumSize(ctx),
	)
}

func AmendAlbumDatesCase(ctx context.Context) *catalog.AmendAlbumDates {
	repository := CatalogRepository(ctx)
	return catalog.NewAmendAlbumDates(
		repository,
		repository,
		repository,
		repository,
		ArchiveTimelineMutationObserver(ctx),
		CommandHandlerAlbumSize(ctx),
	)
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
