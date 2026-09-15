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

type FindAlbumsByIdsPort interface {
	FindAlbumsById(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error)
}

type FindAlbumsByIdsFunc func(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error)

func (f FindAlbumsByIdsFunc) FindAlbumsById(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error) {
	return f(ctx, ids)
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
