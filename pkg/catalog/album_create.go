package catalog

import (
	"context"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	AlbumNameMandatoryErr            = errors.New("Album name is mandatory")
	AlbumStartAndEndDateMandatoryErr = errors.New("Start and End times are mandatory")
	AlbumEndDateMustBeAfterStartErr  = errors.New("Album end must be strictly after its start")
	AlbumFolderNameAlreadyTakenErr   = errors.New("Album folder name is already taken")
)

// AlbumCreated is fired after a new album has been persisted and its overlapping medias
// have been transferred: it carries the newly created album and the medias that were
// actually moved into it (if any).
type AlbumCreated struct {
	CreatedAlbum      Album
	TransferredMedias TransferredMedias
}

type AlbumCreatedObserver interface {
	OnAlbumCreated(ctx context.Context, event AlbumCreated) error
}

type AlbumCreatedObserverFunc func(ctx context.Context, event AlbumCreated) error

func (f AlbumCreatedObserverFunc) OnAlbumCreated(ctx context.Context, event AlbumCreated) error {
	return f(ctx, event)
}

// NewAlbumCreate creates the service to create a new album, including the transfer of medias.
func NewAlbumCreate(
	TimelineRepository TimelineRepository,
	TransferMedias TransferMediasService,
	AlbumCreatedObservers ...AlbumCreatedObserver,
) *CreateAlbum {
	return &CreateAlbum{
		TimelineRepository: TimelineRepository,
		BulkCreateAlbum: &BulkCreateAlbum{
			TimelineRepository:    TimelineRepository,
			TransferMediasService: TransferMedias,
			AlbumCreatedObservers: AlbumCreatedObservers,
		},
	}
}

// CreateAlbum inserts a new album and, if it overlaps with existing albums, transfers the
// overlapping medias into it. On success, an AlbumCreated event is fired to every observer.
type CreateAlbum struct {
	TimelineRepository TimelineRepository
	BulkCreateAlbum    *BulkCreateAlbum
}

func (c *CreateAlbum) Create(ctx context.Context, request CreateAlbumRequest) (*AlbumId, error) {
	timeline, err := c.TimelineRepository.LoadTimeline(ctx, request.Owner)
	if err != nil {
		return nil, err
	}

	return c.BulkCreateAlbum.Create(ctx, timeline, request)
}

type AlbumCreatedAsTimelineMutation struct {
	TimelineMutationObserver TimelineMutationObserver
}

func (a *AlbumCreatedAsTimelineMutation) OnAlbumCreated(ctx context.Context, event AlbumCreated) error {
	if event.TransferredMedias.IsEmpty() {
		return nil
	}
	return a.TimelineMutationObserver.OnTransferredMedias(ctx, event.TransferredMedias)
}

type BulkCreateAlbum struct {
	TimelineRepository    TimelineRepository
	TransferMediasService TransferMediasService
	AlbumCreatedObservers []AlbumCreatedObserver
}

func (c *BulkCreateAlbum) Create(ctx context.Context, timeline *TimelineAggregate, request CreateAlbumRequest) (*AlbumId, error) {

	album, records, err := timeline.CreateNewAlbum(request)
	if err != nil {
		return nil, err
	}

	if err = c.TimelineRepository.InsertAlbum(ctx, album); err != nil {
		return nil, err
	}

	transferred, err := c.TransferMediasService.TransferMedias(ctx, records)
	if err != nil {
		return nil, err
	}

	log.WithField("Owner", request.Owner).Infof("Album %s created", album)

	event := AlbumCreated{
		CreatedAlbum:      album,
		TransferredMedias: transferred,
	}
	for _, observer := range c.AlbumCreatedObservers {
		if err = observer.OnAlbumCreated(ctx, event); err != nil {
			return nil, err
		}
	}

	return &album.AlbumId, nil
}
