package catalog

import log "github.com/sirupsen/logrus"

func transferMedias(filter *FindMediaRequest, folderName string) (int, error) {
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
