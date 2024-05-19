package pkgfactory

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/backupcatalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

type Factory struct{}

func (f *Factory) NewCreatorReferencer(ctx context.Context, owner ownermodel.Owner) (backup.CatalogReferencer, error) {
	queries := AlbumQueries(ctx)
	writeRepo := CatalogRepository(ctx)
	referencer, err := catalog.NewAlbumAutoPopulateReferencer(
		owner,
		queries,
		writeRepo,
		writeRepo,
		ArchiveTimelineMutationObserver(),
	)

	// TODO is the albums recounted after backup complete ?

	return &backupcatalog.CatalogReferencerAdapter{
		Owner: owner,
		InsertMediaSimulator: &catalog.MediasInsertSimulator{
			FindExistingSignaturePort: writeRepo,
		},
		StatefulAlbumReferencer: referencer,
	}, errors.Wrapf(err, "NewCreatorReferencer(%s) failed", owner)
}

func (f *Factory) NewDryRunReferencer(ctx context.Context, owner ownermodel.Owner) (backup.CatalogReferencer, error) {
	//TODO implement me
	panic("implement me")
}

func NewReferencerFactory() backup.ReferencerFactory {
	return new(Factory)
}

func NewInsertMediaAdapter(ctx context.Context) backup.InsertMediaPort {
	return &backupcatalog.InsertMediaAdapter{
		CatalogInsertMedia: InsertMediasCase(ctx),
	}
}
