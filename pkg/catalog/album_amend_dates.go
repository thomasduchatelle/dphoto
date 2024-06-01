package catalog

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type DatesUpdate struct {
	UpdatedAlbum  Album
	PreviousStart time.Time
	PreviousEnd   time.Time
}

func (a *DatesUpdate) DatesNotChanged() bool {
	return a.UpdatedAlbum.Start.Equal(a.PreviousStart) && a.UpdatedAlbum.End.Equal(a.PreviousEnd)
}

type AlbumDatesAmendedObserver interface {
	OnAlbumDatesAmended(ctx context.Context, amendedAlbum DatesUpdate) error
}

type AlbumDatesAmendedObserverFunc func(ctx context.Context, amendedAlbum DatesUpdate) error

func (f AlbumDatesAmendedObserverFunc) OnAlbumDatesAmended(ctx context.Context, amendedAlbum DatesUpdate) error {
	return f(ctx, amendedAlbum)
}

func NewAmendAlbumDates(
	findAlbumsByOwner FindAlbumsByOwnerPort,
	countMediasBySelectors CountMediasBySelectorsPort,
	amendAlbumDateRepository AmendAlbumDateRepositoryPort,
	transferMedias TransferMediasRepositoryPort,
	timelineMutationObservers ...TimelineMutationObserver,
) *AmendAlbumDates {

	return &AmendAlbumDates{
		FindAlbumsByOwnerPort:       findAlbumsByOwner,
		AmendAlbumDatesWithTimeline: new(AmendAlbumDatesStateless),
		Observers: []AlbumDatesAmendedObserverWithTimeline{
			&AmendAlbumMediaTransfer{
				CountMediasBySelectors: countMediasBySelectors,
				MediaTransfer: &MediaTransferExecutor{
					TransferMediasRepository:  transferMedias,
					TimelineMutationObservers: timelineMutationObservers,
				},
			},
			&AlbumDatesAmendedObserverWrapper{AlbumDatesAmendedObserver: &AmendAlbumDatesExecutor{
				AmendAlbumDateRepository: amendAlbumDateRepository,
			}},
		},
	}
}

type AmendAlbumDatesWithTimeline interface {
	AmendAlbumDates(ctx context.Context, timeline *TimelineAggregate, albumId AlbumId, start, end time.Time, observers ...AlbumDatesAmendedObserver) error
}

type AlbumDatesAmendedObserverWithTimeline interface {
	OnAlbumDatesAmended(ctx context.Context, timeline *TimelineAggregate, amendedAlbum DatesUpdate) error
}

type AlbumDatesAmendedObserverWrapper struct {
	AlbumDatesAmendedObserver
}

func (a *AlbumDatesAmendedObserverWrapper) OnAlbumDatesAmended(ctx context.Context, _ *TimelineAggregate, amendedAlbum DatesUpdate) error {
	return a.AlbumDatesAmendedObserver.OnAlbumDatesAmended(ctx, amendedAlbum)
}

// AmendAlbumDates is building the TimelineAggregate and passing it to methods requiring it.
type AmendAlbumDates struct {
	FindAlbumsByOwnerPort       FindAlbumsByOwnerPort
	AmendAlbumDatesWithTimeline AmendAlbumDatesWithTimeline
	Observers                   []AlbumDatesAmendedObserverWithTimeline
}

func (a *AmendAlbumDates) AmendAlbumDates(ctx context.Context, albumId AlbumId, start, end time.Time) error {
	albums, err := a.FindAlbumsByOwnerPort.FindAlbumsByOwner(ctx, albumId.Owner)
	if err != nil {
		return err
	}

	timeline := NewLazyTimelineAggregate(albums)

	return a.AmendAlbumDatesWithTimeline.AmendAlbumDates(ctx, timeline, albumId, start, end, AlbumDatesAmendedObserverFunc(func(ctx context.Context, amendedAlbum DatesUpdate) error {
		for _, observer := range a.Observers {
			err := observer.OnAlbumDatesAmended(ctx, timeline, amendedAlbum)
			if err != nil {
				return err
			}
		}
		return nil
	}))
}

type AmendAlbumDatesStateless struct{}

func (a *AmendAlbumDatesStateless) AmendAlbumDates(ctx context.Context, timeline *TimelineAggregate, albumId AlbumId, start, end time.Time, observers ...AlbumDatesAmendedObserver) error {
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

	for _, observer := range observers {
		err = observer.OnAlbumDatesAmended(ctx, *amendedAlbum)
		if err != nil {
			return err
		}
	}

	return nil
}

type AmendAlbumMediaTransfer struct {
	CountMediasBySelectors CountMediasBySelectorsPort
	MediaTransfer          MediaTransfer
}

func (a *AmendAlbumMediaTransfer) OnAlbumDatesAmended(ctx context.Context, timeline *TimelineAggregate, updatedAlbum DatesUpdate) error {
	records, orphaned, err := timeline.AmendDates(updatedAlbum)
	if err != nil {
		return err
	}

	if len(orphaned) > 0 {
		count, err := a.CountMediasBySelectors.CountMediasBySelectors(ctx, updatedAlbum.UpdatedAlbum.Owner, orphaned)
		if err != nil {
			return err
		}
		if count > 0 {
			var orphanedDesc []string
			for _, o := range orphaned {
				orphanedDesc = append(orphanedDesc, o.String())
			}
			return errors.Wrapf(OrphanedMediasError, "%d medias belongs to %s and would be orphaned in the range %s ; aborting amending date operation.", count, updatedAlbum.UpdatedAlbum.AlbumId, strings.Join(orphanedDesc, ", "))
		}
	}

	if len(records) > 0 {
		err = a.MediaTransfer.Transfer(ctx, records)
		if err != nil {
			return err
		}
	}

	return nil
}

type AmendAlbumDateRepositoryPort interface {
	AmendDates(ctx context.Context, album AlbumId, start, end time.Time) error
}

type AmendAlbumDatesExecutor struct {
	AmendAlbumDateRepository AmendAlbumDateRepositoryPort
}

func (a *AmendAlbumDatesExecutor) OnAlbumDatesAmended(ctx context.Context, update DatesUpdate) error {
	return a.AmendAlbumDateRepository.AmendDates(ctx, update.UpdatedAlbum.AlbumId, update.UpdatedAlbum.Start, update.UpdatedAlbum.End)
}
