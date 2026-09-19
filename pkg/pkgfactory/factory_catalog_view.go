package pkgfactory

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/acl/catalogacl"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviewsadapters/catalogviewsdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/singletons"
)

func AlbumViewRepository(ctx context.Context) *catalogviewsdynamodb.AlbumViewRepository {
	return singletons.MustSingleton(func() (*catalogviewsdynamodb.AlbumViewRepository, error) {
		return &catalogviewsdynamodb.AlbumViewRepository{
			Client:    AWSFactory(ctx).GetDynamoDBClient(),
			TableName: AWSNames.DynamoDBName(),
		}, nil
	})
}

func CatalogToACLAdapter(ctx context.Context) *catalogacl.ReverseReader {
	return singletons.MustSingleton(func() (*catalogacl.ReverseReader, error) {
		return &catalogacl.ReverseReader{
			ScopeRepository: AclQueries(ctx),
		}, nil
	})
}

func AlbumView(ctx context.Context) *catalogviews.AlbumView {
	return singletons.MustSingleton(func() (*catalogviews.AlbumView, error) {
		albumQueries := AlbumQueries(ctx)
		return catalogviews.NewAlbumView(
			AlbumViewRepository(ctx),
			CatalogToACLAdapter(ctx),
			albumQueries,
			albumQueries,
			CatalogToACLAdapter(ctx),
		), nil
	})
}

func OwnerDriftReconciler(ctx context.Context, dry bool, options ...catalogviews.DriftOption) *catalogviews.OwnerDriftReconciler {
	albumQueries := AlbumQueries(ctx)
	repository := AlbumViewRepository(ctx)

	drifts := make([]catalogviews.DriftOption, len(options)+1)
	copy(drifts, options)
	drifts[len(options)] = catalogviews.DriftOptionDryMode(dry, repository)

	return catalogviews.NewDriftReconciler(
		albumQueries,
		repository,
		CatalogToACLAdapter(ctx),
		albumQueries,
		albumQueries,
		drifts...,
	)
}
