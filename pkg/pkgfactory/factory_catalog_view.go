package pkgfactory

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviewsadapters/catalogviewstoacl"
)

func AlbumView(ctx context.Context) *catalogviews.AlbumView {
	repository := CatalogQueries(ctx)

	adapter := &catalogviewstoacl.FindAlbumSharingToAdapter{
		ScopeRepository: AclQueries(ctx),
	}
	return catalogviews.NewAlbumView(
		repository,
		adapter,
		repository,
		adapter,
		repository,
	)

}
