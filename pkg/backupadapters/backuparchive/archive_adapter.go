package backuparchive

import (
	"github.com/thomasduchatelle/dphoto/pkg/archive"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
)

type adapter struct {
}

func New() backup.ArchiveMediaPort {
	return new(adapter)
}

func (a *adapter) ArchiveMedia(owner string, media *backup.BackingUpMediaRequest) (string, error) {
	return archive.Store(&archive.StoreRequest{
		DateTime:         media.AnalysedMedia.Details.DateTime,
		FolderName:       media.CatalogReference.AlbumFolderName(),
		Id:               media.CatalogReference.MediaId(),
		Open:             media.AnalysedMedia.FoundMedia.ReadMedia,
		OriginalFilename: media.AnalysedMedia.FoundMedia.MediaPath().Filename,
		Owner:            owner,
		SignatureSha256:  media.AnalysedMedia.Sha256Hash,
	})
}
