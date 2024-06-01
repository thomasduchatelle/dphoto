package catalog

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"strings"
	"time"
)

var (
	NoAlbumLookedUpError = errors.New("no album matching")
)

type AlbumReference struct {
	AlbumId          *AlbumId // AlbumId if no album fit the time and the implementation doesn't support creation.
	AlbumJustCreated bool     // AlbumJustCreated is true if the album was created during the reference process (depending on the implementation capability).
}

func NewAlbumAutoPopulateReferencer(
	owner ownermodel.Owner,
	findAlbumsByOwner FindAlbumsByOwnerPort,
	insertAlbumPort InsertAlbumPort,
	transferMediasPort TransferMediasRepositoryPort,
	timelineMutationObservers ...TimelineMutationObserver,
) (*StatefulAlbumReferencer, error) {

	albums, err := findAlbumsByOwner.FindAlbumsByOwner(context.Background(), owner)
	if err != nil {
		return nil, errors.Wrapf(err, "NewAlbumAutoPopulateReferencer(...) failed")
	}

	timeline, err := NewInitialisedTimelineAggregate(albums)

	return &StatefulAlbumReferencer{
		Owner:             owner,
		TimelineAggregate: timeline,
		LookupStrategies: []AlbumLookupStrategy{
			new(TimelineLookupStrategy),
			&AlbumAutoCreateLookupStrategy{
				Delegate: &CreateAlbumStateless{
					Observers: []CreateAlbumObserverWithTimeline{
						&CreateAlbumObserverWrapper{CreateAlbumObserver: &CreateAlbumExecutor{
							InsertAlbumPort: insertAlbumPort,
						}},
						&CreateAlbumMediaTransfer{
							MediaTransfer: &MediaTransferExecutor{
								TransferMediasRepository:  transferMediasPort,
								TimelineMutationObservers: timelineMutationObservers,
							},
						},
					},
				},
			},
		},
	}, errors.Wrapf(err, "NewAlbumAutoPopulateReferencer(...) failed")
}

func NewAlbumDryRunReferencer(
	owner ownermodel.Owner,
	findAlbumsByOwner FindAlbumsByOwnerPort,
) (*StatefulAlbumReferencer, error) {

	albums, err := findAlbumsByOwner.FindAlbumsByOwner(context.Background(), owner)
	if err != nil {
		return nil, errors.Wrapf(err, "NewAlbumDryRunReferencer(%s) failed", owner)
	}

	timeline, err := NewInitialisedTimelineAggregate(albums)

	return &StatefulAlbumReferencer{
		Owner:             owner,
		TimelineAggregate: timeline,
		LookupStrategies: []AlbumLookupStrategy{
			new(TimelineLookupStrategy),
			new(DryRunLookupStrategy),
		},
	}, errors.Wrapf(err, "NewAlbumDryRunReferencer(...) failed")
}

type AlbumLookupStrategy interface {
	// LookupAlbum returns the AlbumReference for the given mediaTime, or NoAlbumLookedUpError if it can't find any (or a technical error)
	LookupAlbum(ctx context.Context, owner ownermodel.Owner, timeline *TimelineAggregate, mediaTime time.Time) (AlbumReference, error)
}

type StatefulAlbumReferencer struct {
	Owner             ownermodel.Owner
	TimelineAggregate *TimelineAggregate
	LookupStrategies  []AlbumLookupStrategy
}

func (a *StatefulAlbumReferencer) FindReference(ctx context.Context, mediaTime time.Time) (AlbumReference, error) {
	for _, strategy := range a.LookupStrategies {
		albumReference, err := strategy.LookupAlbum(ctx, a.Owner, a.TimelineAggregate, mediaTime)
		if !errors.Is(err, NoAlbumLookedUpError) {
			return albumReference, err
		}
	}

	var strategies []string
	for _, strategy := range a.LookupStrategies {
		strategies = append(strategies, fmt.Sprintf("%T", strategy))
	}

	return AlbumReference{}, errors.Wrapf(NoAlbumLookedUpError, "no strategy found a matching album for %s with strategies %s", mediaTime, strings.Join(strategies, ", "))
}

type TimelineLookupStrategy struct{}

func (t TimelineLookupStrategy) LookupAlbum(ctx context.Context, owner ownermodel.Owner, timeline *TimelineAggregate, mediaTime time.Time) (AlbumReference, error) {
	album, exists, err := timeline.FindAt(mediaTime)
	if err != nil {
		return AlbumReference{}, err
	}
	if exists {
		return AlbumReference{
			AlbumId:          &album.AlbumId,
			AlbumJustCreated: false,
		}, nil
	}

	return AlbumReference{}, NoAlbumLookedUpError
}

type AlbumAutoCreateLookupStrategy struct {
	Delegate CreateAlbumWithTimeline
}

func (a *AlbumAutoCreateLookupStrategy) LookupAlbum(ctx context.Context, owner ownermodel.Owner, timeline *TimelineAggregate, mediaTime time.Time) (AlbumReference, error) {
	year := mediaTime.Year()
	quarter := (mediaTime.Month() - 1) / 3

	createRequest := CreateAlbumRequest{
		Owner:            owner,
		Name:             fmt.Sprintf("Q%d %d", quarter+1, year),
		Start:            time.Date(year, quarter*3+1, 1, 0, 0, 0, 0, time.UTC),
		End:              time.Date(year, (quarter+1)*3+1, 1, 0, 0, 0, 0, time.UTC),
		ForcedFolderName: fmt.Sprintf("/%d-Q%d", year, quarter+1),
	}

	albumId, err := a.Delegate.Create(ctx, timeline, createRequest)
	return AlbumReference{
		AlbumId:          albumId,
		AlbumJustCreated: true,
	}, err
}

type DryRunLookupStrategy struct{}

func (d *DryRunLookupStrategy) LookupAlbum(ctx context.Context, owner ownermodel.Owner, timeline *TimelineAggregate, mediaTime time.Time) (AlbumReference, error) {
	return AlbumReference{
		AlbumId:          &AlbumId{owner, NewFolderName("new-album")},
		AlbumJustCreated: true,
	}, nil
}
