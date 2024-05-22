package catalog

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type CreateAlbumValidator struct{}

// Create creates a new album
func (c *CreateAlbumValidator) Create(ctx context.Context, request CreateAlbumRequest) (Album, error) {
	if err := request.IsValid(); err != nil {
		return Album{}, err
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

	log.Infof("Album created: %s [%s]", request, createdAlbum.AlbumId)
	return createdAlbum, nil
}
