package catalog

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

var (
	AlbumNameMandatoryErr            = errors.New("Album name is mandatory")
	AlbumStartAndEndDateMandatoryErr = errors.New("Start and End times are mandatory")
	AlbumEndDateMustBeAfterStartErr  = errors.New("Album end must be strictly after its start")
	AlbumFolderNameAlreadyTakenErr   = errors.New("Album folder name is already taken")
)

// NewAlbumCreate creates the service to create a new album, including the transfer of medias
func NewAlbumCreate(
	FindAlbumsByOwnerPort FindAlbumsByOwnerPort,
	InsertAlbumPort InsertAlbumPort,
	TransferMediasPort TransferMediasRepositoryPort,
	TimelineMutationObservers ...TimelineMutationObserver,
) *CreateAlbum {
	return &CreateAlbum{
		FindAlbumsByOwnerPort: FindAlbumsByOwnerPort,
		Observers: []func(timeline *TimelineAggregate) CreateAlbumObserver{
			func(_ *TimelineAggregate) CreateAlbumObserver {
				return &CreateAlbumExecutor{
					InsertAlbumPort: InsertAlbumPort,
				}
			},
			func(timeline *TimelineAggregate) CreateAlbumObserver {
				return &CreateAlbumMediaTransfer{
					Timeline: timeline,
					MediaTransfer: &MediaTransferExecutor{
						TransferMediasRepository:  TransferMediasPort,
						TimelineMutationObservers: TimelineMutationObservers,
					},
				}
			},
		},
	}
}

type FindAlbumsByOwnerPort interface {
	FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*Album, error)
}

type FindAlbumsByOwnerFunc func(ctx context.Context, owner ownermodel.Owner) ([]*Album, error)

func (f FindAlbumsByOwnerFunc) FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*Album, error) {
	return f(ctx, owner)
}

type InsertAlbumPort interface {
	InsertAlbum(ctx context.Context, album Album) error
}

type InsertAlbumPortFunc func(ctx context.Context, album Album) error

func (f InsertAlbumPortFunc) InsertAlbum(ctx context.Context, album Album) error {
	return f(ctx, album)
}

type CreateAlbumObserver interface {
	ObserveCreateAlbum(ctx context.Context, createdAlbum Album) error
}

type CreateAlbumObserverFunc func(ctx context.Context, createdAlbum Album) error

func (f CreateAlbumObserverFunc) ObserveCreateAlbum(ctx context.Context, createdAlbum Album) error {
	return f(ctx, createdAlbum)
}

type CreateAlbum struct {
	FindAlbumsByOwnerPort FindAlbumsByOwnerPort
	MediaTransfer         MediaTransfer
	Observers             []func(timeline *TimelineAggregate) CreateAlbumObserver
}

// Create creates a new album
func (c *CreateAlbum) Create(ctx context.Context, request CreateAlbumRequest) (*AlbumId, error) {
	albums, err := c.FindAlbumsByOwnerPort.FindAlbumsByOwner(ctx, request.Owner)
	if err != nil {
		return nil, err
	}

	timeline := NewLazyTimelineAggregate(albums)

	observers := make([]CreateAlbumObserver, len(c.Observers)+1)
	for index, factory := range c.Observers {
		observers[index] = factory(timeline)
	}
	observers[len(c.Observers)] = CreateAlbumObserverFunc(func(ctx context.Context, album Album) error {
		log.Infof("Album created: %s => %s", request, album.AlbumId)
		return nil
	})

	return timeline.CreateNewAlbum(ctx, request, observers...)
}

type CreateAlbumExecutor struct {
	InsertAlbumPort InsertAlbumPort
}

func (c *CreateAlbumExecutor) ObserveCreateAlbum(ctx context.Context, createdAlbum Album) error {
	return c.InsertAlbumPort.InsertAlbum(ctx, createdAlbum)
}

type CreateAlbumMediaTransfer struct {
	Timeline      *TimelineAggregate
	MediaTransfer MediaTransfer
}

func (c *CreateAlbumMediaTransfer) ObserveCreateAlbum(ctx context.Context, createdAlbum Album) error {
	records, err := c.Timeline.AddNew(createdAlbum)
	if err != nil {
		return err
	}

	return c.MediaTransfer.Transfer(ctx, records)
}
