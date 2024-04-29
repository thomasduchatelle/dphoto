package catalog

import (
	log "github.com/sirupsen/logrus"
)

//type MediaTransfer struct {
//	FindMediaIdsBySelectorPort FindMediaIdsBySelectorPort
//	TransferMediaToAlbumPort TransferMediaToAlbumPort
//}
//
//type FindMediaIdsBySelectorPort interface {
//	FindMediaIdsBySelector(ctx context.Context, selector ...MediaSelector) ([]MediaId, error)
//}
//
//type TransferMediaToAlbumPort interface {
//	TransferMediaToAlbum(ctx context.Context, albumId AlbumId, mediaIds []MediaId) error
//}
//
//func (t *MediaTransfer) Observe(ctx context.Context, transfers TransferredMedias) error {
//	for albumId, selectors := range transfers {
//		mediaIds, err := t.FindMediaIdsBySelectorPort.FindMediaIdsBySelector(ctx, selectors...)
//		if err != nil {
//			return err
//		}
//
//		if len(mediaIds) > 0 {
//			t.TransferMediaToAlbumPort.TransferMediaToAlbum(ctx, )
//		}
//
//	}
//}

// TODO That could be an observer ?

func transferMedias(filter *FindMediaRequest, folderName FolderName) (int, error) {
	ids, err := repositoryPort.FindMediaIds(filter)
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
	return len(ids), repositoryPort.TransferMedias(filter.Owner, ids, folderName)
}
