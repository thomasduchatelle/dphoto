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

type InsertAlbumPort interface {
	InsertAlbum(ctx context.Context, album Album) error
}

type InsertAlbumPortFunc func(ctx context.Context, album Album) error

func (f InsertAlbumPortFunc) InsertAlbum(ctx context.Context, album Album) error {
	return f(ctx, album)
}

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
	FindAlbumsByOwnerPort FindAlbumsByOwnerPort,
	InsertAlbumPort InsertAlbumPort,
	TransferMedias TransferMediasService,
	AlbumCreatedObservers ...AlbumCreatedObserver,
) *CreateAlbum {
	return &CreateAlbum{
		FindAlbumsByOwnerPort: FindAlbumsByOwnerPort,
		BulkCreateAlbum: &BulkCreateAlbum{
			InsertAlbumPort:       InsertAlbumPort,
			TransferMediasService: TransferMedias,
			AlbumCreatedObservers: AlbumCreatedObservers,
		},
	}
}

// CreateAlbum inserts a new album and, if it overlaps with existing albums, transfers the
// overlapping medias into it. On success, an AlbumCreated event is fired to every observer.
type CreateAlbum struct {
	FindAlbumsByOwnerPort FindAlbumsByOwnerPort
	BulkCreateAlbum       *BulkCreateAlbum
}

func (c *CreateAlbum) Create(ctx context.Context, request CreateAlbumRequest) (*AlbumId, error) {
	albums, err := c.FindAlbumsByOwnerPort.FindAlbumsByOwner(ctx, request.Owner)
	if err != nil {
		return nil, err
	}

	timeline := NewLazyTimelineAggregate(albums)

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
	InsertAlbumPort       InsertAlbumPort
	TransferMediasService TransferMediasService
	AlbumCreatedObservers []AlbumCreatedObserver
}

func (c *BulkCreateAlbum) Create(ctx context.Context, timeline *TimelineAggregate, request CreateAlbumRequest) (*AlbumId, error) {

	album, records, err := timeline.CreateNewAlbum(request)
	if err != nil {
		return nil, err
	}

	if err = c.InsertAlbumPort.InsertAlbum(ctx, album); err != nil {
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
