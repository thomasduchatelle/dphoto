package catalogviews

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"slices"
)

func NewAlbumView(
	FindAlbumByOwnerPort FindAlbumByOwnerPort,
	GetAlbumSharingGridPort GetAlbumSharingGridPort,
	FindAlbumsByIdsPort FindAlbumsByIdsPort,
	SharedWithUserPort SharedWithUserPort,
	GetAvailabilitiesByUserPort GetAvailabilitiesByUserPort,
) *AlbumView {
	return &AlbumView{Providers: []ListAlbumsProvider{
		&MediaCounterInjector{
			GetAvailabilitiesByUserPort: GetAvailabilitiesByUserPort,
			ProviderFactories: []ProviderFactory{
				ProviderFactoryFunc(func(ctx context.Context, mediaCounterPort MediaCounterPort) ListAlbumsProvider {
					return &OwnedAlbumListProvider{
						FindAlbumByOwnerPort:    FindAlbumByOwnerPort,
						GetAlbumSharingGridPort: GetAlbumSharingGridPort,
						MediaCounterPort:        mediaCounterPort,
					}
				}),
				ProviderFactoryFunc(func(ctx context.Context, mediaCounterPort MediaCounterPort) ListAlbumsProvider {
					return &SharedAlbumListProvider{
						FindAlbumsByIdsPort: FindAlbumsByIdsPort,
						SharedWithUserPort:  SharedWithUserPort,
						MediaCounterPort:    mediaCounterPort,
					}
				}),
			},
		},
	}}
}

type ListAlbumsProvider interface {
	ListAlbums(ctx context.Context, user usermodel.CurrentUser, filter ListAlbumsFilter) ([]*VisibleAlbum, error)
}

type ListAlbumsProviderFunc func(ctx context.Context, user usermodel.CurrentUser, filter ListAlbumsFilter) ([]*VisibleAlbum, error)

func (f ListAlbumsProviderFunc) ListAlbums(ctx context.Context, user usermodel.CurrentUser, filter ListAlbumsFilter) ([]*VisibleAlbum, error) {
	return f(ctx, user, filter)
}

type AlbumView struct {
	Providers []ListAlbumsProvider
}

// ListAlbums returns albums visible by the user (owned by current user, and shared to him)
func (v *AlbumView) ListAlbums(ctx context.Context, user usermodel.CurrentUser, filter ListAlbumsFilter) ([]*VisibleAlbum, error) {
	var albums []*VisibleAlbum
	for _, provider := range v.Providers {
		view, err := provider.ListAlbums(ctx, user, filter)
		if err != nil {
			return nil, err
		}

		albums = append(albums, view...)
	}

	slices.SortFunc(albums, func(a, b *VisibleAlbum) int {
		if a.Start.Equal(b.Start) {
			return b.End.Compare(a.End)
		}

		return b.Start.Compare(a.Start)
	})

	return albums, nil
}
