package catalog

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type AlbumRenamed struct {
	ExistingAlbum     Album
	RenamedAlbum      Album
	TransferredMedias TransferredMedias
}

type AlbumRenamedObserver interface {
	OnAlbumRenamed(ctx context.Context, event AlbumRenamed) error
}

type AlbumRenamedObserverFunc func(ctx context.Context, event AlbumRenamed) error

func (f AlbumRenamedObserverFunc) OnAlbumRenamed(ctx context.Context, event AlbumRenamed) error {
	return f(ctx, event)
}

// NewRenameAlbum creates the service to rename an album.
func NewRenameAlbum(
	TimelineRepository TimelineRepository,
	TransferMedias TransferMediasService,
	AlbumRenamedObservers ...AlbumRenamedObserver,
) *RenameAlbum {
	return &RenameAlbum{
		TimelineRepository:    TimelineRepository,
		TransferMediasService: TransferMedias,
		AlbumRenamedObservers: AlbumRenamedObservers,
	}
}

type RenameAlbum struct {
	TimelineRepository    TimelineRepository
	TransferMediasService TransferMediasService
	AlbumRenamedObservers []AlbumRenamedObserver
}

func (r *RenameAlbum) RenameAlbum(ctx context.Context, request RenameAlbumRequest) error {
	if request.NewName != "" && !request.RenameFolder && request.ForcedFolderName == "" {
		return r.renameInPlace(ctx, request)
	}

	return r.replaceAlbum(ctx, request)
}

func (r *RenameAlbum) renameInPlace(ctx context.Context, request RenameAlbumRequest) error {
	previous, renamed, err := r.TimelineRepository.UpdateAlbumName(ctx, request.CurrentId, request.NewName)
	if err != nil {
		return err
	}

	event := AlbumRenamed{
		ExistingAlbum: previous,
		RenamedAlbum:  renamed,
	}
	for _, observer := range r.AlbumRenamedObservers {
		if err = observer.OnAlbumRenamed(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

func (r *RenameAlbum) replaceAlbum(ctx context.Context, request RenameAlbumRequest) error {
	timeline, err := r.TimelineRepository.LoadTimeline(ctx, request.CurrentId.Owner)
	if err != nil {
		return err
	}

	nameUpdate, err := timeline.RenameAlbum(request)
	if err != nil {
		return err
	}

	if err = r.TimelineRepository.InsertAlbum(ctx, nameUpdate.RenamedAlbum); err != nil {
		return err
	}

	transferred, err := r.TransferMediasService.TransferMedias(ctx, nameUpdate.MediaTransfer)
	if err != nil {
		return err
	}

	if err = r.TimelineRepository.DeleteAlbum(ctx, nameUpdate.ExistingAlbum.AlbumId); err != nil {
		return err
	}

	log.WithField("AlbumId", request.CurrentId).Infof("Album renamed: %s", request.NewName)

	event := AlbumRenamed{
		ExistingAlbum:     nameUpdate.ExistingAlbum,
		RenamedAlbum:      nameUpdate.RenamedAlbum,
		TransferredMedias: transferred,
	}
	for _, observer := range r.AlbumRenamedObservers {
		if err = observer.OnAlbumRenamed(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// AlbumRenamedAsTimelineMutation adapts an AlbumRenamedObserver notification into a
// TimelineMutationObserver call, forwarding only the TransferredMedias carried by the event.
// This bridges the AlbumRenamed event to the observers that still listen on the legacy
// TimelineMutationObserver interface (archive relocator, view counters, ...).
type AlbumRenamedAsTimelineMutation struct {
	TimelineMutationObserver TimelineMutationObserver
}

func (a *AlbumRenamedAsTimelineMutation) OnAlbumRenamed(ctx context.Context, event AlbumRenamed) error {
	if event.TransferredMedias.IsEmpty() {
		return nil
	}
	return a.TimelineMutationObserver.OnTransferredMedias(ctx, event.TransferredMedias)
}
