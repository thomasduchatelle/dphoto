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
		Observers: []CreateAlbumObserver{
			&CreateAlbumExecutor{
				InsertAlbumPort: InsertAlbumPort,
			},
			&CreateAlbumMediaTransfer{
				MediaTransfer: &MediaTransferExecutor{
					TransferMediasRepository:  TransferMediasPort,
					TimelineMutationObservers: TimelineMutationObservers,
				},
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

type MoveMediaPort interface {
	MoveMedia(ctx context.Context, albumId AlbumId, mediaIds []MediaId) error
}

type MoveMediaPortFunc func(ctx context.Context, albumId AlbumId, mediaIds []MediaId) error

func (f MoveMediaPortFunc) MoveMedia(ctx context.Context, albumId AlbumId, mediaIds []MediaId) error {
	return f(ctx, albumId, mediaIds)
}

type CreateAlbumObserver interface {
	ObserveCreateAlbum(ctx context.Context, createdAlbum Album, existingAlbums []*Album) error
}

type CreateAlbumObserverFunc func(ctx context.Context, album Album, existingAlbums []*Album) error

func (f CreateAlbumObserverFunc) ObserveCreateAlbum(ctx context.Context, createdAlbum Album, existingAlbums []*Album) error {
	return f(ctx, createdAlbum, existingAlbums)
}

type MediaTransfer interface {
	Transfer(ctx context.Context, records MediaTransferRecords) error
}

type MediaTransferFunc func(ctx context.Context, records MediaTransferRecords) error

func (f MediaTransferFunc) Transfer(ctx context.Context, records MediaTransferRecords) error {
	return f(ctx, records)
}

type CreateAlbum struct {
	FindAlbumsByOwnerPort FindAlbumsByOwnerPort
	Observers             []CreateAlbumObserver
}

// Create creates a new album
func (c *CreateAlbum) Create(ctx context.Context, request CreateAlbumRequest) (*AlbumId, error) {
	albums, err := c.FindAlbumsByOwnerPort.FindAlbumsByOwner(ctx, request.Owner)
	if err != nil {
		return nil, err
	}

	createdAlbum, err := request.Convert(albums)
	if err != nil {
		return nil, err
	}

	for _, observer := range c.Observers {
		if err := observer.ObserveCreateAlbum(ctx, createdAlbum, albums); err != nil {
			return nil, err
		}
	}

	log.Infof("Album created: %s [%s]", request, createdAlbum.AlbumId)
	return &createdAlbum.AlbumId, nil
}

type CreateAlbumExecutor struct {
	InsertAlbumPort InsertAlbumPort
}

func (c *CreateAlbumExecutor) ObserveCreateAlbum(ctx context.Context, createdAlbum Album, existingAlbums []*Album) error {
	return c.InsertAlbumPort.InsertAlbum(ctx, createdAlbum)
}

type CreateAlbumMediaTransfer struct {
	MediaTransfer MediaTransfer
}

func (c *CreateAlbumMediaTransfer) ObserveCreateAlbum(ctx context.Context, createdAlbum Album, existingAlbums []*Album) error {
	records, err := NewTimelineAggregate(existingAlbums).AddNew(createdAlbum)
	if err != nil {
		return err
	}

	return c.MediaTransfer.Transfer(ctx, records)
}
