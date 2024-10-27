package backup

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

type Uploader struct {
	Owner           ownermodel.Owner
	InsertMediaPort InsertMediaPort
}

func (u *Uploader) Upload(buffer []*BackingUpMediaRequest, progressChannel chan *progressEvent) error {
	catalogRequests := make([]*CatalogMediaRequest, len(buffer), len(buffer))

	for i, request := range buffer {
		newFilename, err := archivePort.ArchiveMedia(u.Owner.Value(), request)
		if err != nil {
			return errors.Wrapf(err, "archiving media %s failed", request.AnalysedMedia.FoundMedia.String())
		}

		catalogRequests[i] = &CatalogMediaRequest{
			BackingUpMediaRequest: request,
			ArchiveFilename:       newFilename,
		}

		progressChannel <- &progressEvent{
			Type:      trackUploaded,
			Count:     1,
			Size:      request.AnalysedMedia.FoundMedia.Size(),
			Album:     request.CatalogReference.AlbumFolderName(),
			MediaType: request.AnalysedMedia.Type,
		}
	}

	return errors.Wrapf(u.InsertMediaPort.IndexMedias(context.TODO(), u.Owner, catalogRequests), "failed to catalog medias")
}
