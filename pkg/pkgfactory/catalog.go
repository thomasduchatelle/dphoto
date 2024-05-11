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

func ArchiveTimelineMutationObserver() *catalogarchivesync.Observer {
	return singletons.MustSingleton(func() (*catalogarchivesync.Observer, error) {
		return new(catalogarchivesync.Observer), nil
	})
}

func CreateAlbumCase(ctx context.Context) *catalog.CreateAlbum {
	repository := CatalogRepository(ctx)
	return catalog.NewAlbumCreate(
		repository,
		repository,
		repository,
		ArchiveTimelineMutationObserver(),
	)
}

func CreateAlbumDeleteCase(ctx context.Context) *catalog.DeleteAlbum {
	repository := CatalogRepository(ctx)
	return catalog.NewDeleteAlbum(
		repository,
		repository,
		repository,
		repository,
		ArchiveTimelineMutationObserver(),
	)
}

func RenameAlbumCase(ctx context.Context) *catalog.RenameAlbum {
	// TODO Sharing and other artefacts should be transferred as well when renaming (recreating) an album
	repository := CatalogRepository(ctx)
	return catalog.NewRenameAlbum(
		repository,
		repository,
		repository,
		repository,
		repository,
		ArchiveTimelineMutationObserver(),
	)
}

func AmendAlbumDatesCase(ctx context.Context) *catalog.AmendAlbumDates {
	repository := CatalogRepository(ctx)
	return catalog.NewAmendAlbumDates(
		repository,
		repository,
		repository,
		repository,
		ArchiveTimelineMutationObserver(),
	)
}
