package catalog

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type RenameAlbumRequest struct {
	CurrentId        AlbumId
	NewName          string
	RenameFolder     bool   // RenameFolder set to TRUE will create a new album with a FolderName generated from the NewName
	ForcedFolderName string // ForcedFolderName will create a new album with a specific FolderName (RenameFolder is ignored)
}

// NewRenameAlbum creates the service to rename an album
func NewRenameAlbum(
	FindAlbumById FindAlbumByIdPort,
	UpdateAlbumName UpdateAlbumNamePort,
	InsertAlbumPort InsertAlbumPort,
	DeleteAlbumRepositoryPort DeleteAlbumRepositoryPort,
	TransferMedias TransferMediasRepositoryPort,
	FindAlbumsByOwner FindAlbumsByOwnerPort,
	TimelineMutationObservers ...TimelineMutationObserver,
) *RenameAlbum {

	return &RenameAlbum{
		FindAlbumById:   FindAlbumById,
		UpdateAlbumName: UpdateAlbumName,
		RenameAlbumObservers: []RenameAlbumObserver{
			&RenameAlbumReplacer{
				CreateAlbum: CreateAlbum{
					FindAlbumsByOwnerPort: FindAlbumsByOwner,
					Observers: []CreateAlbumObserver{
						&CreateAlbumExecutor{
							InsertAlbumPort: InsertAlbumPort,
						},
					},
				},
				MediaTransfer: &MediaTransferExecutor{
					TransferMediasRepository:  TransferMedias,
					TimelineMutationObservers: TimelineMutationObservers,
				},
				DeleteAlbumRepositoryPort: DeleteAlbumRepositoryPort,
			},
		},
	}
}

func (r RenameAlbumRequest) String() string {
	return fmt.Sprintf("%s -> %s", r.CurrentId.String(), r.NewName)
}

func (r RenameAlbumRequest) IsValid() error {
	if r.NewName == "" {
		return AlbumNameMandatoryErr
	}

	return nil
}

type RenameAlbumObserver interface {
	OnRenameAlbum(ctx context.Context, current AlbumId, creationRequest CreateAlbumRequest) error
}

type FindAlbumByIdPort interface {
	FindAlbumById(ctx context.Context, id AlbumId) (*Album, error)
}

type FindAlbumByIdFunc func(ctx context.Context, id AlbumId) (*Album, error)

func (f FindAlbumByIdFunc) FindAlbumById(ctx context.Context, id AlbumId) (*Album, error) {
	return f(ctx, id)
}

type UpdateAlbumNamePort interface {
	UpdateAlbumName(ctx context.Context, albumId AlbumId, newName string) error
}

type RenameAlbum struct {
	FindAlbumById        FindAlbumByIdPort
	UpdateAlbumName      UpdateAlbumNamePort
	RenameAlbumObservers []RenameAlbumObserver
}

func (r *RenameAlbum) RenameAlbum(ctx context.Context, request RenameAlbumRequest) error {
	if err := request.IsValid(); err != nil {
		return err
	}
	existing, err := r.FindAlbumById.FindAlbumById(ctx, request.CurrentId)
	if err != nil {
		return err
	}

	if !request.RenameFolder && request.ForcedFolderName == "" {
		return r.UpdateAlbumName.UpdateAlbumName(ctx, request.CurrentId, request.NewName)
	}

	createRequest := CreateAlbumRequest{
		Owner:            request.CurrentId.Owner,
		Name:             request.NewName,
		Start:            existing.Start,
		End:              existing.End,
		ForcedFolderName: request.ForcedFolderName,
	}

	for _, observer := range r.RenameAlbumObservers {
		if err = observer.OnRenameAlbum(ctx, request.CurrentId, createRequest); err != nil {
			return err
		}
	}

	log.WithField("AlbumId", request.CurrentId).Infof("Album renamed: %s", request.NewName)
	return nil
}

type RenameAlbumReplacer struct {
	CreateAlbum               CreateAlbum
	MediaTransfer             MediaTransfer
	DeleteAlbumRepositoryPort DeleteAlbumRepositoryPort
}

func (r *RenameAlbumReplacer) OnRenameAlbum(ctx context.Context, current AlbumId, creationRequest CreateAlbumRequest) error {
	newAlbumId, err := r.CreateAlbum.Create(ctx, creationRequest)
	if err != nil {
		return err
	}

	records := MediaTransferRecords{
		*newAlbumId: []MediaSelector{
			{
				FromAlbums: []AlbumId{current},
				Start:      creationRequest.Start,
				End:        creationRequest.End,
			},
		},
	}
	err = r.MediaTransfer.Transfer(ctx, records)
	if err != nil {
		return err
	}

	return r.DeleteAlbumRepositoryPort.DeleteAlbum(ctx, current)
}
