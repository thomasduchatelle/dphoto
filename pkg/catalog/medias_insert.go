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

// AssignIdsToNewMedias filters out signatures that are already known and compute a unique ID for the others.
func (i *InsertMedias) AssignIdsToNewMedias(ctx context.Context, owner ownermodel.Owner, signatures []*MediaSignature) (map[MediaSignature]MediaId, error) {
	// TODO delete this method
	existingSignaturesSlice, err := FindSignatures(owner, signatures)
	if err != nil {
		return nil, err
	}

	existingSignatures := make(map[MediaSignature]interface{})
	for _, sign := range existingSignaturesSlice {
		existingSignatures[*sign] = nil
	}

	assignedIds := make(map[MediaSignature]MediaId)
	for _, sign := range signatures {
		if _, exists := existingSignatures[*sign]; !exists {
			assignedIds[*sign], err = GenerateMediaId(*sign)
			if err != nil {
				return nil, err
			}
		}
	}

	return assignedIds, nil
}
