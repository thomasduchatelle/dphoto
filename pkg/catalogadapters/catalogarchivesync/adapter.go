// Package catalogarchivesync is calling archive domain directly
package catalogarchivesync

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

type Observer struct {
}

func (a *Observer) Observe(ctx context.Context, transfers catalog.TransferredMedias) error {
	for targetAlbumId, ids := range transfers {
		convertedIds := make([]string, len(ids), len(ids))
		for i, id := range ids {
			convertedIds[i] = string(id)
		}

		err := archive.Relocate(targetAlbumId.Owner.String(), convertedIds, targetAlbumId.FolderName.String())
		if err != nil {
			return errors.Wrapf(err, "failed to relocate images to %s", targetAlbumId)
		}
	}

	return nil
}
