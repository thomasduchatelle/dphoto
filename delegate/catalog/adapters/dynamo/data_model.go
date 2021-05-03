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
	"duchatelle.io/dphoto/dphoto/catalog"
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
	AlbumIndexPK string
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
	MoveTransactionStatus MediaMoveTransactionStatus // Ready is false until all media to be moved have a MediaMoveOrderData created and their album updated
}

type MediaMoveOrderData struct {
	TablePk
	MoveTransaction   string // MoveTransaction is a copy of SK to only index MediaMoveOrderData (and not the whole table content)
	DestinationFolder string // DestinationFolder is the folder name of the album to which media must be moved.
}

type MediaDetailsData map[string]interface{}
type dynamoObject map[string]*dynamodb.AttributeValue

func (r *Rep) albumPrimaryKey(foldername string) TablePk {
	return TablePk{
		PK: r.RootOwner,
		SK: fmt.Sprintf("ALBUM#%s", foldername),
	}
}

func (r *Rep) mediaPrimaryKey(signature *catalog.MediaSignature) TablePk {
	return TablePk{
		PK: fmt.Sprintf("%s#MEDIA#%s", r.RootOwner, r.mediaBusinessSignature(signature)),
		SK: "#METADATA",
	}
}

func (r *Rep) mediaLocationPrimaryKey(signature *catalog.MediaSignature) TablePk {
	return TablePk{
		PK: fmt.Sprintf("%s#MEDIA#%s", r.RootOwner, r.mediaBusinessSignature(signature)),
		SK: "LOCATION",
	}
}

func (r *Rep) moveTransactionPrimaryKey(moveTransactionId string) TablePk {
	return TablePk{
		PK: moveTransactionId,
		SK: "#METADATA#",
	}
}

func (r *Rep) mediaMoveOrderPrimaryKey(signature *catalog.MediaSignature, moveTransactionId string) TablePk {
	return TablePk{
		PK: fmt.Sprintf("%s#MEDIA#%s", r.RootOwner, r.mediaBusinessSignature(signature)),
		SK: moveTransactionId,
	}
}

// mediaPrimaryKeyFromSubEntry takes the key of any entries related to the media (metadata, location, or move order) and return location entry key
func (r *Rep) mediaLocationKeyFromMediaKey(mediaKey map[string]*dynamodb.AttributeValue) (map[string]*dynamodb.AttributeValue, error) {
	pk, ok := mediaKey["PK"]
	if !ok {
		return nil, errors.Errorf("mediaKey must contains key 'PK': %+v", mediaKey)
	}

	return map[string]*dynamodb.AttributeValue{
		"PK": pk,
		"SK": r.mustAttribute("LOCATION"),
	}, nil
}

func (r *Rep) albumIndexedKey(owner, folderName string) AlbumIndexKey {
	return AlbumIndexKey{
		AlbumIndexPK: fmt.Sprintf("%s#%s", owner, folderName),
		AlbumIndexSK: fmt.Sprintf("#METADATA#ALBUM#%s", folderName),
	}
}

func (r *Rep) mediaAlbumIndexedKey(owner string, folderName string, dateTime time.Time, signature *catalog.MediaSignature) AlbumIndexKey {
	return AlbumIndexKey{
		AlbumIndexPK: fmt.Sprintf("%s#%s", owner, folderName),
		AlbumIndexSK: fmt.Sprintf("MEDIA#%s#%s", dateTime.Format(IsoTime), r.mediaBusinessSignature(signature)),
	}
}

// mediaBusinessSignature generate a string representing uniquely the album.Media
func (r *Rep) mediaBusinessSignature(signature *catalog.MediaSignature) string {
	return fmt.Sprintf("%s#%v", signature.SignatureSha256, signature.SignatureSize)
}

func (r *Rep) marshalAlbum(album *catalog.Album) (map[string]*dynamodb.AttributeValue, error) {
	if isBlank(album.FolderName) {
		return nil, errors.WithStack(errors.New("folderName must not be blank"))
	}

	return dynamodbattribute.MarshalMap(&AlbumData{
		TablePk:         r.albumPrimaryKey(album.FolderName),
		AlbumIndexKey:   r.albumIndexedKey(r.RootOwner, album.FolderName),
		AlbumName:       album.Name,
		AlbumFolderName: album.FolderName,
		AlbumStart:      album.Start,
		AlbumEnd:        album.End,
	})
}

func (r *Rep) unmarshalAlbum(attributes map[string]*dynamodb.AttributeValue) (*catalog.Album, error) {
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
func (r *Rep) marshalMedia(media *catalog.CreateMediaRequest) (map[string]*dynamodb.AttributeValue, map[string]*dynamodb.AttributeValue, error) {
	if isBlank(media.Signature.SignatureSha256) || media.Signature.SignatureSize == 0 || media.Details.DateTime.IsZero() {
		return nil, nil, errors.WithStack(errors.Errorf("media must have a valid signature and date [sha256=%v ; size=%v ; time=%v]", media.Signature.SignatureSha256, media.Signature.SignatureSize, media.Details.DateTime))
	}

	var details map[string]interface{}
	err := mapstructure.Decode(media.Details, &details)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to encode details values from media %+v", media.Details)
	}

	mediaEntry, err := dynamodbattribute.MarshalMap(&MediaData{
		TablePk:       r.mediaPrimaryKey(&media.Signature),
		AlbumIndexKey: r.mediaAlbumIndexedKey(r.RootOwner, media.Location.FolderName, media.Details.DateTime, &media.Signature),
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
		TablePk:       r.mediaLocationPrimaryKey(&media.Signature),
		FolderName:    media.Location.FolderName,
		Filename:      media.Location.Filename,
		SignatureSize: media.Signature.SignatureSize,
		SignatureHash: media.Signature.SignatureSha256,
	})
	return mediaEntry, locationEntry, err
}

func (r *Rep) unmarshalMediaMetaData(attributes map[string]*dynamodb.AttributeValue) (*catalog.MediaMeta, error) {
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

func (r *Rep) marshalMediaLocationFromMoveOrder(moveOrder *catalog.MovedMedia) (dynamoObject, error) {
	return dynamodbattribute.MarshalMap(&MediaLocationData{
		TablePk:       r.mediaLocationPrimaryKey(&moveOrder.Signature),
		FolderName:    moveOrder.TargetFolderName,
		Filename:      moveOrder.Filename,
		SignatureSize: moveOrder.Signature.SignatureSize,
		SignatureHash: moveOrder.Signature.SignatureSha256,
	})
}

func (r *Rep) marshalMoveTransaction(moveTransactionId string, status MediaMoveTransactionStatus) (dynamoObject, error) {
	return dynamodbattribute.MarshalMap(&MediaMoveTransactionData{
		TablePk:               r.moveTransactionPrimaryKey(moveTransactionId),
		MoveTransactionStatus: status,
	})
}

// mediaKey must have the attributes of TablePk
func (r *Rep) marshalMoveOrder(mediaKey dynamoObject, moveTransactionId, destinationFolder string) (dynamoObject, error) {
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

func (r *Rep) unmarshalMoveOrder(attributes dynamoObject) (*MediaMoveOrderData, error) {
	var order MediaMoveOrderData
	err := dynamodbattribute.UnmarshalMap(attributes, &order)
	return &order, err
}

func (r *Rep) unmarshalMediaItems(items []map[string]*dynamodb.AttributeValue) (location *MediaLocationData, orders []*MediaMoveOrderData, err error) {
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
