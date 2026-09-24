// Package state provides tools to maintain an index of all medias that have been backed up.
package catalog

import (
	"context"
	"time"

	"github.com/pkg/errors"
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

// TimelinePersistencePort is the persistence surface a TimelineRepositoryAdapter needs
// from the outside world. It is satisfied by the DynamoDB Repository as well as by the
// in-memory fake used in tests.
type TimelinePersistencePort interface {
	FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*Album, error)
	InsertAlbum(ctx context.Context, album Album) error
	DeleteAlbum(ctx context.Context, albumId AlbumId) error
	UpdateAlbumName(ctx context.Context, albumId AlbumId, newName string) (Album, error)
	AmendDates(ctx context.Context, albumId AlbumId, start, end time.Time) error
}

// TimelineRepositoryAdapter wraps a persistence adapter (typically the DynamoDB one) to
// satisfy the TimelineRepository port used by the album use cases. LoadTimeline reads the
// owner's albums and builds a TimelineAggregate eagerly; the write methods are inherited
// from the embedded adapter.
type TimelineRepositoryAdapter struct {
	TimelinePersistencePort
}

func (t *TimelineRepositoryAdapter) LoadTimeline(ctx context.Context, owner ownermodel.Owner) (*TimelineAggregate, error) {
	albums, err := t.TimelinePersistencePort.FindAlbumsByOwner(ctx, owner)
	if err != nil {
		return nil, errors.Wrapf(err, "LoadTimeline(%s) failed", owner)
	}
	return NewTimelineAggregate(albums)
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
