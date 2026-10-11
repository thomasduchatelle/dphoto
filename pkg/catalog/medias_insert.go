package catalog

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

func NewInsertMedias(
	InsertMediasRepository InsertMediasRepositoryPort,
	CoverService CoverServicePort,
	InsertMediasObservers ...InsertMediasObserver,
) *InsertMedias {

	return &InsertMedias{
		InsertMediasRepository: InsertMediasRepository,
		CoverService:           CoverService,
		InsertMediasObservers:  InsertMediasObservers,
	}
}

// MediasInserted is the lifecycle event fired after medias have been persisted. It
// carries the inserted MediaIds per album and the resulting cover sets of the albums
// whose covers were updated (albums whose cover set did not change are omitted).
type MediasInserted struct {
	Inserted map[AlbumId][]MediaId
	Covers   map[AlbumId][]Cover
}

type InsertMediasObserver interface {
	OnMediasInserted(ctx context.Context, event MediasInserted) error
}

type InsertMediasObserverFunc func(ctx context.Context, event MediasInserted) error

func (f InsertMediasObserverFunc) OnMediasInserted(ctx context.Context, event MediasInserted) error {
	return f(ctx, event)
}

// InsertMedias is a use case to pre-generate ids and store media metadata.
type InsertMedias struct {
	InsertMediasRepository InsertMediasRepositoryPort
	CoverService           CoverServicePort
	InsertMediasObservers  []InsertMediasObserver
}

type InsertMediasRepositoryPort interface {
	// InsertMedias bulks insert medias
	InsertMedias(ctx context.Context, owner ownermodel.Owner, media []CreateMediaRequest) error
}

func (i *InsertMedias) Insert(ctx context.Context, owner ownermodel.Owner, medias []CreateMediaRequest) error {
	if err := i.InsertMediasRepository.InsertMedias(ctx, owner, medias); err != nil {
		return err
	}

	insertedMedias := make(map[AlbumId][]MediaId)
	for _, media := range medias {
		albumId := AlbumId{Owner: owner, FolderName: media.FolderName}
		insertedMedias[albumId] = append(insertedMedias[albumId], media.Id)
	}

	affectedAlbums := make([]AlbumId, 0, len(insertedMedias))
	for albumId := range insertedMedias {
		affectedAlbums = append(affectedAlbums, albumId)
	}
	coversPerAlbum, err := i.CoverService.ForcedRandomise(ctx, affectedAlbums...)
	if err != nil {
		return errors.Wrapf(err, "InsertMedias failed to randomise covers of %v", affectedAlbums)
	}

	event := MediasInserted{
		Inserted: insertedMedias,
		Covers:   coversPerAlbum,
	}
	for _, observer := range i.InsertMediasObservers {
		if err := observer.OnMediasInserted(ctx, event); err != nil {
			return err
		}
	}

	return nil
}
