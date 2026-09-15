package catalogviews

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type ListSummariesForUserPort interface {
	ListSummariesForUser(ctx context.Context, userId usermodel.UserId) ([]UserAlbumSummary, error)
}

type ProviderFactory interface {
	NewProvider(ctx context.Context, mediaCounterPort MediaCounterPort) ListAlbumsProvider
}

type ProviderFactoryFunc func(ctx context.Context, mediaCounterPort MediaCounterPort) ListAlbumsProvider

func (f ProviderFactoryFunc) NewProvider(ctx context.Context, mediaCounterPort MediaCounterPort) ListAlbumsProvider {
	return f(ctx, mediaCounterPort)
}

// MediaCounterInjector is covering the function while the view only contains some of the data (the counts)
type MediaCounterInjector struct {
	ListSummariesForUserPort ListSummariesForUserPort
	ProviderFactories        []ProviderFactory
}

func (o *MediaCounterInjector) ListAlbums(ctx context.Context, user usermodel.CurrentUser, filter ListAlbumsFilter) ([]*VisibleAlbum, error) {
	userSummaries, err := o.ListSummariesForUserPort.ListSummariesForUser(ctx, user.UserId)
	if err != nil {
		return nil, err
	}

	view := make([]AlbumSummary, len(userSummaries))
	for i, userSummary := range userSummaries {
		view[i] = userSummary.AlbumSummary
	}

	var visibleAlbums []*VisibleAlbum
	for _, factory := range o.ProviderFactories {
		provider := factory.NewProvider(ctx, &MediaCounterFromView{Summaries: view})
		albums, err := provider.ListAlbums(ctx, user, filter)
		if err != nil {
			return nil, err
		}

		visibleAlbums = append(visibleAlbums, albums...)
	}

	return visibleAlbums, nil
}

type MediaCounterFromView struct {
	Summaries []AlbumSummary
}

func (m *MediaCounterFromView) CountMedia(ctx context.Context, album ...catalog.AlbumId) (map[catalog.AlbumId]int, error) {
	counts := make(map[catalog.AlbumId]int)
	for _, albumId := range album {
		for _, summary := range m.Summaries {
			if summary.AlbumId.IsEqual(albumId) {
				counts[albumId] = summary.MediaCount
				break
			}
		}
	}

	return counts, nil
}
