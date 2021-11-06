package backup

import "github.com/thomasduchatelle/dphoto/delegate/backup/interactors"

// MovePhysicalStorage is a pass through to the adapter.
func MovePhysicalStorage(owner string, folderName, filename, destinationFolderName string) (string, error) {
	return interactors.OnlineStoragePort.MoveFile(owner, folderName, filename, destinationFolderName)
}
