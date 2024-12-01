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
	Owner            ownermodel.Owner
	InsertMediaPort  InsertMediaPort
	ArchivePort      BArchiveAdapter
	UploaderObserver uploaderObserver
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

		err = u.UploaderObserver.OnBackingUpMediaRequestUploaded(ctx, request)
		if err != nil {
			return err
		}
	}

	err := u.InsertMediaPort.IndexMedias(ctx, u.Owner, catalogRequests)
	return errors.Wrapf(err, "failed to catalog medias")
}

func (u *uploader) Upload(buffer []*BackingUpMediaRequest, progressChannel chan *progressEvent) error {
	catalogRequests := make([]*CatalogMediaRequest, len(buffer), len(buffer))

	for i, request := range buffer {
		newFilename, err := u.ArchivePort.ArchiveMedia(u.Owner.Value(), request)
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
