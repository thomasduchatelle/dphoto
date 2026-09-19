package catalog

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

var (
	OrphanedMediasErr = errors.New("albums cannot be deleted or amended if it orphans medias.")
)

type CountMediasBySelectorsPort interface {
	CountMediasBySelectors(ctx context.Context, owner ownermodel.Owner, selectors []MediaSelector) (int, error)
}

type CountMediasBySelectorsFunc func(ctx context.Context, owner ownermodel.Owner, selectors []MediaSelector) (int, error)

func (f CountMediasBySelectorsFunc) CountMediasBySelectors(ctx context.Context, owner ownermodel.Owner, selectors []MediaSelector) (int, error) {
	return f(ctx, owner, selectors)
}

// AlbumDeleted is fired after an album has been deleted and its medias transferred to
// the surrounding albums: it carries the id of the album that has been removed and the
// medias that were actually moved to another album.
type AlbumDeleted struct {
	DeletedAlbumId    AlbumId
	TransferredMedias TransferredMedias
}

type AlbumDeletedObserver interface {
	OnAlbumDeleted(ctx context.Context, event AlbumDeleted) error
}

type AlbumDeletedObserverFunc func(ctx context.Context, event AlbumDeleted) error

func (f AlbumDeletedObserverFunc) OnAlbumDeleted(ctx context.Context, event AlbumDeleted) error {
	return f(ctx, event)
}

// NewDeleteAlbum creates the service to delete an album, transferring its medias to the
// surrounding albums when possible.
func NewDeleteAlbum(
	TimelineRepository TimelineRepository,
	CountMediasBySelectors CountMediasBySelectorsPort,
	TransferMedias TransferMediasService,
	AlbumDeletedObservers ...AlbumDeletedObserver,
) *DeleteAlbum {
	return &DeleteAlbum{
		TimelineRepository:     TimelineRepository,
		CountMediasBySelectors: CountMediasBySelectors,
		TransferMediasService:  TransferMedias,
		AlbumDeletedObservers:  AlbumDeletedObservers,
	}
}

// DeleteAlbum removes an album from the catalog: any media it currently holds is moved
// to the surrounding albums (as computed by the TimelineAggregate), and an AlbumDeleted
// event is fired once the row has been removed.
type DeleteAlbum struct {
	TimelineRepository     TimelineRepository
	CountMediasBySelectors CountMediasBySelectorsPort
	TransferMediasService  TransferMediasService
	AlbumDeletedObservers  []AlbumDeletedObserver
}

func (d *DeleteAlbum) DeleteAlbum(ctx context.Context, albumId AlbumId) error {
	timeline, err := d.TimelineRepository.LoadTimeline(ctx, albumId.Owner)
	if err != nil {
		return err
	}

	records, orphaned, err := timeline.RemoveAlbum(albumId)
	if err != nil {
		return err
	}

	if len(orphaned) > 0 {
		count, err := d.CountMediasBySelectors.CountMediasBySelectors(ctx, albumId.Owner, orphaned)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.Wrapf(OrphanedMediasErr, "%d medias from %s cannot be reallocated", count, albumId)
		}
	}

	transferred, err := d.TransferMediasService.TransferMedias(ctx, records)
	if err != nil {
		return err
	}

	if err = d.TimelineRepository.DeleteAlbum(ctx, albumId); err != nil {
		return err
	}

	log.WithField("Owner", albumId.Owner).Infof("Album %s deleted", albumId)

	event := AlbumDeleted{
		DeletedAlbumId:    albumId,
		TransferredMedias: transferred,
	}
	for _, observer := range d.AlbumDeletedObservers {
		if err = observer.OnAlbumDeleted(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// AlbumDeletedAsTimelineMutation adapts an AlbumDeletedObserver notification into a
// TimelineMutationObserver call, forwarding only the TransferredMedias carried by the
// event. This bridges the AlbumDeleted event to the observers that still listen on the
// legacy TimelineMutationObserver interface (archive relocator, view counters, ...).
type AlbumDeletedAsTimelineMutation struct {
	TimelineMutationObserver TimelineMutationObserver
}

func (a *AlbumDeletedAsTimelineMutation) OnAlbumDeleted(ctx context.Context, event AlbumDeleted) error {
	if event.TransferredMedias.IsEmpty() {
		return nil
	}
	return a.TimelineMutationObserver.OnTransferredMedias(ctx, event.TransferredMedias)
}
