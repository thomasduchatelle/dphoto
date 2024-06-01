package backupcatalog

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

type CatalogInsertMedia interface {
	Insert(ctx context.Context, owner ownermodel.Owner, medias []catalog.CreateMediaRequest) error
}

type InsertMediaAdapter struct {
	CatalogInsertMedia CatalogInsertMedia
}

func (a *InsertMediaAdapter) IndexMedias(ctx context.Context, owner ownermodel.Owner, requests []*backup.CatalogMediaRequest) error {
	creates := make([]catalog.CreateMediaRequest, len(requests))
	for i, request := range requests {
		reference, valid := request.BackingUpMediaRequest.CatalogReference.(Reference)
		if !valid {
			return errors.Errorf("%T is not a reference supported by InsertMediaAdapter adapter", request.BackingUpMediaRequest.CatalogReference)
		}

		creates[i] = catalog.CreateMediaRequest{
			Id:         reference.MediaReference.ProvisionalMediaId,
			Signature:  reference.MediaReference.Signature,
			FolderName: reference.AlbumReference.AlbumId.FolderName,
			Filename:   request.ArchiveFilename,
			Type:       catalog.MediaType(request.BackingUpMediaRequest.AnalysedMedia.Type),
			Details: catalog.MediaDetails{
				Width:         request.BackingUpMediaRequest.AnalysedMedia.Details.Width,
				Height:        request.BackingUpMediaRequest.AnalysedMedia.Details.Height,
				DateTime:      request.BackingUpMediaRequest.AnalysedMedia.Details.DateTime,
				Orientation:   catalog.MediaOrientation(request.BackingUpMediaRequest.AnalysedMedia.Details.Orientation),
				Make:          request.BackingUpMediaRequest.AnalysedMedia.Details.Make,
				Model:         request.BackingUpMediaRequest.AnalysedMedia.Details.Model,
				GPSLatitude:   request.BackingUpMediaRequest.AnalysedMedia.Details.GPSLatitude,
				GPSLongitude:  request.BackingUpMediaRequest.AnalysedMedia.Details.GPSLongitude,
				Duration:      request.BackingUpMediaRequest.AnalysedMedia.Details.Duration,
				VideoEncoding: request.BackingUpMediaRequest.AnalysedMedia.Details.VideoEncoding,
			},
		}
	}

	return a.CatalogInsertMedia.Insert(ctx, owner, creates)
}
