package catalog

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"time"
)

var (
	AlbumNameMandatoryErr            = errors.New("Album name is mandatory")
	AlbumStartAndEndDateMandatoryErr = errors.New("Start and End times are mandatory")
	AlbumEndDateMustBeAfterStartErr  = errors.New("Album end must be strictly after its start")
)

// NewAlbumCreate creates the service to create a new album, including the transfer of medias
func NewAlbumCreate(
	FindAlbumsByOwnerPort FindAlbumsByOwnerPort,
	InsertAlbumPort InsertAlbumPort,
	TransferMediasPort TransferMediasRepositoryPort,
	TimelineMutationObservers ...TimelineMutationObserver,
) *CreateAlbum {
	return &CreateAlbum{
		Observers: []CreateAlbumObserver{
			&CreateAlbumExecutor{
				InsertAlbumPort: InsertAlbumPort,
			},
			&CreateAlbumMediaTransfer{
				FindAlbumsByOwnerPort: FindAlbumsByOwnerPort, // FIXME albums already inserted ; it causes duplicates in the timeline
				MediaTransfer: &MediaTransferExecutor{
					TransferMediasRepository:  TransferMediasPort,
					TimelineMutationObservers: TimelineMutationObservers,
				},
			},
		},
	}
}

// CreateAlbumRequest is a request to create a new album
type CreateAlbumRequest struct {
	Owner            ownermodel.Owner
	Name             string
	Start            time.Time
	End              time.Time
	ForcedFolderName string
}

func (c *CreateAlbumRequest) String() string {
	const layout = "2006-01-02T03"
	return fmt.Sprintf("[%s -> %s] %s (%s/%s)", c.Start.Format(layout), c.End.Format(layout), c.Name, c.Owner, c.ForcedFolderName)
}

func (c *CreateAlbumRequest) IsValid() error {
	if err := c.Owner.IsValid(); err != nil {
		return err
	}
	if c.Name == "" {
		return AlbumNameMandatoryErr
	}

	if c.Start.IsZero() || c.End.IsZero() {
		return AlbumStartAndEndDateMandatoryErr
	}

	if !c.End.After(c.Start) {
		return AlbumEndDateMustBeAfterStartErr
	}

	return nil
}

type CreateAlbum struct {
	Observers []CreateAlbumObserver
}

func (c *CreateAlbum) Create(ctx context.Context, request CreateAlbumRequest) (*AlbumId, error) {
	validator := CreateAlbumValidator{}
	album, err := validator.Create(ctx, request)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create album %s, invalid request", request)
	}

	for _, observer := range c.Observers {
		if err := observer.ObserveCreateAlbum(ctx, album); err != nil {
			return nil, errors.Wrapf(err, "cannot create album %s, failed observer", request)
		}
	}

	return &album.AlbumId, nil
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
	ObserveCreateAlbum(ctx context.Context, album Album) error
}

type CreateAlbumObserverFunc func(ctx context.Context, album Album) error

func (f CreateAlbumObserverFunc) ObserveCreateAlbum(ctx context.Context, album Album) error {
	return f(ctx, album)
}

type MediaTransfer interface {
	Transfer(ctx context.Context, records MediaTransferRecords) error
}

type MediaTransferFunc func(ctx context.Context, records MediaTransferRecords) error

func (f MediaTransferFunc) Transfer(ctx context.Context, records MediaTransferRecords) error {
	return f(ctx, records)
}

type CreateAlbumExecutor struct {
	InsertAlbumPort InsertAlbumPort
}

func (c *CreateAlbumExecutor) ObserveCreateAlbum(ctx context.Context, album Album) error {
	return c.InsertAlbumPort.InsertAlbum(ctx, album)
}

type CreateAlbumMediaTransfer struct {
	MediaTransfer         MediaTransfer
	FindAlbumsByOwnerPort FindAlbumsByOwnerPort
}

func (c *CreateAlbumMediaTransfer) ObserveCreateAlbum(ctx context.Context, createdAlbum Album) error {
	albums, err := c.FindAlbumsByOwnerPort.FindAlbumsByOwner(ctx, createdAlbum.AlbumId.Owner)
	if err != nil {
		return err
	}

	records, err := NewTimelineAggregate(albums).AddNew(createdAlbum)
	if err != nil {
		return err
	}

	return c.MediaTransfer.Transfer(ctx, records)
}
