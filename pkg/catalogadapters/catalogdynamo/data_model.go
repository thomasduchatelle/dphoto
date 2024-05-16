// Package catalogdynamo package store all the data in a single multi-tenant table:
package catalogdynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"strings"
	"time"
)

const (
	IsoTime = "2006-01-02T15:04:05"
)

// AlbumIndexKey is a secondary key to index medias per albums
type AlbumIndexKey struct {
	AlbumIndexPK string // AlbumIndexPK is same than album's TablePk.PK
	AlbumIndexSK string // AlbumIndexSK identify the object within the index, and is naturally sorted
}

type AlbumRecord struct {
	appdynamodb.TablePk
	AlbumIndexKey
	AlbumOwner      string // AlbumOwner has been added to the data structure on 18 Apr 2022
	AlbumName       string
	AlbumFolderName string
	AlbumStart      time.Time
	AlbumEnd        time.Time
}

type MediaRecord struct {
	appdynamodb.TablePk
	AlbumIndexKey
	Id            string                 // Id is the unique identifier of the media
	Type          string                 // Type is either PHOTO or VIDEO
	DateTime      time.Time              // DateTime time used in AlbumIndexKey
	Details       map[string]interface{} // Details are other attributes from domain model, stored as it
	Filename      string                 // Filename is the original filename for display purpose only ; physical filename is in MediaLocationData
	SignatureSize int
	SignatureHash string
}

func AlbumPrimaryKey(owner ownermodel.Owner, folderName catalog.FolderName) appdynamodb.TablePk {
	return appdynamodb.TablePk{
		PK: fmt.Sprintf("%s#ALBUM", owner),
		SK: fmt.Sprintf("ALBUM#%s", folderName),
	}
}

func MediaPrimaryKey(owner ownermodel.Owner, id catalog.MediaId) appdynamodb.TablePk {
	return appdynamodb.TablePk{
		PK: appdynamodb.MediaPrimaryKeyPK(string(owner), string(id)),
		SK: "#METADATA",
	}
}

func AlbumIndexedKeyPK(owner ownermodel.Owner, folderName catalog.FolderName) string {
	return fmt.Sprintf("%s#%s", owner, folderName)
}

func AlbumIndexedKey(owner ownermodel.Owner, folderName catalog.FolderName) AlbumIndexKey {
	return AlbumIndexKey{
		AlbumIndexPK: AlbumIndexedKeyPK(owner, folderName),
		AlbumIndexSK: "#METADATA",
	}
}

func MediaAlbumIndexedKey(owner ownermodel.Owner, folderName catalog.FolderName, dateTime time.Time, id catalog.MediaId) AlbumIndexKey {
	return AlbumIndexKey{
		AlbumIndexPK: AlbumIndexedKeyPK(owner, folderName),
		AlbumIndexSK: fmt.Sprintf("MEDIA#%s#%s", dateTime.Format(IsoTime), id),
	}
}

func marshalAlbum(album *catalog.Album) (map[string]types.AttributeValue, error) {
	if err := album.FolderName.IsValid(); err != nil {
		return nil, errors.WithStack(err)
	}

	return attributevalue.MarshalMap(&AlbumRecord{
		TablePk:         AlbumPrimaryKey(album.Owner, album.FolderName),
		AlbumIndexKey:   AlbumIndexedKey(album.Owner, album.FolderName),
		AlbumOwner:      album.Owner.String(),
		AlbumName:       album.Name,
		AlbumFolderName: album.FolderName.String(),
		AlbumStart:      album.Start,
		AlbumEnd:        album.End,
	})
}

func unmarshalAlbum(attributes map[string]types.AttributeValue) (*catalog.Album, error) {
	var data AlbumRecord
	err := attributevalue.UnmarshalMap(attributes, &data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal attributes %+v", attributes)
	}

	return &catalog.Album{
		AlbumId: catalog.AlbumId{
			Owner:      ownermodel.Owner(data.AlbumOwner),
			FolderName: catalog.NewFolderName(data.AlbumFolderName),
		},
		Name:  data.AlbumName,
		Start: data.AlbumStart,
		End:   data.AlbumEnd,
	}, nil
}

// marshalMedia return both Media metadata attributes and location attributes
func marshalMedia(owner ownermodel.Owner, media *catalog.CreateMediaRequest) (map[string]types.AttributeValue, error) {
	if err := owner.IsValid(); err != nil {
		return nil, err
	}
	if isBlank(string(media.Id)) {
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

	return attributevalue.MarshalMap(&MediaRecord{
		TablePk:       MediaPrimaryKey(owner, media.Id),
		AlbumIndexKey: MediaAlbumIndexedKey(owner, media.FolderName, media.Details.DateTime, media.Id),
		Id:            string(media.Id),
		Type:          string(media.Type),
		DateTime:      media.Details.DateTime,
		Details:       details,
		Filename:      media.Filename,
		SignatureSize: media.Signature.SignatureSize,
		SignatureHash: media.Signature.SignatureSha256,
	})
}

func unmarshalMediaMetaData(attributes map[string]types.AttributeValue) (*catalog.MediaMeta, error) {
	var data MediaRecord
	err := attributevalue.UnmarshalMap(attributes, &data)
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
		Id: catalog.MediaId(data.Id),
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

func readMediaId(record map[string]types.AttributeValue) catalog.MediaId {
	if id, ok := record["Id"].(*types.AttributeValueMemberS); ok && id.Value != "" {
		return catalog.MediaId(id.Value)
	}

	return ""
}

// isBlank returns true is value is empty, or contains only spaces
func isBlank(value string) bool {
	return value == "" || strings.Trim(value, " ") == ""
}
