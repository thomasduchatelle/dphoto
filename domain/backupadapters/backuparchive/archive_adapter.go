package backuparchive

import (
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"github.com/thomasduchatelle/dphoto/domain/backup"
)

type adapter struct {
}

func New() backup.ArchiveAdapter {
	return new(adapter)
}

func (a *adapter) ArchiveMedia(owner string, media *backup.BackingUpMediaRequest) (string, error) {
	reader, err := media.AnalysedMedia.FoundMedia.ReadMedia()
	if err != nil {
		return "", errors.Wrapf(err, "Can read the file %s to be stored", media.AnalysedMedia.FoundMedia)
	}

	return archive.Store(&archive.StoreRequest{
		Content:          reader,
		DateTime:         media.AnalysedMedia.Details.DateTime,
		FolderName:       media.FolderName,
		Id:               media.Id,
		OriginalFilename: media.AnalysedMedia.FoundMedia.MediaPath().Filename,
		Owner:            owner,
		Size:             uint(media.AnalysedMedia.FoundMedia.Size()),
	})
}
