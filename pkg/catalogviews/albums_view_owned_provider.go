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

type OwnedAlbumListProvider struct {
	FindAlbumByOwnerPort    FindAlbumByOwnerPort
	GetAlbumSharingGridPort GetAlbumSharingGridPort
}

func (o *OwnedAlbumListProvider) ListAlbums(ctx context.Context, user usermodel.CurrentUser, filter ListAlbumsFilter) ([]*VisibleAlbum, error) {
	if user.Owner == nil {
		return nil, nil
	}

	ownedAlbums, err := o.FindAlbumByOwnerPort.FindAlbumsByOwner(ctx, *user.Owner)
	if err != nil {
		return nil, err
	}

	sharing, err := o.GetAlbumSharingGridPort.GetAlbumSharingGrid(ctx, *user.Owner)

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
