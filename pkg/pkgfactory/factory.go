package pkgfactory

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/awsfactory"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogarchivesync"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/catalogdynamo"
	"github.com/thomasduchatelle/dphoto/pkg/singletons"
)

var (
	AWSConfigFactory = awsfactory.NewContextualConfigFactory() // AWSConfigFactory can be overridden to use other AWS authentication means (default on AWS Default config)
	AWSNames         AWSAdapterNames                           // Names provides the config required by the adapters

)

type AWSAdapterNames interface {
	DynamoDBName() string
	ArchiveMainBucketName() string
	ArchiveCacheBucketName() string
	ArchiveJobsSNSARN() string
	ArchiveJobsSQSURL() string
}

func AWSFactory(ctx context.Context) *awsfactory.AWSFactory {
	return singletons.MustSingleton(func() (*awsfactory.AWSFactory, error) {
		return awsfactory.NewAWSFactory(ctx, AWSConfigFactory)
	})
}

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
