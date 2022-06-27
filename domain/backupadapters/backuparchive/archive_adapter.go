package backuparchive

import (
	"github.com/thomasduchatelle/dphoto/domain/archive"
	"github.com/thomasduchatelle/dphoto/domain/backup"
)

type adapter struct {
}

func New() backup.BArchiveAdapter {
	return new(adapter)
}

func (a *adapter) ArchiveMedia(owner string, media *backup.BackingUpMediaRequest) (string, error) {
	return archive.Store(&archive.StoreRequest{
		DateTime:         media.AnalysedMedia.Details.DateTime,
		FolderName:       media.FolderName,
		Id:               media.Id,
		Open:             media.AnalysedMedia.FoundMedia.ReadMedia,
		OriginalFilename: media.AnalysedMedia.FoundMedia.MediaPath().Filename,
		Owner:            owner,
		SignatureSha256:  media.AnalysedMedia.Sha256Hash,
	})
}
