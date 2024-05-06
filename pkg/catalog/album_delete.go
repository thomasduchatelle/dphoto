package catalog

import (
	"context"
	"github.com/pkg/errors"
)

var (
	OrphanedMediasError  = errors.New("albums cannot be deleted if it make medias orphaned")
	AlbumIsNotEmptyError = errors.New("album is not empty")
)

type CountMediasBySelectorsPort interface {
	CountMediasBySelectors(ctx context.Context, owner Owner, selectors []MediaSelector) (int, error)
}

type CountMediasBySelectorsFunc func(ctx context.Context, owner Owner, selectors []MediaSelector) (int, error)

func (f CountMediasBySelectorsFunc) CountMediasBySelectors(ctx context.Context, owner Owner, selectors []MediaSelector) (int, error) {
	return f(ctx, owner, selectors)
}

type AlbumCanBeDeletedObserver interface {
	Observe(ctx context.Context, deletedAlbum AlbumId, transfers MediaTransferRecords) error
}

func NewDeleteAlbum(
	FindAlbumsByOwner FindAlbumsByOwnerPort,
	CountMediasBySelectors CountMediasBySelectorsPort,
	TransferMediasPort TransferMediasPort,
	DeleteAlbumRepository DeleteAlbumRepositoryPort,
	TimelineMutationObservers ...TimelineMutationObserver,
) *DeleteAlbum {
	return &DeleteAlbum{
		FindAlbumsByOwner:      FindAlbumsByOwner,
		CountMediasBySelectors: CountMediasBySelectors,
		AlbumCanBeDeletedObserver: []AlbumCanBeDeletedObserver{
			&DeleteAlbumMediaTransfer{
				TransferMedias:            TransferMediasPort,
				TimelineMutationObservers: TimelineMutationObservers,
			},
			&DeleteAlbumMetadata{
				DeleteAlbumRepository: DeleteAlbumRepository,
			},
		},
	}
}

type DeleteAlbum struct {
	FindAlbumsByOwner         FindAlbumsByOwnerPort
	CountMediasBySelectors    CountMediasBySelectorsPort
	AlbumCanBeDeletedObserver []AlbumCanBeDeletedObserver
}

// DeleteAlbum delete an album, medias it contains are dispatched to other albums.
func (d *DeleteAlbum) DeleteAlbum(ctx context.Context, albumId AlbumId) error {
	albums, err := d.FindAlbumsByOwner.FindAlbumsByOwner(ctx, albumId.Owner)
	if err != nil {
		return err
	}

	transfers, orphaned, err := NewTimelineMutator().RemoveAlbum(albums, albumId)
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

	for _, observer := range d.AlbumCanBeDeletedObserver {
		err := observer.Observe(ctx, albumId, transfers)
		if err != nil {
			return err
		}
	}

	return nil
}

type DeleteAlbumMediaTransfer struct {
	TransferMedias            TransferMediasPort
	TimelineMutationObservers []TimelineMutationObserver
}

func (d *DeleteAlbumMediaTransfer) Observe(ctx context.Context, deletedAlbum AlbumId, records MediaTransferRecords) error {
	transfers, err := d.TransferMedias.TransferMediasFromRecords(ctx, records)
	if err != nil || transfers.IsEmpty() {
		return err
	}

	for _, observer := range d.TimelineMutationObservers {
		err = observer.Observe(ctx, transfers)
		if err != nil {
			return err
		}
	}

	return nil
}

type DeleteAlbumRepositoryPort interface {
	DeleteAlbum(ctx context.Context, albumId AlbumId) error
}

type DeleteAlbumRepositoryFunc func(ctx context.Context, albumId AlbumId) error

func (f DeleteAlbumRepositoryFunc) DeleteAlbum(ctx context.Context, albumId AlbumId) error {
	return f(ctx, albumId)
}

type DeleteAlbumMetadata struct {
	DeleteAlbumRepository DeleteAlbumRepositoryPort
}

func (d *DeleteAlbumMetadata) Observe(ctx context.Context, deletedAlbum AlbumId, transfers MediaTransferRecords) error {
	return d.DeleteAlbumRepository.DeleteAlbum(ctx, deletedAlbum)
}
