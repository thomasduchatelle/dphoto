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
	queries := CatalogQueries(ctx)
	writeRepo := CatalogRepository(ctx)
	referencer, err := catalog.NewAlbumAutoPopulateReferencer(
		owner,
		queries,
		writeRepo,
		writeRepo,
		ArchiveTimelineMutationObserver(),
	)

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

func NewCatalogAdapter(ctx context.Context) backup.CatalogAdapter {
	return &backupcatalog.Adapter{
		InsertMediaCase: InsertMediasCase(ctx),
	}
}
