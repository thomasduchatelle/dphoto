// Package catalogdynamo package store all the data in a single multi-tenant table:
package catalogdynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"strings"
	"time"
)

// DynamoDB table structure
// - OWNER (OWNER)
//   > Album X meta
//   > Album Y meta
// - MEDIA (OWNER#ID)
//   > #META = SIGNATURE + ALBUM + DETAILS
// 	 > #LOCATION = key (see archive extension)

const (
	IsoTime = "2006-01-02T15:04:05"
)

// TablePk are the primary and sort keys of the table
type TablePk struct {
	PK string // PK is the Partition key ; see what's used depending on object types
	SK string // SK is the Sort key ; see what's used depending on object types
}

// AlbumIndexKey is a secondary key to index medias per albums
type AlbumIndexKey struct {
	AlbumIndexPK string // AlbumIndexPK is same than album's TablePk.PK
	AlbumIndexSK string // AlbumIndexSK identify the object within the index, and is naturally sorted
}

type AlbumRecord struct {
	TablePk
	AlbumIndexKey
	AlbumOwner      string // AlbumOwner has been added to the data structure on 18 Apr 2022
	AlbumName       string
	AlbumFolderName string
	AlbumStart      time.Time
	AlbumEnd        time.Time
}

type MediaRecord struct {
	TablePk
	AlbumIndexKey
	Id            string                 // Id is the unique identifier of the media
	Type          string                 // Type is either PHOTO or VIDEO
	DateTime      time.Time              // DateTime time used in AlbumIndexKey
	Details       map[string]interface{} // Details are other attributes from domain model, stored as it
	Filename      string                 // Filename is the original filename for display purpose only ; physical filename is in MediaLocationData
	SignatureSize int
	SignatureHash string
}

func AlbumPrimaryKey(owner string, folderName string) TablePk {
	return TablePk{
		PK: owner,
		SK: fmt.Sprintf("ALBUM#%s", folderName),
	}
}

func MediaPrimaryKeyPK(owner string, id string) string {
	return fmt.Sprintf("%s#MEDIA#%s", owner, id)
}

func MediaPrimaryKey(owner string, id string) TablePk {
	return TablePk{
		PK: MediaPrimaryKeyPK(owner, id),
		SK: "#METADATA",
	}
}

func AlbumIndexedKeyPK(owner string, folderName string) string {
	return fmt.Sprintf("%s#%s", owner, folderName)
}

func AlbumIndexedKey(owner, folderName string) AlbumIndexKey {
	return AlbumIndexKey{
		AlbumIndexPK: AlbumIndexedKeyPK(owner, folderName),
		AlbumIndexSK: fmt.Sprintf("#METADATA#ALBUM#%s", folderName),
	}
}

func mediaAlbumIndexedKey(owner string, folderName string, dateTime time.Time, id string) AlbumIndexKey {
	return AlbumIndexKey{
		AlbumIndexPK: AlbumIndexedKeyPK(owner, folderName),
		AlbumIndexSK: fmt.Sprintf("MEDIA#%s#%s", dateTime.Format(IsoTime), id),
	}
}

func marshalAlbum(album *catalog.Album) (map[string]*dynamodb.AttributeValue, error) {
	if isBlank(album.FolderName) {
		return nil, errors.WithStack(errors.New("folderName must not be blank"))
	}

	return dynamodbattribute.MarshalMap(&AlbumRecord{
		TablePk:         AlbumPrimaryKey(album.Owner, album.FolderName),
		AlbumIndexKey:   AlbumIndexedKey(album.Owner, album.FolderName),
		AlbumOwner:      album.Owner,
		AlbumName:       album.Name,
		AlbumFolderName: album.FolderName,
		AlbumStart:      album.Start,
		AlbumEnd:        album.End,
	})
}

func unmarshalAlbum(attributes map[string]*dynamodb.AttributeValue) (*catalog.Album, error) {
	var data AlbumRecord
	err := dynamodbattribute.UnmarshalMap(attributes, &data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal attributes %+v", attributes)
	}

	return &catalog.Album{
		Owner:      data.AlbumOwner,
		Name:       data.AlbumName,
		FolderName: data.AlbumFolderName,
		Start:      data.AlbumStart,
		End:        data.AlbumEnd,
	}, nil
}

// marshalMedia return both Media metadata attributes and location attributes
func marshalMedia(owner string, media *catalog.CreateMediaRequest) (map[string]*dynamodb.AttributeValue, error) {
	if isBlank(owner) {
		return nil, errors.Errorf("owner is mandatory")
	}
	if isBlank(media.Id) {
		return nil, errors.Errorf("media ID is mndatory")
	}
	if isBlank(media.Filename) {
		return nil, errors.Errorf("media filename is mndatory")
	}
	if isBlank(media.Signature.SignatureSha256) || media.Signature.SignatureSize == 0 || media.Details.DateTime.IsZero() {
		return nil, errors.WithStack(errors.Errorf("media must have a valid signature and date [sha256=%v ; size=%v ; time=%v]", media.Signature.SignatureSha256, media.Signature.SignatureSize, media.Details.DateTime))
	}

	var details map[string]interface{}
	err := mapstructure.Decode(media.Details, &details)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to encode details values from media %+v", media.Details)
	}

	return dynamodbattribute.MarshalMap(&MediaRecord{
		TablePk:       MediaPrimaryKey(owner, media.Id),
		AlbumIndexKey: mediaAlbumIndexedKey(owner, media.FolderName, media.Details.DateTime, media.Id),
		Id:            media.Id,
		Type:          string(media.Type),
		DateTime:      media.Details.DateTime,
		Details:       details,
		Filename:      media.Filename,
		SignatureSize: media.Signature.SignatureSize,
		SignatureHash: media.Signature.SignatureSha256,
	})
}

func unmarshalMediaMetaData(attributes map[string]*dynamodb.AttributeValue) (*catalog.MediaMeta, error) {
	var data MediaRecord
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
		Id: data.Id,
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

// isBlank returns true is value is empty, or contains only spaces
func isBlank(value string) bool {
	return value == "" || strings.Trim(value, " ") == ""
}
