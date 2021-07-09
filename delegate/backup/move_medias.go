package backup

import "duchatelle.io/dphoto/dphoto/backup/interactors"

// MovePhysicalStorage is a pass through to the adapter.
func MovePhysicalStorage(owner string, folderName, filename, destinationFolderName string) (string, error) {
	return interactors.OnlineStoragePort.MoveFile(owner, folderName, filename, destinationFolderName)
}
