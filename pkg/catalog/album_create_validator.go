package catalog

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type CreateAlbumValidator struct {
	Observers []CreateAlbumObserver
}

// Create creates a new album
func (c *CreateAlbumValidator) Create(ctx context.Context, request CreateAlbumRequest) (*AlbumId, error) {
	if err := request.IsValid(); err != nil {
		return nil, err
	}

	folderName := generateFolderName(request.Name, request.Start)
	if request.ForcedFolderName != "" {
		folderName = NewFolderName(request.ForcedFolderName)
	}

	albumId := AlbumId{
		Owner:      request.Owner,
		FolderName: folderName,
	}
	createdAlbum := Album{
		AlbumId: albumId,
		Name:    request.Name,
		Start:   request.Start,
		End:     request.End,
	}

	for _, observer := range c.Observers {
		err := observer.ObserveCreateAlbum(ctx, createdAlbum)
		if err != nil {
			return nil, err
		}
	}

	log.Infof("Album created: %s [%s]", request, createdAlbum.AlbumId)
	return &albumId, nil
}