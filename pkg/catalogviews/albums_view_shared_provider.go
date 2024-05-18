package catalogviews

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type FindAlbumsByIdsPort interface {
	FindAlbumsById(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error)
}

type FindAlbumsByIdsFunc func(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error)

func (f FindAlbumsByIdsFunc) FindAlbumsById(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error) {
	return f(ctx, ids)
}

type SharedWithUserPort interface {
	ListAlbumIdsSharedWithUser(ctx context.Context, userId usermodel.UserId) ([]catalog.AlbumId, error)
}

type SharedWithUserFunc func(ctx context.Context, userId usermodel.UserId) ([]catalog.AlbumId, error)

func (f SharedWithUserFunc) ListAlbumIdsSharedWithUser(ctx context.Context, userId usermodel.UserId) ([]catalog.AlbumId, error) {
	return f(ctx, userId)
}

type SharedAlbumListProvider struct {
	FindAlbumsByIdsPort FindAlbumsByIdsPort
	SharedWithUserPort  SharedWithUserPort
	MediaCounterPort    MediaCounterPort
}

func (s *SharedAlbumListProvider) ListAlbums(ctx context.Context, user usermodel.CurrentUser, filter ListAlbumsFilter) ([]*VisibleAlbum, error) {
	if filter.OnlyDirectlyOwned {
		return nil, nil
	}

	shares, err := s.SharedWithUserPort.ListAlbumIdsSharedWithUser(ctx, user.UserId)
	if err != nil || len(shares) == 0 {
		return nil, err
	}

	countMedia, err := s.MediaCounterPort.CountMedia(ctx, shares...)
	if err != nil {
		return nil, err
	}

	albums, err := s.FindAlbumsByIdsPort.FindAlbumsById(ctx, shares)

	var view []*VisibleAlbum
	for _, album := range albums {
		count, _ := countMedia[album.AlbumId]
		view = append(view, &VisibleAlbum{
			Album:      *album,
			MediaCount: count,
		})
	}

	return view, err
}
