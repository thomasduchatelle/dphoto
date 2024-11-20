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
	Owner             ownermodel.Owner
	InsertMediaPort   InsertMediaPort
	ArchivePort       ArchiveMediaPort
	UploaderObservers []uploaderObserver // UploaderObservers are called after the media is uploaded, but before the media is catalogued
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

	err := u.InsertMediaPort.IndexMedias(ctx, u.Owner, catalogRequests)
	return errors.Wrapf(err, "failed to catalog medias")
}
