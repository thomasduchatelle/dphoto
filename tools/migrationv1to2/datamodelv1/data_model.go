package datamodelv1

import (
	"fmt"
	"time"
)

// Package dynamodb package store all the data in a single multi-tenant table:
// - OWNER
//   > Album X meta
//   > Album Y meta
// - MEDIA (OWNER#SIGNATURE)
//   > #META
//   > LOCATION
//   > MOVE LOCATION
//   > MOVE LOCATION
// - MOVE TRANSACTION (...#uniqueID)

const (
	IsoTime              = "2006-01-02T15:04:05"
	transactionReady     = MediaMoveTransactionStatus("MOVE_TRANSACTION_READY")
	transactionPreparing = MediaMoveTransactionStatus("MOVE_TRANSACTION_PREPARING")
)

type TablePk struct {
	PK string // Partition key ; see what's used depending on object types
	SK string // Sort key ; see what's used depending on object types
}

type AlbumIndexKey struct {
	AlbumIndexPK string // AlbumIndexPK is same than album's TablePk.PK
	AlbumIndexSK string
}

type AlbumData struct {
	TablePk
	AlbumIndexKey
	AlbumOwner      string
	AlbumName       string
	AlbumFolderName string
	AlbumStart      time.Time
	AlbumEnd        time.Time
}

type MediaData struct {
	TablePk
	AlbumIndexKey
	Type          string           // Type is either PHOTO or VIDEO
	DateTime      time.Time        // DateTime time used in AlbumIndexKey
	Details       MediaDetailsData // Details are other attributes from domain model, stored as it
	Filename      string           // Filename is the original filename for display purpose only ; physical filename is in MediaLocationData
	SignatureSize int
	SignatureHash string
}

// MediaDetailsData is a sub-object ; not stored directly
type MediaDetailsData map[string]interface{}

type MediaLocationData struct {
	TablePk
	FolderName    string // FolderName is where the media is physically located: its current album folder or previous album if the physical move haven't been flushed yet
	Filename      string // Filename is the physical name of the image
	SignatureSize int
	SignatureHash string
}

type MediaMoveTransactionStatus string

type MediaMoveTransactionData struct {
	TablePk
	MoveTransactionStatus MediaMoveTransactionStatus // MoveTransactionStatus is false until all media to be moved have a MediaMoveOrderData created and their album updated
}

type MediaMoveOrderData struct {
	TablePk                  // TablePk.PK is the same than the media, TablePk.SK is the transaction
	MoveTransaction   string // MoveTransaction is a copy of TablePk.SK used by 'MoveOrder' index (thus, only orders are in the index, not transactions)
	DestinationFolder string // DestinationFolder is the folder name of the album to which media must be moved.
}

func albumPrimaryKey(owner string, folderName string) TablePk {
	return TablePk{
		PK: owner,
		SK: fmt.Sprintf("ALBUM#%s", folderName),
	}
}

//func mediaPrimaryKey(owner string, signature *catalogmodel.MediaSignature) TablePk {
//	return TablePk{
//		PK: fmt.Sprintf("%s#MEDIA#%s", owner, mediaBusinessSignature(signature)),
//		SK: "#METADATA",
//	}
//}
//
//func mediaLocationPrimaryKey(owner string, signature *catalogmodel.MediaSignature) TablePk {
//	return TablePk{
//		PK: fmt.Sprintf("%s#MEDIA#%s", owner, mediaBusinessSignature(signature)),
//		SK: "LOCATION",
//	}
//}

func albumIndexedKey(owner, folderName string) AlbumIndexKey {
	return AlbumIndexKey{
		AlbumIndexPK: fmt.Sprintf("%s#%s", owner, folderName),
		AlbumIndexSK: fmt.Sprintf("#METADATA#ALBUM#%s", folderName),
	}
}

//func mediaAlbumIndexedKey(owner string, folderName string, dateTime time.Time, signature *catalogmodel.MediaSignature) AlbumIndexKey {
//	return AlbumIndexKey{
//		AlbumIndexPK: fmt.Sprintf("%s#%s", owner, folderName),
//		AlbumIndexSK: fmt.Sprintf("MEDIA#%s#%s", dateTime.Format(IsoTime), mediaBusinessSignature(signature)),
//	}
//}
