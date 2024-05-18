package catalogviews

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"slices"
)

func NewAlbumView(
	FindAlbumByOwnerPort FindAlbumByOwnerPort,
	GetAlbumSharingGridPort GetAlbumSharingGridPort,
	FindAlbumsByIdsPort FindAlbumsByIdsPort,
	SharedWithUserPort SharedWithUserPort,
) *AlbumView {
	return &AlbumView{Providers: []ListAlbumsProvider{
		&OwnedAlbumListProvider{
			FindAlbumByOwnerPort:    FindAlbumByOwnerPort,
			GetAlbumSharingGridPort: GetAlbumSharingGridPort,
		},
		&SharedAlbumListProvider{
			FindAlbumsByIdsPort: FindAlbumsByIdsPort,
			SharedWithUserPort:  SharedWithUserPort,
		},
	}}
}

type VisibleAlbum struct {
	catalog.Album
	MediaCount         int                // Count is the number of medias on the album
	Visitors           []usermodel.UserId // Visitors are the users that can see the album ; only visible to the owner of the album
	OwnedByCurrentUser bool               // OwnedByCurrentUser is set to true when the user is an owner of the album
}

type ListAlbumsFilter struct {
	OnlyDirectlyOwned bool // OnlyDirectlyOwned provides a sub-view where only resources directly owned by user are displayed and accessible
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
