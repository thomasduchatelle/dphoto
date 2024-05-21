package pkgfactory

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviewsadapters/catalogviewsdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviewsadapters/catalogviewstoacl"
)

func AlbumViewRepository(ctx context.Context) *catalogviewsdynamodb.AlbumViewRepository {
	return &catalogviewsdynamodb.AlbumViewRepository{
		Client:    AWSFactory(ctx).GetDynamoDBClient(),
		TableName: AWSNames.DynamoDBName(),
	}
}

func CatalogToACLAdapter(ctx context.Context) *catalogviewstoacl.CatalogToACLAdapter {
	return &catalogviewstoacl.CatalogToACLAdapter{
		ScopeRepository: AclQueries(ctx),
	}
}

func AlbumView(ctx context.Context) *catalogviews.AlbumView {
	albumQueries := AlbumQueries(ctx)
	albumViewRepository := AlbumViewRepository(ctx)
	adapter := CatalogToACLAdapter(ctx)

	return catalogviews.NewAlbumView(
		albumQueries,
		adapter,
		albumQueries,
		adapter,
		albumViewRepository,
	)
}

func CommandHandlerAlbumSize(ctx context.Context) *catalogviews.CommandHandlerAlbumSize {
	albumQueries := AlbumQueries(ctx)
	albumViewRepository := AlbumViewRepository(ctx)
	adapter := CatalogToACLAdapter(ctx)

	return &catalogviews.CommandHandlerAlbumSize{
		MediaCounterPort:              albumQueries,
		ListUserWhoCanAccessAlbumPort: adapter,
		ViewWriteRepository:           albumViewRepository,
	}
}