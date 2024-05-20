package catalogviews

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type GetAvailabilitiesByUserPort interface {
	GetAvailabilitiesByUser(ctx context.Context, user usermodel.UserId) ([]AlbumSize, error)
}

type GetAvailabilitiesByUserFunc func(ctx context.Context, user usermodel.UserId) ([]AlbumSize, error)

func (f GetAvailabilitiesByUserFunc) GetAvailabilitiesByUser(ctx context.Context, user usermodel.UserId) ([]AlbumSize, error) {
	return f(ctx, user)
}

type ProviderFactory interface {
	NewProvider(ctx context.Context, mediaCounterPort MediaCounterPort) ListAlbumsProvider
}

type ProviderFactoryFunc func(ctx context.Context, mediaCounterPort MediaCounterPort) ListAlbumsProvider

func (f ProviderFactoryFunc) NewProvider(ctx context.Context, mediaCounterPort MediaCounterPort) ListAlbumsProvider {
	return f(ctx, mediaCounterPort)
}

type MediaCounterInjector struct {
	GetAvailabilitiesByUserPort GetAvailabilitiesByUserPort
	ProviderFactories           []ProviderFactory
}

func (o *MediaCounterInjector) ListAlbums(ctx context.Context, user usermodel.CurrentUser, filter ListAlbumsFilter) ([]*VisibleAlbum, error) {
	view, err := o.GetAvailabilitiesByUserPort.GetAvailabilitiesByUser(ctx, user.UserId)
	if err != nil {
		return nil, err
	}

	var visibleAlbums []*VisibleAlbum
	for _, factory := range o.ProviderFactories {
		provider := factory.NewProvider(ctx, &MediaCounterFromView{AlbumSizes: view})
		albums, err := provider.ListAlbums(ctx, user, filter)
		if err != nil {
			return nil, err
		}

		visibleAlbums = append(visibleAlbums, albums...)
	}

	return visibleAlbums, nil
}

type MediaCounterFromView struct {
	AlbumSizes []AlbumSize
}

func (m *MediaCounterFromView) CountMedia(ctx context.Context, album ...catalog.AlbumId) (map[catalog.AlbumId]int, error) {
	counts := make(map[catalog.AlbumId]int)
	for _, albumId := range album {
		for _, size := range m.AlbumSizes {
			if size.AlbumId.IsEqual(albumId) {
				counts[albumId] = size.MediaCount
				break
			}
		}
	}

	return counts, nil
}
