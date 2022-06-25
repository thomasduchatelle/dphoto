package catalog

import log "github.com/sirupsen/logrus"

func transferMedias(filter *FindMediaRequest, filderName string) (int, error) {
	ids, err := repositoryPort.FindMediaIds(filter)
	if err != nil {
		return 0, err
	}

	err = archivePort.MoveMedias(filter.Owner, ids, filderName)
	if err != nil {
		return 0, err
	}

	defer func() {
		log.WithFields(log.Fields{
			"Owner":      filter,
			"FolderName": filderName,
		}).Infoln(len(ids), "medias virtually moved to new album")
	}()
	return len(ids), repositoryPort.TransferMedias(filter.Owner, ids, filderName)
}
