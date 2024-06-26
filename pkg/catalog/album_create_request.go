package catalog

import (
	"fmt"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
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
