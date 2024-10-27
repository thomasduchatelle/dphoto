package pkgfactory

import (
	"context"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/backup"
	_ "github.com/thomasduchatelle/dphoto/pkg/backupadapters/analysers"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/backuparchive"
	"github.com/thomasduchatelle/dphoto/pkg/backupadapters/backupcatalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

type MultiFilesBackup func(ctx context.Context, owner ownermodel.Owner, volumeSource backup.SourceVolume, optionsSlice ...backup.Options) (backup.CompletionReport, error)

type MultiFilesScanner func(ctx context.Context, owner string, volume backup.SourceVolume, optionSlice ...backup.Options) ([]*backup.ScannedFolder, error)

func NewMultiFilesBackup(ctx context.Context) MultiFilesBackup {
	factory.InitArchive(ctx)
	backup.Init(backuparchive.New(), NewReferencerFactory(), NewInsertMediaAdapter(ctx))

	return func(ctx context.Context, owner ownermodel.Owner, volume backup.SourceVolume, optionsSlice ...backup.Options) (backup.CompletionReport, error) {
		options := backupDefaultOptionsForAWS(optionsSlice)

		return backup.Backup(owner, volume, options...)
	}
}

type BackupReferencerFactory struct{}

func (f *BackupReferencerFactory) NewCreatorReferencer(ctx context.Context, owner ownermodel.Owner) (backup.CatalogReferencer, error) {
	queries := AlbumQueries(ctx)
	writeRepo := CatalogRepository(ctx)
	referencer, err := catalog.NewAlbumAutoPopulateReferencer(
		owner,
		queries,
		writeRepo,
		writeRepo,
		ArchiveTimelineMutationObserver(nil),
		CommandHandlerAlbumSize(ctx),
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

func (f *BackupReferencerFactory) NewDryRunReferencer(ctx context.Context, owner ownermodel.Owner) (backup.CatalogReferencer, error) {
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
	}, errors.Wrapf(err, "NewDryRunReferencer(%s) failed", owner)
}

func NewReferencerFactory() backup.ReferencerFactory {
	return new(BackupReferencerFactory)
}

func NewInsertMediaAdapter(ctx context.Context) backup.InsertMediaPort {
	return &backupcatalog.InsertMediaAdapter{
		CatalogInsertMedia: InsertMediasCase(ctx),
	}
}

func NewMultiFilesScanner(ctx context.Context) MultiFilesScanner {
	backup.Init(backuparchive.New(), NewReferencerFactory(), NewInsertMediaAdapter(ctx))

	return func(ctx context.Context, owner string, volume backup.SourceVolume, optionSlice ...backup.Options) ([]*backup.ScannedFolder, error) {
		batchScanner := new(backup.BatchScanner)
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
