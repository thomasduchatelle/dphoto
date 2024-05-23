package catalog

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"time"
)

func NewAlbumAutoPopulateReferencer(
	owner ownermodel.Owner,
	findAlbumsByOwner FindAlbumsByOwnerPort,
	insertAlbumPort InsertAlbumPort,
	transferMediasPort TransferMediasRepositoryPort,
	timelineMutationObservers ...TimelineMutationObserver,
) (*AutoCreateAlbumReferencer, error) {

	albums, err := findAlbumsByOwner.FindAlbumsByOwner(context.Background(), owner)
	if err != nil {
		return nil, errors.Wrapf(err, "NewAlbumAutoPopulateReferencer(...) failed")
	}

	timelineAggregate, err := NewInitialisedTimelineAggregate(albums)

	return &AutoCreateAlbumReferencer{
		owner:             owner,
		timelineAggregate: timelineAggregate,
		Observer: &StatefulCreateAlbum{
			Timeline: timelineAggregate,
			CreateAlbumObserver: []CreateAlbumObserver{
				&CreateAlbumExecutor{
					InsertAlbumPort: insertAlbumPort,
				},
				&CreateAlbumMediaTransfer{
					Timeline: timelineAggregate,
					MediaTransfer: &MediaTransferExecutor{
						TransferMediasRepository:  transferMediasPort,
						TimelineMutationObservers: timelineMutationObservers,
					},
				},
			},
		},
	}, errors.Wrapf(err, "NewAlbumAutoPopulateReferencer(...) failed")
}

type AlbumReference struct {
	AlbumId          *AlbumId // AlbumId if no album fit the time and the implementation doesn't support creation.
	AlbumJustCreated bool     // AlbumJustCreated is true if the album was created during the reference process (depending on the implementation capability).
}

type AutoCreateAlbumObserver interface {
	Create(ctx context.Context, request CreateAlbumRequest) (*AlbumId, error)
}

type AutoCreateAlbumReferencer struct {
	owner             ownermodel.Owner
	timelineAggregate *TimelineAggregate
	Observer          AutoCreateAlbumObserver
}

func (a *AutoCreateAlbumReferencer) FindReference(ctx context.Context, mediaTime time.Time) (AlbumReference, error) {
	album, exists, err := a.timelineAggregate.FindAt(mediaTime)
	if err != nil {
		return AlbumReference{}, err
	}
	if exists {
		return AlbumReference{
			AlbumId:          &album.AlbumId,
			AlbumJustCreated: false,
		}, nil
	}

	year := mediaTime.Year()
	quarter := (mediaTime.Month() - 1) / 3

	createRequest := CreateAlbumRequest{
		Owner:            a.owner,
		Name:             fmt.Sprintf("Q%d %d", quarter+1, year),
		Start:            time.Date(year, quarter*3+1, 1, 0, 0, 0, 0, time.UTC),
		End:              time.Date(year, (quarter+1)*3+1, 1, 0, 0, 0, 0, time.UTC),
		ForcedFolderName: fmt.Sprintf("/%d-Q%d", year, quarter+1),
	}

	albumId, err := a.Observer.Create(ctx, createRequest)
	return AlbumReference{
		AlbumId:          albumId,
		AlbumJustCreated: true,
	}, err
}

type StatefulCreateAlbum struct {
	Timeline            *TimelineAggregate
	CreateAlbumObserver []CreateAlbumObserver
}

func (c *StatefulCreateAlbum) Create(ctx context.Context, request CreateAlbumRequest) (*AlbumId, error) {
	newAlbum, err := c.Timeline.CreateNewAlbum(ctx, request, c.CreateAlbumObserver...)
	return newAlbum, errors.Wrapf(err, "StatefulCreateAlbum.Create(%s) failed", request)
}
