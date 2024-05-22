package catalog

import (
	"fmt"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"slices"
	"time"
)

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

func (c *CreateAlbumRequest) Convert(existingAlbums []*Album) (Album, error) {
	if err := c.IsValid(); err != nil {
		return Album{}, err
	}

	folderName := generateFolderName(c.Name, c.Start)
	if c.ForcedFolderName != "" {
		folderName = NewFolderName(c.ForcedFolderName)
	}

	albumId := AlbumId{
		Owner:      c.Owner,
		FolderName: folderName,
	}

	nameIsAlreadyTaken := slices.ContainsFunc(existingAlbums, func(album *Album) bool {
		return album.AlbumId.IsEqual(albumId)
	})
	if nameIsAlreadyTaken {
		return Album{}, AlbumFolderNameAlreadyTakenErr
	}

	return Album{
		AlbumId: albumId,
		Name:    c.Name,
		Start:   c.Start,
		End:     c.End,
	}, nil
}
