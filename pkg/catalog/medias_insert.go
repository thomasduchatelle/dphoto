package catalog

import (
	"context"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

func NewInsertMedias(
	InsertMediasRepository InsertMediasRepositoryPort,
	InsertMediasObservers ...InsertMediasObserver,
) *InsertMedias {

	return &InsertMedias{
		InsertMediasRepository: InsertMediasRepository,
		InsertMediasObservers:  InsertMediasObservers,
	}
}

type InsertMediasObserver interface {
	OnMediasInserted(context.Context, map[AlbumId][]MediaId) error
}

// InsertMedias is a use case to pre-generate ids and store media metadata
type InsertMedias struct {
	InsertMediasRepository InsertMediasRepositoryPort
	InsertMediasObservers  []InsertMediasObserver
}

type InsertMediasRepositoryPort interface {
	// InsertMedias bulks insert medias
	InsertMedias(ctx context.Context, owner ownermodel.Owner, media []CreateMediaRequest) error
}

func (i *InsertMedias) Insert(ctx context.Context, owner ownermodel.Owner, medias []CreateMediaRequest) error {
	err := i.InsertMediasRepository.InsertMedias(ctx, owner, medias)
	if err != nil {
		return err
	}

	insertedMedias := make(map[AlbumId][]MediaId)
	for _, media := range medias {
		albumId := AlbumId{Owner: owner, FolderName: media.FolderName}
		if list, exists := insertedMedias[albumId]; exists {
			insertedMedias[albumId] = append(list, media.Id)
		} else {
			insertedMedias[albumId] = []MediaId{media.Id}
		}
	}

	for _, observer := range i.InsertMediasObservers {
		err = observer.OnMediasInserted(ctx, insertedMedias)
		if err != nil {
			return err
		}
	}

	return nil
}
