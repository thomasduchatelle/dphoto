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

type FindAlbumsByOwnerPort interface {
	FindAlbumsByOwner(ctx context.Context, owner Owner) ([]*Album, error)
}

type FindAlbumsByOwnerPortFunc func(ctx context.Context, owner Owner) ([]*Album, error)

func (f FindAlbumsByOwnerPortFunc) FindAlbumsByOwner(ctx context.Context, owner Owner) ([]*Album, error) {
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

type TransferredMedias map[AlbumId][]MediaId

type TransferMediasPort interface {
	TransferMediasFromRecords(ctx context.Context, records MediaTransferRecords) (TransferredMedias, error)
}

type CreateAlbum struct {
	FindAlbumsByOwnerPort     FindAlbumsByOwnerPort
	InsertAlbumPort           InsertAlbumPort
	TransferMediasPort        TransferMediasPort
	TimelineMutationObservers []TimelineMutationObserver
}

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
	mutator := NewTimelineMutator(albums)

	err = c.InsertAlbumPort.InsertAlbum(ctx, createdAlbum)
	if err != nil {
		return err
	}

	records, err := mutator.AddNew(createdAlbum)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		log.WithField("Album", createdAlbum.FolderName).
			Infof("Album %s has been created.", createdAlbum.FolderName)
		return nil
	}

	transfers, err := c.TransferMediasPort.TransferMediasFromRecords(ctx, records)
	count := 0

	if len(transfers) > 0 {
		for _, observer := range c.TimelineMutationObservers {
			err = observer.Observe(ctx, transfers)
			if err != nil {
				return err
			}
		}

		for _, ids := range transfers {
			count += len(ids)
		}
	}

	log.WithField("Album", createdAlbum.FolderName).
		Infof("Album %s has been created, %d medias have been transferred.", createdAlbum.FolderName, count)

	return err
}
