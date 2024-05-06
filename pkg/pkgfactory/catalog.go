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
	return &catalog.CreateAlbum{
		FindAlbumsByOwnerPort: repository,
		InsertAlbumPort:       repository,
		TransferMediasPort:    repository,
		TimelineMutationObservers: []catalog.TimelineMutationObserver{
			ArchiveTimelineMutationObserver(),
		},
	}
}
