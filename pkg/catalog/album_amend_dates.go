package catalog

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func NewAmendAlbumDates(
	findAlbumsByOwner FindAlbumsByOwnerPort,
	countMediasBySelectors CountMediasBySelectorsPort,
	amendAlbumDateRepository AmendAlbumDateRepositoryPort,
	transferMedias TransferMediasRepositoryPort,
	timelineMutationObservers ...TimelineMutationObserver,
) *AmendAlbumDates {

	return &AmendAlbumDates{
		FindAlbumsByOwnerPort: findAlbumsByOwner,
		AmendAlbumDatesObservers: []AmendAlbumDatesObserver{
			&AmendAlbumMediaTransfer{
				CountMediasBySelectors: countMediasBySelectors,
				MediaTransfer: &MediaTransferExecutor{
					TransferMediasRepository:  transferMedias,
					TimelineMutationObservers: timelineMutationObservers,
				},
			},
			&AmendAlbumDatesExecutor{
				AmendAlbumDateRepository: amendAlbumDateRepository,
			},
		},
	}
}

type AmendAlbumDatesObserver interface {
	OnAlbumDatesAmended(ctx context.Context, existingTimeline []*Album, updatedAlbum Album) error
}

type AmendAlbumDatesObserverFunc func(ctx context.Context, existingTimeline []*Album, updatedAlbum Album) error

func (f AmendAlbumDatesObserverFunc) OnAlbumDatesAmended(ctx context.Context, existingTimeline []*Album, updatedAlbum Album) error {
	return f(ctx, existingTimeline, updatedAlbum)
}

type AmendAlbumDates struct {
	FindAlbumsByOwnerPort    FindAlbumsByOwnerPort
	AmendAlbumDatesObservers []AmendAlbumDatesObserver
}

func (a *AmendAlbumDates) AmendAlbumDates(ctx context.Context, albumId AlbumId, start, end time.Time) error {
	currentTimeline, err := a.FindAlbumsByOwnerPort.FindAlbumsByOwner(ctx, albumId.Owner)
	if err != nil {
		return err
	}

	var album Album
	for _, it := range currentTimeline {
		if it.AlbumId.IsEqual(albumId) {
			album = *it
			break
		}
	}
	if !album.AlbumId.IsEqual(albumId) {
		return AlbumNotFoundError
	}

	if album.Start.Equal(start) && album.End.Equal(end) {
		log.WithFields(log.Fields{
			"AlbumId": albumId,
			"Start":   start,
			"End":     end,
		}).Infoln("Album date unchanged, nothing to do.")
		return nil
	}

	album.Start = start
	album.End = end

	for _, observer := range a.AmendAlbumDatesObservers {
		err = observer.OnAlbumDatesAmended(ctx, currentTimeline, album)
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

func (a *AmendAlbumMediaTransfer) OnAlbumDatesAmended(ctx context.Context, existingTimeline []*Album, updatedAlbum Album) error {
	timeline, err := NewInitialisedTimelineAggregate(existingTimeline)
	if err != nil {
		return errors.Wrapf(err, "OnAlbumDatesAmended(%s) failed", updatedAlbum)
	}

	records, orphaned, err := timeline.AmendDates(updatedAlbum)
	if err != nil {
		return err
	}

	if len(orphaned) > 0 {
		count, err := a.CountMediasBySelectors.CountMediasBySelectors(ctx, updatedAlbum.Owner, orphaned)
		if err != nil {
			return err
		}
		if count > 0 {
			var orphanedDesc []string
			for _, o := range orphaned {
				orphanedDesc = append(orphanedDesc, o.String())
			}
			return errors.Wrapf(OrphanedMediasError, "%d medias belongs to %s and would be orphaned in the range %s ; aborting amending date operation.", count, updatedAlbum.AlbumId, strings.Join(orphanedDesc, ", "))
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

func (a *AmendAlbumDatesExecutor) OnAlbumDatesAmended(ctx context.Context, _ []*Album, updatedAlbum Album) error {
	return a.AmendAlbumDateRepository.AmendDates(ctx, updatedAlbum.AlbumId, updatedAlbum.Start, updatedAlbum.End)
}
