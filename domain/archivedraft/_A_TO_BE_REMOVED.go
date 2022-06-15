package archivedraft

type MovedAdapter interface {
	// UpdateMedias updates metadata and mark the media to be moved, the AlbumFolderName is never updated (part of the primary key)
	UpdateMedias(filter *UpdateMediaFilter, newFolderName string) (string, int, error)
	// FindMediaLocations get the list of all locations media might be in (expected only 1 unless media is pending a physical re-location)
	FindMediaLocations(owner string, signature MediaSignature) ([]*MediaLocation, error)

	FindReadyMoveTransactions(owner string) ([]*MoveTransaction, error)
	FindFilesToMove(transactionId, pageToken string) ([]*MovedMedia, string, error)
	UpdateMediasLocation(owner string, transactionId string, moves []*MovedMedia) error
	DeleteEmptyMoveTransaction(transactionId string) error
}

// MovedMedia is a record of a media that will be, or have been, physically moved
type MovedMedia struct {
	Signature        MediaSignature
	SourceFolderName string
	SourceFilename   string
	TargetFolderName string
	TargetFilename   string
}

type MediaSignatureAndLocation struct {
	Location  MediaLocation
	Signature MediaSignature
}

type MoveTransaction struct {
	TransactionId string
	Count         int // Count is the number of medias to be moved as part of this transaction
}

type MoveMediaOperator interface {
	// Move must perform the physical move of the file to a different directory ; return the final name if it has been changed
	Move(source, dest MediaLocation) (string, error)

	// UpdateStatus informs of the global status of the move operation
	UpdateStatus(done, total int) error
	// Continue requests if the operation should continue or be interrupted
	Continue() bool
}
