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
		InsertAlbumPort:       InsertAlbumPort,
		TransferMediasService: TransferMedias,
		AlbumCreatedObservers: AlbumCreatedObservers,
	}
}

// CreateAlbum inserts a new album and, if it overlaps with existing albums, transfers the
// overlapping medias into it. On success, an AlbumCreated event is fired to every observer.
type CreateAlbum struct {
	FindAlbumsByOwnerPort FindAlbumsByOwnerPort
	InsertAlbumPort       InsertAlbumPort
	TransferMediasService TransferMediasService
	AlbumCreatedObservers []AlbumCreatedObserver
}

func (c *CreateAlbum) Create(ctx context.Context, request CreateAlbumRequest) (*AlbumId, error) {
	albums, err := c.FindAlbumsByOwnerPort.FindAlbumsByOwner(ctx, request.Owner)
	if err != nil {
		return nil, err
	}

	timeline := NewLazyTimelineAggregate(albums)

	album, err := timeline.CreateNewAlbum(request)
	if err != nil {
		return nil, err
	}

	if err = c.InsertAlbumPort.InsertAlbum(ctx, album); err != nil {
		return nil, err
	}

	records, err := timeline.AddNew(album)
	if err != nil {
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

// AlbumCreatedAsTimelineMutation adapts an AlbumCreatedObserver notification into a
// TimelineMutationObserver call, forwarding only the TransferredMedias carried by the event.
// This bridges the AlbumCreated event to the observers that still listen on the legacy
// TimelineMutationObserver interface (archive relocator, view counters, ...).
type AlbumCreatedAsTimelineMutation struct {
	TimelineMutationObserver TimelineMutationObserver
}

func (a *AlbumCreatedAsTimelineMutation) OnAlbumCreated(ctx context.Context, event AlbumCreated) error {
	if event.TransferredMedias.IsEmpty() {
		return nil
	}
	return a.TimelineMutationObserver.OnTransferredMedias(ctx, event.TransferredMedias)
}

// The following types are kept alive solely so album_referencer.go still compiles. They
// will be removed once the referencer is refactored to build albums inline like CreateAlbum.
// TODO remove once album_referencer is refactored

type CreateAlbumObserver interface {
	ObserveCreateAlbum(ctx context.Context, createdAlbum Album) error
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
