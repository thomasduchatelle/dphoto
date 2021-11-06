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
package dynamo

import (
	"github.com/thomasduchatelle/dphoto/delegate/catalog"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

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

func mediaPrimaryKey(owner string, signature *catalog.MediaSignature) TablePk {
	return TablePk{
		PK: fmt.Sprintf("%s#MEDIA#%s", owner, mediaBusinessSignature(signature)),
		SK: "#METADATA",
	}
}

func mediaLocationPrimaryKey(owner string, signature *catalog.MediaSignature) TablePk {
	return TablePk{
		PK: fmt.Sprintf("%s#MEDIA#%s", owner, mediaBusinessSignature(signature)),
		SK: "LOCATION",
	}
}

func moveTransactionPrimaryKey(moveTransactionId string) TablePk {
	return TablePk{
		PK: moveTransactionId,
		SK: "#METADATA#",
	}
}

func mediaMoveOrderPrimaryKey(rootOwner string, signature *catalog.MediaSignature, moveTransactionId string) TablePk {
	return TablePk{
		PK: fmt.Sprintf("%s#MEDIA#%s", rootOwner, mediaBusinessSignature(signature)),
		SK: moveTransactionId,
	}
}

// mediaPrimaryKeyFromSubEntry takes the key of any entries related to the media (metadata, location, or move order) and return location entry key
func mediaLocationKeyFromMediaKey(mediaKey map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error) {
	pk, ok := mediaKey["PK"]
	if !ok {
		return nil, errors.Errorf("mediaKey must contains key 'PK': %+v", mediaKey)
	}

	return map[string]*dynamodb.AttributeValue{
		"PK": pk,
		"SK": mustAttribute("LOCATION"),
	}, nil
}

func albumIndexedKey(owner, folderName string) AlbumIndexKey {
	return AlbumIndexKey{
		AlbumIndexPK: fmt.Sprintf("%s#%s", owner, folderName),
		AlbumIndexSK: fmt.Sprintf("#METADATA#ALBUM#%s", folderName),
	}
}

func mediaAlbumIndexedKey(owner string, folderName string, dateTime time.Time, signature *catalog.MediaSignature) AlbumIndexKey {
	return AlbumIndexKey{
		AlbumIndexPK: fmt.Sprintf("%s#%s", owner, folderName),
		AlbumIndexSK: fmt.Sprintf("MEDIA#%s#%s", dateTime.Format(IsoTime), mediaBusinessSignature(signature)),
	}
}

// mediaBusinessSignature generate a string representing uniquely the album.Media
func mediaBusinessSignature(signature *catalog.MediaSignature) string {
	return fmt.Sprintf("%s#%v", signature.SignatureSha256, signature.SignatureSize)
}

func marshalAlbum(owner string, album *catalog.Album) (map[string]*dynamodb.AttributeValue, error) {
	if isBlank(album.FolderName) {
		return nil, errors.WithStack(errors.New("folderName must not be blank"))
	}

	return dynamodbattribute.MarshalMap(&AlbumData{
		TablePk:         albumPrimaryKey(owner, album.FolderName),
		AlbumIndexKey:   albumIndexedKey(owner, album.FolderName),
		AlbumName:       album.Name,
		AlbumFolderName: album.FolderName,
		AlbumStart:      album.Start,
		AlbumEnd:        album.End,
	})
}

func unmarshalAlbum(attributes map[string]*dynamodb.AttributeValue) (*catalog.Album, error) {
	var data AlbumData
	err := dynamodbattribute.UnmarshalMap(attributes, &data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal attributes %+v", attributes)
	}

	return &catalog.Album{
		Name:       data.AlbumName,
		FolderName: data.AlbumFolderName,
		Start:      data.AlbumStart,
		End:        data.AlbumEnd,
	}, nil
}

// marshalMedia return both Media metadata attributes and location attributes
func marshalMedia(owner string, media *catalog.CreateMediaRequest) (map[string]*dynamodb.AttributeValue, map[string]*dynamodb.AttributeValue, error) {
	if isBlank(media.Signature.SignatureSha256) || media.Signature.SignatureSize == 0 || media.Details.DateTime.IsZero() {
		return nil, nil, errors.WithStack(errors.Errorf("media must have a valid signature and date [sha256=%v ; size=%v ; time=%v]", media.Signature.SignatureSha256, media.Signature.SignatureSize, media.Details.DateTime))
	}

	var details map[string]interface{}
	err := mapstructure.Decode(media.Details, &details)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to encode details values from media %+v", media.Details)
	}

	mediaEntry, err := dynamodbattribute.MarshalMap(&MediaData{
		TablePk:       mediaPrimaryKey(owner, &media.Signature),
		AlbumIndexKey: mediaAlbumIndexedKey(owner, media.Location.FolderName, media.Details.DateTime, &media.Signature),
		Type:          string(media.Type),
		DateTime:      media.Details.DateTime,
		Details:       details,
		Filename:      media.Location.Filename,
		SignatureSize: media.Signature.SignatureSize,
		SignatureHash: media.Signature.SignatureSha256,
	})
	if err != nil {
		return nil, nil, err
	}

	locationEntry, err := dynamodbattribute.MarshalMap(&MediaLocationData{
		TablePk:       mediaLocationPrimaryKey(owner, &media.Signature),
		FolderName:    media.Location.FolderName,
		Filename:      media.Location.Filename,
		SignatureSize: media.Signature.SignatureSize,
		SignatureHash: media.Signature.SignatureSha256,
	})
	return mediaEntry, locationEntry, err
}

