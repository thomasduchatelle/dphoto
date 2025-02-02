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

type CreateAlbumObserver interface {
	ObserveCreateAlbum(ctx context.Context, createdAlbum Album) error
}

type CreateAlbumObserverFunc func(ctx context.Context, createdAlbum Album) error

func (f CreateAlbumObserverFunc) ObserveCreateAlbum(ctx context.Context, createdAlbum Album) error {
	return f(ctx, createdAlbum)
}

// NewAlbumCreate creates the service to create a new album, including the transfer of medias
func NewAlbumCreate(
	FindAlbumsByOwnerPort FindAlbumsByOwnerPort,
	InsertAlbumPort InsertAlbumPort,
	TransferMediasPort TransferMediasRepositoryPort,
	TimelineMutationObservers ...TimelineMutationObserver,
) *CreateAlbum {
	return &CreateAlbum{
		FindAlbumsByOwnerPort: FindAlbumsByOwnerPort,
		CreateAlbumWithTimeline: &CreateAlbumStateless{
			Observers: []CreateAlbumObserverWithTimeline{
				&CreateAlbumObserverWrapper{CreateAlbumObserver: &CreateAlbumExecutor{
					InsertAlbumPort: InsertAlbumPort,
				}},
				&CreateAlbumMediaTransfer{
					MediaTransfer: &MediaTransferExecutor{
						TransferMediasRepository:  TransferMediasPort,
						TimelineMutationObservers: TimelineMutationObservers,
					},
				},
			},
		},
	}
}

type InsertAlbumPort interface {
	InsertAlbum(ctx context.Context, album Album) error
}

type InsertAlbumPortFunc func(ctx context.Context, album Album) error

func (f InsertAlbumPortFunc) InsertAlbum(ctx context.Context, album Album) error {
	return f(ctx, album)
}

type CreateAlbumObserverWithTimeline interface {
	ObserveCreateAlbum(ctx context.Context, timeline *TimelineAggregate, createdAlbum Album) error
}

type CreateAlbumObserverWrapper struct {
	CreateAlbumObserver
}

func (c *CreateAlbumObserverWrapper) ObserveCreateAlbum(ctx context.Context, _ *TimelineAggregate, createdAlbum Album) error {
	return c.CreateAlbumObserver.ObserveCreateAlbum(ctx, createdAlbum)
}

type CreateAlbumWithTimeline interface {
	Create(ctx context.Context, timeline *TimelineAggregate, request CreateAlbumRequest) (*AlbumId, error)
}

type CreateAlbum struct {
	FindAlbumsByOwnerPort   FindAlbumsByOwnerPort
	CreateAlbumWithTimeline CreateAlbumWithTimeline
}

// Create creates a new album
func (c *CreateAlbum) Create(ctx context.Context, request CreateAlbumRequest) (*AlbumId, error) {
	albums, err := c.FindAlbumsByOwnerPort.FindAlbumsByOwner(ctx, request.Owner)
	if err != nil {
		return nil, err
	}

	timeline := NewLazyTimelineAggregate(albums)

	return c.CreateAlbumWithTimeline.Create(ctx, timeline, request)
}

type CreateAlbumStateless struct {
	Observers []CreateAlbumObserverWithTimeline
}

func (c *CreateAlbumStateless) Create(ctx context.Context, timeline *TimelineAggregate, request CreateAlbumRequest) (*AlbumId, error) {
	album, err := timeline.CreateNewAlbum(request)
	if err != nil {
		return nil, err
	}

	for index, observer := range c.Observers {
		if err = observer.ObserveCreateAlbum(ctx, timeline, album); err != nil {
			return nil, errors.Wrapf(err, "CreateNewAlbum(%s) failed at observer %d/%d", request, index, len(c.Observers))
		}
	}

	log.WithField("Owner", request.Owner).Infof("Album %s created", album)

	return &album.AlbumId, nil
}

type CreateAlbumExecutor struct {
	InsertAlbumPort InsertAlbumPort
}

func (c *CreateAlbumExecutor) ObserveCreateAlbum(ctx context.Context, createdAlbum Album) error {
	return c.InsertAlbumPort.InsertAlbum(ctx, createdAlbum)
}

type CreateAlbumMediaTransfer struct {
	MediaTransfer MediaTransfer
}

func (c *CreateAlbumMediaTransfer) ObserveCreateAlbum(ctx context.Context, timeline *TimelineAggregate, createdAlbum Album) error {
	records, err := timeline.AddNew(createdAlbum)
	if err != nil {
		return err
	}

	return c.MediaTransfer.Transfer(ctx, records)
}
