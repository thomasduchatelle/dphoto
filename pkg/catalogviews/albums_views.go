package catalogviews

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"sort"
)

type VisibleAlbum struct {
	catalog.Album
	MediaCount         int                // Count is the number of medias on the album
	Visitors           []usermodel.UserId // Visitors are the users that can see the album ; only visible to the owner of the album
	OwnedByCurrentUser bool               // OwnedByCurrentUser is set to true when the user is an owner of the album
}

type ListAlbumsFilter struct {
	OnlyDirectlyOwned bool // OnlyDirectlyOwned provides a sub-view where only resources directly owned by user are displayed and accessible
}

type AlbumView struct {
	FindAlbumByOwnerPort   FindAlbumByOwnerPort
	FindAlbumsByIdsPort    FindAlbumsByIdsPort
	FindAlbumSharingToPort FindAlbumSharingToPort
	SharedWithUserPort     SharedWithUserPort
}

type FindAlbumByOwnerPort interface {
	FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error)
}

type FindAlbumsByIdsPort interface {
	FindAlbumsById(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error)
}

type FindAlbumSharingToPort interface {
	GetAlbumSharingGrid(ctx context.Context, owner ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error)
}

type SharedWithUserPort interface {
	ListAlbumIdsSharedWithUser(ctx context.Context, userId usermodel.UserId) ([]catalog.AlbumId, error)
}

// ListAlbums returns albums visible by the user (owned by current user, and shared to him)
func (v *AlbumView) ListAlbums(ctx context.Context, user usermodel.CurrentUser, filter ListAlbumsFilter) ([]*VisibleAlbum, error) {
	view, err := v.listOwnedAlbums(ctx, user)
	if err != nil {
		return nil, err
	}

	if !filter.OnlyDirectlyOwned {
		albums, err := v.listSharedWithUserAlbums(ctx, user)
		if err != nil {
			return nil, err
		}

		view = append(view, albums...)
		sort.Slice(view, func(i, j int) bool {
			return view[i].Start.Before(view[j].Start)
		})
	}

	return view, nil
}

func (v *AlbumView) listOwnedAlbums(ctx context.Context, user usermodel.CurrentUser) ([]*VisibleAlbum, error) {
	if user.Owner == nil {
		return nil, nil
	}

	ownedAlbums, err := v.FindAlbumByOwnerPort.FindAlbumsByOwner(ctx, *user.Owner)
	if err != nil {
		return nil, err
	}

	sharing, err := v.FindAlbumSharingToPort.GetAlbumSharingGrid(ctx, *user.Owner)

	var view []*VisibleAlbum
	for _, album := range ownedAlbums {
		sharedTo, _ := sharing[album.AlbumId]
		view = append(view, &VisibleAlbum{
			Album:              *album,
			MediaCount:         album.TotalCount,
			Visitors:           sharedTo,
			OwnedByCurrentUser: true,
		})
	}

	return view, err
}

func (v *AlbumView) listSharedWithUserAlbums(ctx context.Context, user usermodel.CurrentUser) ([]*VisibleAlbum, error) {
	shares, err := v.SharedWithUserPort.ListAlbumIdsSharedWithUser(ctx, user.UserId)
	if err != nil || len(shares) == 0 {
		return nil, err
	}

	albums, err := v.FindAlbumsByIdsPort.FindAlbumsById(ctx, shares)

	var view []*VisibleAlbum
	for _, album := range albums {
		view = append(view, &VisibleAlbum{
			Album:      *album,
			MediaCount: album.TotalCount,
		})
	}

	return view, err
}
