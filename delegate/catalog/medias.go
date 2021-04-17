package catalog

import log "github.com/sirupsen/logrus"

// ListMedias return a page of medias within an album
func ListMedias(folderName string, request PageRequest) (*MediaPage, error) {
	return Repository.FindMedias(folderName, request)
}

// InsertMedias stores metadata and location of photo and videos
func InsertMedias(medias []CreateMediaRequest) error {
	return Repository.InsertMedias(medias)
}

// FindSignatures returns a list of the medias already known ; they can't be duplicated
func FindSignatures(signatures []*MediaSignature) ([]*MediaSignature, error) {
	return Repository.FindExistingSignatures(signatures)
}

type MoveMediaOperator interface {
	// Move must perform the physical move of the file to a different directory
	Move(source, dest MediaLocation) error

	// UpdateStatus informs of the global status of the move operation
	UpdateStatus(done, total int) error
	// Continue requests if the operation should continue or be interrupted
	Continue() bool
}

// RelocateMovedMedias drives the physical re-location of all medias that have been flagged.
func RelocateMovedMedias(operator MoveMediaOperator) (int, error) {
	transactions, err := Repository.FindReadyMoveTransactions()
	if err != nil {
		return 0, err
	}

	if len(transactions) == 0 {
		log.Infoln("No physical move to perform, aborting.")
		return 0, err
	}

	transactionId := transactions[0].TransactionId
	total := transactions[0].Count

	err = operator.UpdateStatus(0, total)
	if err != nil {
		return 0, err
	}

	count := 0
	pageToken := ""
	for operator.Continue() && (pageToken != "" || count == 0){
		var moves []*MovedMedia
		moves, pageToken, err = Repository.FindFilesToMove(transactionId, pageToken)
		if len(moves) == 0 {
			break
		}

		for _, move := range moves {
			err = operator.Move(MediaLocation{
				FolderName: move.SourceFolderName,
				Filename:   move.Filename,
			}, MediaLocation{
				FolderName: move.TargetFolderName,
				Filename:   move.Filename,
			})
			if err != nil {
				return count, err
			}
		}

		err = Repository.UpdateMediasLocation(transactionId, moves)
		if err != nil {
			return count, err
		}

		count += len(moves)
		err = operator.UpdateStatus(count, total)
		if err != nil {
			return count, err
		}
	}

	return count, nil
}
