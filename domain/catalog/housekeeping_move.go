package catalog

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// FindMoveTransactions lists transactions of media requiring to be physically moved.
func FindMoveTransactions(owner string) ([]*MoveTransaction, error) {
	return Repository.FindReadyMoveTransactions(owner)
}

// RelocateMovedMedias drives the physical re-location of all medias that have been flagged.
func RelocateMovedMedias(owner string, operator MoveMediaOperator, transactionId string) (int, error) {
	transactions, err := Repository.FindReadyMoveTransactions(owner)
	if err != nil {
		return 0, err
	}

	if len(transactions) == 0 {
		log.WithField("HouseKeepingTransaction", transactionId).Infoln("No physical move to perform, aborting.")
		return 0, err
	}

	if transactionId != transactions[0].TransactionId {
		return 0, errors.Errorf("Transactions must be proceed in creation order, %s is the first, not %s.", transactions[0].TransactionId, transactionId)
	}
	total := transactions[0].Count

	err = operator.UpdateStatus(0, total)
	if err != nil {
		return 0, err
	}

	count := 0
	pageToken := ""
	for operator.Continue() && (pageToken != "" || count == 0) {
		var moves []*MovedMedia
		moves, pageToken, err = Repository.FindFilesToMove(transactionId, pageToken)
		if len(moves) == 0 {
			break
		}

		for _, move := range moves {
			move.TargetFilename, err = operator.Move(MediaLocation{
				FolderName: move.SourceFolderName,
				Filename:   move.SourceFilename,
			}, MediaLocation{
				FolderName: move.TargetFolderName,
				Filename:   move.TargetFilename,
			})
			if err != nil {
				return count, err
			}
		}

		err = Repository.UpdateMediasLocation(owner, transactionId, moves)
		if err != nil {
			return count, err
		}

		count += len(moves)
		err = operator.UpdateStatus(count, total)
		if err != nil {
			return count, err
		}
	}

	log.WithField("MoveTransactionId", transactionId).Infoln("Move transaction completed.")
	err = Repository.DeleteEmptyMoveTransaction(transactionId)

	return count, err
}
