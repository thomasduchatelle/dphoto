package catalog

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

var (
	OrphanedMediasError  = errors.New("albums cannot be deleted if it make medias orphaned")
	AlbumIsNotEmptyError = errors.New("album is not empty")
)

type CountMediasBySelectorsPort interface {
	CountMediasBySelectors(ctx context.Context, owner ownermodel.Owner, selectors []MediaSelector) (int, error)
}

type CountMediasBySelectorsFunc func(ctx context.Context, owner ownermodel.Owner, selectors []MediaSelector) (int, error)

func (f CountMediasBySelectorsFunc) CountMediasBySelectors(ctx context.Context, owner ownermodel.Owner, selectors []MediaSelector) (int, error) {
	return f(ctx, owner, selectors)
}

type DeleteAlbumRepositoryPort interface {
	DeleteAlbum(ctx context.Context, albumId AlbumId) error
}

type DeleteAlbumRepositoryFunc func(ctx context.Context, albumId AlbumId) error

func (f DeleteAlbumRepositoryFunc) DeleteAlbum(ctx context.Context, albumId AlbumId) error {
	return f(ctx, albumId)
}

type DeleteAlbumObserver interface {
	OnDeleteAlbum(ctx context.Context, deletedAlbum AlbumId, transfers MediaTransferRecords) error
}

// NewDeleteAlbum creates a new DeleteAlbum service.
func NewDeleteAlbum(
	FindAlbumsByOwner FindAlbumsByOwnerPort,
	CountMediasBySelectors CountMediasBySelectorsPort,
	TransferMediasPort TransferMediasRepositoryPort,
	DeleteAlbumRepository DeleteAlbumRepositoryPort,
	TimelineMutationObservers ...TimelineMutationObserver,
) *DeleteAlbum {

	return &DeleteAlbum{
		FindAlbumsByOwner:      FindAlbumsByOwner,
		CountMediasBySelectors: CountMediasBySelectors,
		Observers: []DeleteAlbumObserver{
			&DeleteAlbumMediaTransfer{
				MediaTransferExecutor: MediaTransferExecutor{
					TransferMediasRepository:  TransferMediasPort,
					TimelineMutationObservers: TimelineMutationObservers,
				},
			},
			&DeleteAlbumMetadata{
				DeleteAlbumRepository: DeleteAlbumRepository,
			},
		},
	}
}

type DeleteAlbum struct {
	FindAlbumsByOwner      FindAlbumsByOwnerPort
	CountMediasBySelectors CountMediasBySelectorsPort
	Observers              []DeleteAlbumObserver
}

// DeleteAlbum delete an album, medias it contains are dispatched to other albums.
func (d *DeleteAlbum) DeleteAlbum(ctx context.Context, albumId AlbumId) error {
	albums, err := d.FindAlbumsByOwner.FindAlbumsByOwner(ctx, albumId.Owner)
	if err != nil {
		return err
	}

	transfers, orphaned, err := NewLazyTimelineAggregate(albums).RemoveAlbum(albumId)
	if err != nil {
		return err
	}

	if len(orphaned) > 0 {
		count, err := d.CountMediasBySelectors.CountMediasBySelectors(ctx, albumId.Owner, orphaned)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.Wrapf(OrphanedMediasError, "%d orphaned medias prevent %s deletion", count, albumId)
		}
	}

	for _, observer := range d.Observers {
		err := observer.OnDeleteAlbum(ctx, albumId, transfers)
		if err != nil {
			return err
		}
	}

	log.Infof("Album deleted: %s", albumId)
	return nil
}

type DeleteAlbumMediaTransfer struct {
	MediaTransferExecutor
}

func (d *DeleteAlbumMediaTransfer) OnDeleteAlbum(ctx context.Context, deletedAlbum AlbumId, records MediaTransferRecords) error {
	return d.MediaTransferExecutor.Transfer(ctx, records)
}

type DeleteAlbumMetadata struct {
	DeleteAlbumRepository DeleteAlbumRepositoryPort
}

func (d *DeleteAlbumMetadata) OnDeleteAlbum(ctx context.Context, deletedAlbum AlbumId, transfers MediaTransferRecords) error {
	return d.DeleteAlbumRepository.DeleteAlbum(ctx, deletedAlbum)
}
