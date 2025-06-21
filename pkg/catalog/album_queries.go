// Package state provides tools to maintain an index of all medias that have been backed up.
package catalog

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

var (
	repositoryPort RepositoryAdapter
)

// Init must be called before using this package.
func Init(repositoryAdapter RepositoryAdapter) {
	repositoryPort = repositoryAdapter
}

// RepositoryAdapter brings persistence layer to catalog package
type RepositoryAdapter interface {
	FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*Album, error)

	// FindAlbumByIds only returns found albums
	FindAlbumByIds(ctx context.Context, ids ...AlbumId) ([]*Album, error)

	CountMedia(ctx context.Context, album ...AlbumId) (map[AlbumId]int, error)
}

type AlbumQueries struct {
	Repository RepositoryAdapter
}

func (a *AlbumQueries) FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*Album, error) {
	return a.Repository.FindAlbumsByOwner(ctx, owner)
}

func (a *AlbumQueries) FindAlbumsById(ctx context.Context, ids []AlbumId) ([]*Album, error) {
	return a.Repository.FindAlbumByIds(ctx, ids...)
}

func (a *AlbumQueries) CountMedia(ctx context.Context, album ...AlbumId) (map[AlbumId]int, error) {
	return a.Repository.CountMedia(ctx, album...)
}

func (a *AlbumQueries) FindAlbum(ctx context.Context, albumId AlbumId) (*Album, error) {
	albums, err := a.Repository.FindAlbumByIds(ctx, albumId)
	if err != nil {
		return nil, err
	}

	if len(albums) == 0 {
		return nil, AlbumNotFoundErr
	}

	return albums[0], nil
}
