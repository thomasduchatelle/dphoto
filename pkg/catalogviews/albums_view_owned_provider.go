package catalogviews

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type FindAlbumByOwnerPort interface {
	FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error)
}

type FindAlbumByOwnerFunc func(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error)

func (f FindAlbumByOwnerFunc) FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
	return f(ctx, owner)
}

type GetAlbumSharingGridPort interface {
	GetAlbumSharingGrid(ctx context.Context, owner ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error)
}

type GetAlbumSharingGridFunc func(ctx context.Context, owner ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error)

func (f GetAlbumSharingGridFunc) GetAlbumSharingGrid(ctx context.Context, owner ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error) {
	return f(ctx, owner)
}

type MediaCounterPort interface {
	CountMedia(ctx context.Context, album ...catalog.AlbumId) (map[catalog.AlbumId]int, error)
}

type MediaCounterFunc func(ctx context.Context, album ...catalog.AlbumId) (map[catalog.AlbumId]int, error)

func (f MediaCounterFunc) CountMedia(ctx context.Context, album ...catalog.AlbumId) (map[catalog.AlbumId]int, error) {
	return f(ctx, album...)
}

type OwnedAlbumListProvider struct {
	FindAlbumByOwnerPort    FindAlbumByOwnerPort
	GetAlbumSharingGridPort GetAlbumSharingGridPort
	MediaCounterPort        MediaCounterPort
}

func (o *OwnedAlbumListProvider) ListAlbums(ctx context.Context, user usermodel.CurrentUser, filter ListAlbumsFilter) ([]*VisibleAlbum, error) {
	if user.Owner == nil {
		return nil, nil
	}

	ownedAlbums, err := o.FindAlbumByOwnerPort.FindAlbumsByOwner(ctx, *user.Owner)
	if err != nil {
		return nil, err
	}

	albumIds := make([]catalog.AlbumId, len(ownedAlbums))
	for i, album := range ownedAlbums {
		albumIds[i] = album.AlbumId
	}

	mediaCount, err := o.MediaCounterPort.CountMedia(ctx, albumIds...)
	if err != nil {
		return nil, err
	}

	sharing, err := o.GetAlbumSharingGridPort.GetAlbumSharingGrid(ctx, *user.Owner)

	var view []*VisibleAlbum
	for _, album := range ownedAlbums {
		count, _ := mediaCount[album.AlbumId]
		sharedTo, _ := sharing[album.AlbumId]
		view = append(view, &VisibleAlbum{
			Album:              *album,
			MediaCount:         count,
			Visitors:           sharedTo,
			OwnedByCurrentUser: true,
		})
	}

	return view, err
}
