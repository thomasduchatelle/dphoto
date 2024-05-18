// Package catalog provides tools to maintain an index of all medias that have been backed up.
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

	// FindMedias is a paginated search for media with their details
	FindMedias(ctx context.Context, request *FindMediaRequest) (medias []*MediaMeta, err error)
	// FindMediaIds is a paginated search to only get the media ids
	FindMediaIds(ctx context.Context, request *FindMediaRequest) (ids []MediaId, err error)
	// FindMediaCurrentAlbum returns the folderName the media is currently in
	FindMediaCurrentAlbum(ctx context.Context, owner ownermodel.Owner, mediaId MediaId) (id *AlbumId, err error)
	// FindExistingSignatures returns the signatures that are already known
	FindExistingSignatures(ctx context.Context, owner ownermodel.Owner, signatures []*MediaSignature) ([]*MediaSignature, error)

	CountMedia(ctx context.Context, album ...AlbumId) (map[AlbumId]int, error)
}

// FindAllAlbums find all albums owned by root user
func FindAllAlbums(owner ownermodel.Owner) ([]*Album, error) {
	return repositoryPort.FindAlbumsByOwner(context.TODO(), owner)
}

// FindAlbums get several albums by their business keys
func FindAlbums(keys []AlbumId) ([]*Album, error) {
	return repositoryPort.FindAlbumByIds(context.TODO(), keys...)
}

// FindAlbum get an album by its business key (its folder name), or returns AlbumNotFoundError
func FindAlbum(id AlbumId) (*Album, error) {
	albums, err := repositoryPort.FindAlbumByIds(context.TODO(), id)
	if err != nil {
		return nil, err
	}
	if len(albums) == 0 {
		return nil, AlbumNotFoundError
	}
	return albums[0], nil
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
