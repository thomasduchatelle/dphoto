package backup

import "github.com/pkg/errors"

func newBackupUploader(owner string) runnerUploader {
	return func(buffer []*BackingUpMediaRequest, progressChannel chan *ProgressEvent) error {
		catalogRequests := make([]*CatalogMediaRequest, len(buffer), len(buffer))

		for i, request := range buffer {
			newFilename, err := archivePort.ArchiveMedia(owner, request)
			if err != nil {
				return errors.Wrapf(err, "archiving media %s failed", request.AnalysedMedia.FoundMedia.String())
			}

			catalogRequests[i] = &CatalogMediaRequest{
				BackingUpMediaRequest: request,
				ArchiveFilename:       newFilename,
			}

			progressChannel <- &ProgressEvent{
				Type:      ProgressEventUploaded,
				Count:     1,
				Size:      request.AnalysedMedia.FoundMedia.Size(),
				Album:     request.FolderName,
				MediaType: request.AnalysedMedia.Type,
			}
		}

		return errors.Wrapf(catalogPort.IndexMedias(owner, catalogRequests), "failed to catalog medias")
	}
}