func unmarshalMediaMetaData(attributes map[string]*dynamodb.AttributeValue) (*catalog.MediaMeta, error) {
	var data MediaData
	err := dynamodbattribute.UnmarshalMap(attributes, &data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal attributes %+v", attributes)
	}

	var details catalog.MediaDetails
	err = mapstructure.Decode(data.Details, &details)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal details %+v", data.Details)
	}

	details.DateTime = data.DateTime // note: mapstructure do not support times
	media := catalog.MediaMeta{
		Signature: catalog.MediaSignature{
			SignatureSha256: data.SignatureHash,
			SignatureSize:   data.SignatureSize,
		},
		Filename: data.Filename,
		Type:     catalog.MediaType(data.Type),
		Details:  details,
	}

	return &media, nil
}

func marshalMediaLocationFromMoveOrder(owner string, moveOrder *catalog.MovedMedia) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(&MediaLocationData{
		TablePk:       mediaLocationPrimaryKey(owner, &moveOrder.Signature),
		FolderName:    moveOrder.TargetFolderName,
		Filename:      moveOrder.TargetFilename,
		SignatureSize: moveOrder.Signature.SignatureSize,
		SignatureHash: moveOrder.Signature.SignatureSha256,
	})
}

func marshalMoveTransaction(moveTransactionId string, status MediaMoveTransactionStatus) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(&MediaMoveTransactionData{
		TablePk:               moveTransactionPrimaryKey(moveTransactionId),
		MoveTransactionStatus: status,
	})
}

// marshalMoveOrder creates a marshalled MediaMoveOrderData from its transaction (mediaKey must have the attributes of TablePk)
func marshalMoveOrder(mediaKey map[string]*dynamodb.AttributeValue, moveTransactionId, destinationFolder string) (map[string]*dynamodb.AttributeValue, error) {
	if _, pkOk := mediaKey["PK"]; !pkOk {
		return nil, errors.Errorf("mediaPk do not contains mandatory PK key")
	}

	return dynamodbattribute.MarshalMap(&MediaMoveOrderData{
		TablePk: TablePk{
			PK: *mediaKey["PK"].S,
			SK: moveTransactionId,
		},
		MoveTransaction:   moveTransactionId,
		DestinationFolder: destinationFolder,
	})
}

func unmarshalMoveOrder(attributes map[string]*dynamodb.AttributeValue) (*MediaMoveOrderData, error) {
	var order MediaMoveOrderData
	err := dynamodbattribute.UnmarshalMap(attributes, &order)
	return &order, err
}

func unmarshalMediaItems(items []map[string]*dynamodb.AttributeValue) (location *MediaLocationData, orders []*MediaMoveOrderData, err error) {
	for _, item := range items {
		if sk, ok := item["SK"]; ok && sk.S != nil {
			switch {
			case strings.HasPrefix(*sk.S, "LOCATION"):
				location = &MediaLocationData{}
				err = dynamodbattribute.UnmarshalMap(item, &location)
				if err != nil {
					return
				}

			case strings.HasPrefix(*sk.S, "MOVE_ORDER"):
				var data MediaMoveOrderData
				err = dynamodbattribute.UnmarshalMap(item, &data)
				if err != nil {
					return
				}

				orders = append(orders, &data)

			default:
				log.WithFields(log.Fields{
					"SK": *sk.S,
				}).Warnln("Unknown item type")
			}
		}
	}

	return
}

// isBlank returns true is value is empty, or contains only spaces
func isBlank(value string) bool {
	return value == "" || strings.Trim(value, " ") == ""
}
