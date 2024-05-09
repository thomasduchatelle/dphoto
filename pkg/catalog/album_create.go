package catalog

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	TransferMediasPort TransferMediasPort,
	TimelineMutationObservers ...TimelineMutationObserver,
) *CreateAlbum {
	return &CreateAlbum{
		FindAlbumsByOwnerPort: FindAlbumsByOwnerPort,
		Observers: []CreateAlbumObserver{
			&CreateAlbumExecutor{
				InsertAlbumPort: InsertAlbumPort,
			},
			&CreateAlbumMediaTransfer{
				MediaTransfer: MediaTransfer{
					TransferMedias:            TransferMediasPort,
					TimelineMutationObservers: TimelineMutationObservers,
				},
			},
		},
	}

}

// CreateAlbumRequest is a request to create a new album
type CreateAlbumRequest struct {
	Owner            Owner
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

type FindAlbumsByOwnerPort interface {
	FindAlbumsByOwner(ctx context.Context, owner Owner) ([]*Album, error)
}

type FindAlbumsByOwnerFunc func(ctx context.Context, owner Owner) ([]*Album, error)

func (f FindAlbumsByOwnerFunc) FindAlbumsByOwner(ctx context.Context, owner Owner) ([]*Album, error) {
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
	ObserveCreateAlbum(ctx context.Context, album Album, records MediaTransferRecords) error
}

type CreateAlbumObserverFunc func(ctx context.Context, album Album, records MediaTransferRecords) error

func (f CreateAlbumObserverFunc) ObserveCreateAlbum(ctx context.Context, album Album, records MediaTransferRecords) error {
	return f(ctx, album, records)
}

type CreateAlbum struct {
	FindAlbumsByOwnerPort FindAlbumsByOwnerPort
	Observers             []CreateAlbumObserver
}

// Create creates a new album
func (c *CreateAlbum) Create(ctx context.Context, request CreateAlbumRequest) error {
	if err := request.IsValid(); err != nil {
		return err
	}

	folderName := generateFolderName(request.Name, request.Start)
	if request.ForcedFolderName != "" {
		folderName = NewFolderName(request.ForcedFolderName)
	}

	createdAlbum := Album{
		AlbumId: AlbumId{
			Owner:      request.Owner,
			FolderName: folderName,
		},
		Name:  request.Name,
		Start: request.Start,
		End:   request.End,
	}

	albums, err := c.FindAlbumsByOwnerPort.FindAlbumsByOwner(ctx, request.Owner)
	if err != nil {
		return err
	}

	records, err := NewTimelineMutator().AddNew(albums, createdAlbum)
	if err != nil {
		return err
	}

	for _, observer := range c.Observers {
		err = observer.ObserveCreateAlbum(ctx, createdAlbum, records)
		if err != nil {
			return err
		}
	}

	log.Infof("Album created: %s [%s]", request, createdAlbum.AlbumId)
	return nil
}

type CreateAlbumExecutor struct {
	InsertAlbumPort InsertAlbumPort
}

func (c *CreateAlbumExecutor) ObserveCreateAlbum(ctx context.Context, album Album, records MediaTransferRecords) error {
	return c.InsertAlbumPort.InsertAlbum(ctx, album)
}

type CreateAlbumMediaTransfer struct {
	MediaTransfer
}

func (c *CreateAlbumMediaTransfer) ObserveCreateAlbum(ctx context.Context, album Album, records MediaTransferRecords) error {
	return c.MediaTransfer.Transfer(ctx, records)
}
