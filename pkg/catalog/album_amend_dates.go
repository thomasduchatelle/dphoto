package catalog

import (
	"context"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type DatesUpdate struct {
	UpdatedAlbum  Album
	PreviousStart time.Time
	PreviousEnd   time.Time
}

func (a *DatesUpdate) DatesNotChanged() bool {
	return a.UpdatedAlbum.Start.Equal(a.PreviousStart) && a.UpdatedAlbum.End.Equal(a.PreviousEnd)
}

type AmendAlbumDateRepositoryPort interface {
	AmendDates(ctx context.Context, album AlbumId, start, end time.Time) error
}

// AlbumDatesAmended is fired after an album's dates have been persisted and the medias
// affected by the change have been transferred to their new albums.
type AlbumDatesAmended struct {
	DatesUpdate       DatesUpdate
	TransferredMedias TransferredMedias
}

type AlbumDatesAmendedObserver interface {
	OnAlbumDatesAmended(ctx context.Context, event AlbumDatesAmended) error
}

type AlbumDatesAmendedObserverFunc func(ctx context.Context, event AlbumDatesAmended) error

func (f AlbumDatesAmendedObserverFunc) OnAlbumDatesAmended(ctx context.Context, event AlbumDatesAmended) error {
	return f(ctx, event)
}

// NewAmendAlbumDates creates the service to amend the dates of an album.
func NewAmendAlbumDates(
	findAlbumsByOwner FindAlbumsByOwnerPort,
	countMediasBySelectors CountMediasBySelectorsPort,
	amendAlbumDateRepository AmendAlbumDateRepositoryPort,
	transferMedias TransferMediasService,
	observers ...AlbumDatesAmendedObserver,
) *AmendAlbumDates {
	return &AmendAlbumDates{
		FindAlbumsByOwnerPort:        findAlbumsByOwner,
		CountMediasBySelectorsPort:   countMediasBySelectors,
		AmendAlbumDateRepositoryPort: amendAlbumDateRepository,
		TransferMediasService:        transferMedias,
		AlbumDatesAmendedObservers:   observers,
	}
}

// AmendAlbumDates changes the start/end dates of an album, transfers medias affected by the
// date change to their new album, then fires an AlbumDatesAmended event. When some medias
// would be left orphan (no album covers them anymore) the operation is aborted and
// OrphanedMediasErr is returned before any change is persisted.
type AmendAlbumDates struct {
	FindAlbumsByOwnerPort        FindAlbumsByOwnerPort
	CountMediasBySelectorsPort   CountMediasBySelectorsPort
	AmendAlbumDateRepositoryPort AmendAlbumDateRepositoryPort
	TransferMediasService        TransferMediasService
	AlbumDatesAmendedObservers   []AlbumDatesAmendedObserver
}

func (a *AmendAlbumDates) AmendAlbumDates(ctx context.Context, albumId AlbumId, start, end time.Time) error {
	albums, err := a.FindAlbumsByOwnerPort.FindAlbumsByOwner(ctx, albumId.Owner)
	if err != nil {
		return err
	}

	timeline := NewLazyTimelineAggregate(albums)

	amendedAlbum, err := timeline.ValidateAmendDates(albumId, start, end)
	if err != nil {
		return err
	}

	if amendedAlbum.DatesNotChanged() {
		log.WithFields(log.Fields{
			"AlbumId": albumId,
			"Start":   start,
			"End":     end,
		}).Infof("Album %s dates haven't changed, nothing to do.", albumId)
		return nil
	}

	records, orphaned, err := timeline.AmendDates(*amendedAlbum)
	if err != nil {
		return err
	}

	if len(orphaned) > 0 {
		count, err := a.CountMediasBySelectorsPort.CountMediasBySelectors(ctx, amendedAlbum.UpdatedAlbum.Owner, orphaned)
		if err != nil {
			return err
		}
		if count > 0 {
			return errors.Wrapf(OrphanedMediasErr, "%d medias from %s cannot be reallocated to a different album", count, amendedAlbum.UpdatedAlbum.AlbumId)
		}
	}

	if err = a.AmendAlbumDateRepositoryPort.AmendDates(ctx, amendedAlbum.UpdatedAlbum.AlbumId, amendedAlbum.UpdatedAlbum.Start, amendedAlbum.UpdatedAlbum.End); err != nil {
		return err
	}

	transferred := NewTransferredMedias()
	if len(records) > 0 {
		transferred, err = a.TransferMediasService.TransferMedias(ctx, records)
		if err != nil {
			return err
		}
	}

	log.WithField("Owner", albumId.Owner).Infof("Album %s dates updates to %s -> %s", albumId, amendedAlbum.UpdatedAlbum.Start.Format(time.DateTime), amendedAlbum.UpdatedAlbum.End.Format(time.DateTime))

	event := AlbumDatesAmended{
		DatesUpdate:       *amendedAlbum,
		TransferredMedias: transferred,
	}
	for _, observer := range a.AlbumDatesAmendedObservers {
		if err = observer.OnAlbumDatesAmended(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// AlbumDatesAmendedAsTimelineMutation adapts an AlbumDatesAmendedObserver notification into
// a TimelineMutationObserver call, forwarding only the TransferredMedias carried by the
// event. This bridges the AlbumDatesAmended event to the observers that still listen on the
// legacy TimelineMutationObserver interface (archive relocator, view counters, ...).
type AlbumDatesAmendedAsTimelineMutation struct {
	TimelineMutationObserver TimelineMutationObserver
}

func (a *AlbumDatesAmendedAsTimelineMutation) OnAlbumDatesAmended(ctx context.Context, event AlbumDatesAmended) error {
	if event.TransferredMedias.IsEmpty() {
		return nil
	}
	return a.TimelineMutationObserver.OnTransferredMedias(ctx, event.TransferredMedias)
}
