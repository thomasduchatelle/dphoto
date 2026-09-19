package catalog

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type FindAndRenameAlbumPort interface {
	FindAlbumsByOwnerPort
	InsertAlbumPort
	DeleteAlbumRepositoryPort
	UpdateAlbumName(ctx context.Context, albumId AlbumId, newName string) error
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
	FindAndRenameAlbumPort FindAndRenameAlbumPort,
	TransferMedias TransferMediasService,
	AlbumRenamedObservers ...AlbumRenamedObserver,
) *RenameAlbum {
	return &RenameAlbum{
		FindAndRenameAlbumPort: FindAndRenameAlbumPort,
		TransferMediasService:  TransferMedias,
		AlbumRenamedObservers:  AlbumRenamedObservers,
	}
}

// RenameAlbum applies a rename request. When only the display name changes (RenameFolder is
// false and no ForcedFolderName is provided), it is an in-place update of the album row.
// Otherwise, the album is replaced: a new row is inserted under the new folder name, medias
// are transferred, the old row is deleted, and the AlbumRenamed event is fired.
type RenameAlbum struct {
	FindAndRenameAlbumPort FindAndRenameAlbumPort
	TransferMediasService  TransferMediasService
	AlbumRenamedObservers  []AlbumRenamedObserver
}

func (r *RenameAlbum) RenameAlbum(ctx context.Context, request RenameAlbumRequest) error {
	if request.NewName != "" && !request.RenameFolder && request.ForcedFolderName == "" {
		return r.FindAndRenameAlbumPort.UpdateAlbumName(ctx, request.CurrentId, request.NewName)
	}

	return r.replaceAlbum(ctx, request)
}

func (r *RenameAlbum) replaceAlbum(ctx context.Context, request RenameAlbumRequest) error {
	albums, err := r.FindAndRenameAlbumPort.FindAlbumsByOwner(ctx, request.CurrentId.Owner)
	if err != nil {
		return err
	}

	timeline, err := NewInitialisedTimelineAggregate(albums)
	if err != nil {
		return err
	}

	event, err := timeline.RenameAlbum(request)
	if err != nil {
		return err
	}

	if err = r.FindAndRenameAlbumPort.InsertAlbum(ctx, event.RenamedAlbum); err != nil {
		return err
	}

	if _, err = r.TransferMediasService.TransferMedias(ctx, event.MediaTransfer); err != nil {
		return err
	}

	if err = r.FindAndRenameAlbumPort.DeleteAlbum(ctx, event.ExistingAlbum.AlbumId); err != nil {
		return err
	}

	log.WithField("AlbumId", request.CurrentId).Infof("Album renamed: %s", request.NewName)

	for _, observer := range r.AlbumRenamedObservers {
		if err = observer.OnAlbumRenamed(ctx, *event); err != nil {
			return err
		}
	}
	return nil
}
