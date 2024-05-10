package catalog

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type TransferredMedias map[AlbumId][]MediaId

func (t TransferredMedias) IsEmpty() bool {
	count := 0
	for _, ids := range t {
		count += len(ids)
	}

	return count == 0
}

type TransferMediasPort interface {
	TransferMediasFromRecords(ctx context.Context, records MediaTransferRecords) (TransferredMedias, error)
}

type TransferMediasFunc func(ctx context.Context, records MediaTransferRecords) (TransferredMedias, error)

func (f TransferMediasFunc) TransferMediasFromRecords(ctx context.Context, records MediaTransferRecords) (TransferredMedias, error) {
	return f(ctx, records)
}

type MediaTransferExecutor struct {
	TransferMedias            TransferMediasPort
	TimelineMutationObservers []TimelineMutationObserver
}

func (d *MediaTransferExecutor) Transfer(ctx context.Context, records MediaTransferRecords) error {
	transfers, err := d.TransferMedias.TransferMediasFromRecords(ctx, records)
	if err != nil || transfers.IsEmpty() {
		return err
	}

	for _, observer := range d.TimelineMutationObservers {
		err = observer.Observe(ctx, transfers)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO That could be an observer ?

func transferMedias(filter *FindMediaRequest, folderName FolderName) (int, error) {
	ids, err := repositoryPort.FindMediaIds(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	if len(ids) == 0 {
		log.WithFields(log.Fields{
			"Owner":      filter,
			"FolderName": folderName,
		}).Infoln(len(ids), "no media to transfer to the new album")
		return 0, nil
	}

	err = archivePort.MoveMedias(filter.Owner, ids, folderName)
	if err != nil {
		return 0, err
	}

	defer func() {
		log.WithFields(log.Fields{
			"Owner":      filter,
			"FolderName": folderName,
		}).Infoln(len(ids), "medias virtually moved to new album")
	}()
	return len(ids), repositoryPort.TransferMedias(context.TODO(), filter.Owner, ids, folderName)
}
