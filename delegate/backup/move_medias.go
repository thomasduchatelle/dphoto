package backup

import "duchatelle.io/dphoto/dphoto/backup/interactors"

func MovePhysicalStorage(folderName, filename, destinationFolderName string) (string, error) {
	return interactors.OnlineStoragePort.MoveFile(folderName, filename, destinationFolderName)
}
