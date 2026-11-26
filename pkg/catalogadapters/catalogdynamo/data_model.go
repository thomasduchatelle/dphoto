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

// MediaByTagKey is the secondary key that lists every media carrying a tag for an owner,
// across albums, in O(matches).
type MediaByTagKey struct {
	MediaByTagPK string // MediaByTagPK is {owner}#{tag}
	MediaByTagSK string // MediaByTagSK is MEDIA#{id}, naturally unique per media
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

// CoverRecord persists the whole ordered cover set of an album as a single DynamoDB
// item. See ADR / DATA_MODEL.md.
type CoverRecord struct {
	appdynamodb.TablePk
	AlbumOwner      string
	AlbumFolderName string
	Covers          []CoverRecordItem
}

type CoverRecordItem struct {
	MediaId  string
	Filename string
	Origin   string
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
	Tags          []string `dynamodbav:"Tags,stringset,omitempty"` // Tags is the source of truth of the media's tags, stored as a String Set so ADD/DELETE merge concurrent writes ; TAG# siblings are a derived index
}

// TagRecord is one item per (media, tag), stored as a sibling of the media on its partition
// (PK={owner}#MEDIA#{id}, SK=TAG#{tag}). It is a derived index item feeding the MediaByTagIndex
// GSI to list medias of a tag across albums ; the tag set source of truth is the media's
// #METADATA row. It carries no album reference, so moving a media between albums never touches
// its tags. TagVersion carries the media's #METADATA TagVersion at the time it was written, so
// concurrent add/remove of the same tag reconcile the sibling in version order.
type TagRecord struct {
	appdynamodb.TablePk
	MediaByTagKey
	Id         string // Id is the media id the tag is attached to
	Tag        string // Tag is the normalised tag string
	TagVersion int64  // TagVersion is the media's #METADATA TagVersion that last wrote this sibling ; guards against a stale concurrent delete
}

func TagPrimaryKey(owner ownermodel.Owner, mediaId catalog.MediaId, tag catalog.Tag) appdynamodb.TablePk {
	return appdynamodb.TablePk{
		PK: appdynamodb.MediaPrimaryKeyPK(string(owner), string(mediaId)),
		SK: fmt.Sprintf("TAG#%s", tag),
	}
}

func MediaByTagIndexedKey(owner ownermodel.Owner, mediaId catalog.MediaId, tag catalog.Tag) MediaByTagKey {
	return MediaByTagKey{
		MediaByTagPK: fmt.Sprintf("%s#%s", owner, tag),
		MediaByTagSK: fmt.Sprintf("MEDIA#%s", mediaId),
	}
}

func AlbumPrimaryKey(owner ownermodel.Owner, folderName catalog.FolderName) appdynamodb.TablePk {
	return appdynamodb.TablePk{
		PK: fmt.Sprintf("%s#ALBUM", owner),
		SK: fmt.Sprintf("ALBUM#%s", folderName),
	}
}

func CoverPrimaryKey(owner ownermodel.Owner, folderName catalog.FolderName) appdynamodb.TablePk {
	return appdynamodb.TablePk{
		PK: fmt.Sprintf("%s#ALBUM", owner),
		SK: fmt.Sprintf("ALBUM#%s#COVERS", folderName),
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

func marshalCovers(albumId catalog.AlbumId, covers []catalog.Cover) (map[string]types.AttributeValue, error) {
	if err := albumId.IsValid(); err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]CoverRecordItem, 0, len(covers))
	for _, cover := range covers {
		items = append(items, CoverRecordItem{
			MediaId:  string(cover.MediaId),
			Filename: cover.Filename,
			Origin:   string(cover.Origin),
		})
	}

	return attributevalue.MarshalMap(&CoverRecord{
		TablePk:         CoverPrimaryKey(albumId.Owner, albumId.FolderName),
		AlbumOwner:      albumId.Owner.String(),
		AlbumFolderName: albumId.FolderName.String(),
		Covers:          items,
	})
}

func unmarshalCovers(attributes map[string]types.AttributeValue) ([]catalog.Cover, error) {
	var data CoverRecord
	if err := attributevalue.UnmarshalMap(attributes, &data); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal cover attributes %+v", attributes)
	}

	covers := make([]catalog.Cover, 0, len(data.Covers))
	for _, item := range data.Covers {
		covers = append(covers, catalog.Cover{
			MediaId:  catalog.MediaId(item.MediaId),
			Filename: item.Filename,
			Origin:   catalog.CoverOrigin(item.Origin),
		})
	}
	return covers, nil
}

// marshalMedia return both Media metadata attributes and location attributes
func marshalMedia(owner ownermodel.Owner, media *catalog.CreateMediaRequest) (map[string]types.AttributeValue, error) {
	if err := owner.IsValid(); err != nil {
		return nil, err
	}
	if isBlank(string(media.Id)) {
		return nil, errors.Errorf("media ID is mandatory")
	}
	if isBlank(media.Filename) {
		return nil, errors.Errorf("media filename is mandatory")
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

	tags, err := unmarshalTags(data.Tags)
	if err != nil {
		return nil, err
	}

	media := catalog.MediaMeta{
		Id: catalog.MediaId(data.Id),
		Signature: catalog.MediaSignature{
			SignatureSha256: data.SignatureHash,
			SignatureSize:   data.SignatureSize,
		},
		Filename: data.Filename,
		Type:     catalog.MediaType(data.Type),
		Details:  details,
		Tags:     tags,
	}

	return &media, nil
}

func unmarshalTags(raw []string) ([]catalog.Tag, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	tags := make([]catalog.Tag, 0, len(raw))
	for _, value := range raw {
		tag, err := catalog.NewTag(value)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid persisted tag %q", value)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func marshalTag(owner ownermodel.Owner, mediaId catalog.MediaId, tag catalog.Tag, version int64) (map[string]types.AttributeValue, error) {
	if err := owner.IsValid(); err != nil {
		return nil, errors.WithStack(err)
	}
	if isBlank(string(mediaId)) {
		return nil, errors.Errorf("media ID is mandatory")
	}

	return attributevalue.MarshalMap(&TagRecord{
		TablePk:       TagPrimaryKey(owner, mediaId, tag),
		MediaByTagKey: MediaByTagIndexedKey(owner, mediaId, tag),
		Id:            string(mediaId),
		Tag:           tag.String(),
		TagVersion:    version,
	})
}

// unmarshalTagMediaId reads the media id from a MediaByTagIndex item. The index projects
// KEYS_ONLY, so the id is derived from the MediaByTagSK (MEDIA#{id}) rather than a projected
// attribute.
func unmarshalTagMediaId(attributes map[string]types.AttributeValue) (catalog.MediaId, error) {
	value, ok := attributes["MediaByTagSK"].(*types.AttributeValueMemberS)
	if !ok {
		return "", errors.Errorf("tag index item is missing MediaByTagSK: %+v", attributes)
	}
	id, found := strings.CutPrefix(value.Value, "MEDIA#")
	if !found || id == "" {
		return "", errors.Errorf("tag index item has an unexpected MediaByTagSK %q", value.Value)
	}
	return catalog.MediaId(id), nil
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
