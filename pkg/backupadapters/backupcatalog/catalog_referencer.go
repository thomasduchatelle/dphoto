package backupcatalog

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"time"
)

type InsertMediaSimulator interface {
	SimulateInsertingMedia(ctx context.Context, owner ownermodel.Owner, signatures []catalog.MediaSignature) ([]catalog.MediaFutureReference, error)
}

type StatefulAlbumReferencer interface {
	FindReference(ctx context.Context, mediaTime time.Time) (catalog.AlbumReference, error)
}

type CatalogReferencerAdapter struct {
	Owner                   ownermodel.Owner
	InsertMediaSimulator    InsertMediaSimulator
	StatefulAlbumReferencer StatefulAlbumReferencer
}

func (c *CatalogReferencerAdapter) Reference(ctx context.Context, analysedMedias []*backup.AnalysedMedia) (map[*backup.AnalysedMedia]backup.CatalogReference, error) {
	references := make(map[*backup.AnalysedMedia]backup.CatalogReference)

	var signatures []catalog.MediaSignature
	analysedMediasBySignature := make(map[catalog.MediaSignature][]*backup.AnalysedMedia)
	for _, media := range analysedMedias {
		signature := catalog.MediaSignature{
			SignatureSha256: media.Sha256Hash,
			SignatureSize:   media.FoundMedia.Size(),
		}

		list, duplicate := analysedMediasBySignature[signature]
		analysedMediasBySignature[signature] = append(list, media)

		if !duplicate {
			signatures = append(signatures, signature)
		}
	}

	mediaReferences, err := c.InsertMediaSimulator.SimulateInsertingMedia(ctx, c.Owner, signatures)
	if err != nil {
		return nil, err
	}

	for _, mediaReference := range mediaReferences {
		medias := analysedMediasBySignature[mediaReference.Signature]
		albumReference, err := c.StatefulAlbumReferencer.FindReference(ctx, medias[0].Details.DateTime)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find album reference for media at time %s", medias[0].Details.DateTime)
		}

		for _, media := range medias {
			references[media] = Reference{
				MediaReference: mediaReference,
				AlbumReference: albumReference,
			}
		}
	}

	return references, nil
}
