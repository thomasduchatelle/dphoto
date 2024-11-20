package pkgfactory

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysers"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/backuparchive"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/backupcatalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

type MultiFilesBackup func(ctx context.Context, owner ownermodel.Owner, volumeSource backup.SourceVolume, optionsSlice ...backup.Options) (backup.Report, error)

type MultiFilesScanner func(ctx context.Context, owner string, volume backup.SourceVolume, optionSlice ...backup.Options) ([]*backup.ScannedFolder, error)

func NewMultiFilesBackup(ctx context.Context) MultiFilesBackup {
	factory.InitArchive(ctx)

	return func(ctx context.Context, owner ownermodel.Owner, volume backup.SourceVolume, optionsSlice ...backup.Options) (backup.Report, error) {
		batch := &backup.BatchBackup{
			CataloguerFactory: new(AlbumCreatorCataloguerFactory),
			DetailsReaders:    analysers.ListDetailReaders(),
			InsertMediaPort:   NewInsertMediaAdapter(ctx),
			ArchivePort:       backuparchive.New(),
		}

		return batch.Backup(ctx, owner, volume, backupDefaultOptionsForAWS(optionsSlice)...)
	}
}

type AlbumCreatorCataloguerFactory struct{}

func (f *AlbumCreatorCataloguerFactory) NewOwnerScopedCataloguer(ctx context.Context, owner ownermodel.Owner) (backup.Cataloguer, error) {
	queries := AlbumQueries(ctx)
	writeRepo := CatalogRepository(ctx)
	referencer, err := catalog.NewAlbumAutoPopulateReferencer(
		owner,
		queries,
		writeRepo,
		writeRepo,
		ArchiveTimelineMutationObserver(ctx),
		CommandHandlerAlbumSize(ctx),
	)

	return &backupcatalog.CatalogReferencerAdapter{
		Owner: owner,
		InsertMediaSimulator: &catalog.MediasInsertSimulator{
			FindExistingSignaturePort: writeRepo,
		},
		StatefulAlbumReferencer: referencer,
	}, errors.Wrapf(err, "NewOwnerScopedCataloguer(%s) failed", owner)
}

type DryRunCataloguerFactory struct{}

func (f *DryRunCataloguerFactory) NewOwnerScopedCataloguer(ctx context.Context, owner ownermodel.Owner) (backup.Cataloguer, error) {
	queries := AlbumQueries(ctx)
	writeRepo := CatalogRepository(ctx)
	referencer, err := catalog.NewAlbumDryRunReferencer(
		owner,
		queries,
	)

	return &backupcatalog.CatalogReferencerAdapter{
		Owner: owner,
		InsertMediaSimulator: &catalog.MediasInsertSimulator{
			FindExistingSignaturePort: writeRepo,
		},
		StatefulAlbumReferencer: referencer,
	}, errors.Wrapf(err, "NewDryRunCataloguer(%s) failed", owner)
}

func NewInsertMediaAdapter(ctx context.Context) backup.InsertMediaPort {
	return &backupcatalog.InsertMediaAdapter{
		CatalogInsertMedia: InsertMediasCase(ctx),
	}
}

func NewMultiFilesScanner(ctx context.Context) MultiFilesScanner {
	return func(ctx context.Context, owner string, volume backup.SourceVolume, optionSlice ...backup.Options) ([]*backup.ScannedFolder, error) {
		batchScanner := &backup.BatchScanner{
			CataloguerFactory: new(DryRunCataloguerFactory),
			DetailsReaders:    analysers.ListDetailReaders(),
		}
		return batchScanner.Scan(ctx, ownermodel.Owner(owner), volume, backupDefaultOptionsForAWS(optionSlice)...)
	}
}

func backupDefaultOptionsForAWS(optionsSlice []backup.Options) []backup.Options {
	// AWS Optimisation - comes last as low priority
	options := make([]backup.Options, len(optionsSlice)+1, len(optionsSlice)+1)
	copy(options, optionsSlice)
	options[len(options)-1] = backup.OptionsBatchSize(dynamoutils.DynamoReadBatchSize)
	return options
}
