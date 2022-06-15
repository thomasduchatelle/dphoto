package catalog

import log "github.com/sirupsen/logrus"

func transferMedias(filter *FindMediaRequest, filderName string) (int, error) {
	ids, err := dbPort.FindMediaIds(filter)
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
	return len(ids), dbPort.TransferMedias(filter.Owner, ids, filderName)
}
