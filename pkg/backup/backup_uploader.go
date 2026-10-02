package backup

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

type uploaderObserver interface {
	OnBackingUpMediaRequestUploaded(ctx context.Context, request BackingUpMediaRequest) error
}

type uploader struct {
	Owner              ownermodel.Owner
	InsertMediaPort    InsertMediaPort
	ArchivePort        ArchiveMediaPort
	CompleteCoversPort CompleteCoversPort
	UploaderObservers  []uploaderObserver // UploaderObservers are called after the media is uploaded, but before the media is catalogued
}

func (u *uploader) OnMediaCatalogued(ctx context.Context, requests []BackingUpMediaRequest) error {
	catalogRequests := make([]*CatalogMediaRequest, len(requests), len(requests))

	for i, request := range requests {
		newFilename, err := u.ArchivePort.ArchiveMedia(u.Owner.Value(), &request)
		if err != nil {
			return errors.Wrapf(err, "archiving media %s failed", request.AnalysedMedia.FoundMedia.String())
		}

		catalogRequests[i] = &CatalogMediaRequest{
			BackingUpMediaRequest: &request,
			ArchiveFilename:       newFilename,
		}

		for _, observer := range u.UploaderObservers {
			err = observer.OnBackingUpMediaRequestUploaded(ctx, request)
			if err != nil {
				return err
			}
		}
	}

	if err := u.InsertMediaPort.IndexMedias(ctx, u.Owner, catalogRequests); err != nil {
		return errors.Wrapf(err, "failed to catalog medias")
	}

	return u.completeCovers(ctx, catalogRequests)
}

// completeCovers invokes the catalog cover completion for each album that just received
// at least one IMAGE in this batch, passing the inserted medias as candidates. This keeps
// the hot path cheap: no album-wide query, no work for albums with no additions, and no
// work when the album had no new IMAGE in the batch.
func (u *uploader) completeCovers(ctx context.Context, requests []*CatalogMediaRequest) error {
	if u.CompleteCoversPort == nil {
		return nil
	}

	candidatesByAlbum := make(map[string][]CoverCandidate)
	albumOrder := make([]string, 0)
	for _, request := range requests {
		if request.BackingUpMediaRequest.AnalysedMedia.Type != MediaTypeImage {
			continue
		}
		albumFolderName := request.BackingUpMediaRequest.CatalogReference.AlbumFolderName()
		if _, seen := candidatesByAlbum[albumFolderName]; !seen {
			albumOrder = append(albumOrder, albumFolderName)
		}
		candidatesByAlbum[albumFolderName] = append(candidatesByAlbum[albumFolderName], CoverCandidate{
			MediaId:  request.BackingUpMediaRequest.CatalogReference.MediaId(),
			Filename: request.ArchiveFilename,
		})
	}

	for _, albumFolderName := range albumOrder {
		if err := u.CompleteCoversPort.CompleteCoversFromCandidates(ctx, u.Owner, albumFolderName, candidatesByAlbum[albumFolderName]); err != nil {
			return errors.Wrapf(err, "failed to complete covers for album %s", albumFolderName)
		}
	}

	return nil
}
