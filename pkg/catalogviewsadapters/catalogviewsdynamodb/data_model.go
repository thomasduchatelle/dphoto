package catalogviewsdynamodb

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

const (
	AvailabilityTypeOwner   = "OWNED"
	AvailabilityTypeVisitor = "VISITOR"

	RecordTypeSummary = "SUMMARY"
	RecordTypeCovers  = "COVERS"

	coversSKSuffix = "#COVERS"
)

type AlbumSummaryRecord struct {
	appdynamodb.TablePk
	AlbumOwner       string
	AlbumFolderName  string
	AvailabilityType string
	UserId           string
	RecordType       string `dynamodbav:",omitempty"`
	Count            int
	AlbumName        string `dynamodbav:",omitempty"`
	AlbumStart       string `dynamodbav:",omitempty"`
	AlbumEnd         string `dynamodbav:",omitempty"`
	AlbumViewIndexPK string
}

type AlbumCoversRecord struct {
	appdynamodb.TablePk
	AlbumOwner       string
	AlbumFolderName  string
	AvailabilityType string
	UserId           string
	RecordType       string
	Covers           []CoverRecord `dynamodbav:",omitempty"`
	AlbumViewIndexPK string
}

type CoverRecord struct {
	MediaId  string
	Filename string
	Origin   string
}

func marshalCovers(covers []catalog.Cover) []CoverRecord {
	if len(covers) == 0 {
		return nil
	}
	records := make([]CoverRecord, len(covers))
	for i, cover := range covers {
		records[i] = CoverRecord{
			MediaId:  cover.MediaId.Value(),
			Filename: cover.Filename,
			Origin:   string(cover.Origin),
		}
	}
	return records
}

func unmarshalCovers(records []CoverRecord) []catalog.Cover {
	if len(records) == 0 {
		return nil
	}
	covers := make([]catalog.Cover, len(records))
	for i, record := range records {
		covers[i] = catalog.Cover{
			MediaId:  catalog.MediaId(record.MediaId),
			Filename: record.Filename,
			Origin:   catalog.CoverOrigin(record.Origin),
		}
	}
	return covers
}

func albumsViewPK(user usermodel.UserId) string {
	return fmt.Sprintf("USER#%s#ALBUMS_VIEW", user.Value())
}

func albumViewByAlbumIndexPK(albumId catalog.AlbumId) string {
	return fmt.Sprintf("ALBUM#%s#%s#ALBUMS_VIEW", albumId.Owner.Value(), albumId.FolderName.String())
}

func albumSummaryKey(user catalogviews.Availability, albumId catalog.AlbumId) appdynamodb.TablePk {
	belongType := marshalAvailabilityType(user)

	recordKey := appdynamodb.TablePk{
		PK: albumsViewPK(user.UserId),
		SK: fmt.Sprintf("%s#%s#%s", belongType, albumId.Owner, albumId.FolderName.String()),
	}
	return recordKey
}

func albumCoversKey(user catalogviews.Availability, albumId catalog.AlbumId) appdynamodb.TablePk {
	key := albumSummaryKey(user, albumId)
	key.SK += coversSKSuffix
	return key
}

func marshalAvailabilityType(user catalogviews.Availability) string {
	belongType := AvailabilityTypeOwner
	if !user.AsOwner {
		belongType = AvailabilityTypeVisitor
	}
	return belongType
}

func marshalAlbumSummary(summary catalogviews.AlbumSummaryForUsers) ([]map[string]types.AttributeValue, error) {
	var items []map[string]types.AttributeValue
	for _, user := range summary.Users {
		recordKey := albumSummaryKey(user, summary.AlbumId)

		item, err := attributevalue.MarshalMap(AlbumSummaryRecord{
			TablePk:          recordKey,
			AlbumOwner:       summary.AlbumId.Owner.Value(),
			AlbumFolderName:  summary.AlbumId.FolderName.String(),
			AvailabilityType: marshalAvailabilityType(user),
			UserId:           user.UserId.Value(),
			RecordType:       RecordTypeSummary,
			Count:            summary.MediaCount,
			AlbumName:        summary.Name,
			AlbumStart:       marshalTime(summary.Start),
			AlbumEnd:         marshalTime(summary.End),
			AlbumViewIndexPK: albumViewByAlbumIndexPK(summary.AlbumId),
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal album summary record: %+v", summary)
		}

		items = append(items, item)
	}

	return items, nil
}

func marshalAlbumCovers(user catalogviews.Availability, albumId catalog.AlbumId, covers []catalog.Cover) (map[string]types.AttributeValue, error) {
	recordKey := albumCoversKey(user, albumId)

	item, err := attributevalue.MarshalMap(AlbumCoversRecord{
		TablePk:          recordKey,
		AlbumOwner:       albumId.Owner.Value(),
		AlbumFolderName:  albumId.FolderName.String(),
		AvailabilityType: marshalAvailabilityType(user),
		UserId:           user.UserId.Value(),
		RecordType:       RecordTypeCovers,
		Covers:           marshalCovers(covers),
		AlbumViewIndexPK: albumViewByAlbumIndexPK(albumId),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal album covers record for %s %s", user, albumId)
	}
	return item, nil
}

func unmarshalAlbumSummary(item map[string]types.AttributeValue) (*catalogviews.UserAlbumSummary, error) {
	record := &AlbumSummaryRecord{}
	if err := attributevalue.UnmarshalMap(item, record); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal album summary record: %+v", item)
	}

	availability, err := unmarshalAvailability(record.AvailabilityType, record.UserId, item)
	if err != nil {
		return nil, err
	}

	start, err := unmarshalTime(record.AlbumStart)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal AlbumStart on album summary record: %+v", item)
	}
	end, err := unmarshalTime(record.AlbumEnd)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal AlbumEnd on album summary record: %+v", item)
	}

	return &catalogviews.UserAlbumSummary{
		AlbumSummary: catalogviews.AlbumSummary{
			AlbumId: catalog.AlbumId{
				Owner:      ownermodel.Owner(record.AlbumOwner),
				FolderName: catalog.NewFolderName(record.AlbumFolderName),
			},
			MediaCount: record.Count,
			Name:       record.AlbumName,
			Start:      start,
			End:        end,
		},
		Availability: availability,
	}, nil
}

func unmarshalAlbumCovers(item map[string]types.AttributeValue) (catalog.AlbumId, catalogviews.Availability, []catalog.Cover, error) {
	record := &AlbumCoversRecord{}
	if err := attributevalue.UnmarshalMap(item, record); err != nil {
		return catalog.AlbumId{}, catalogviews.Availability{}, nil, errors.Wrapf(err, "failed to unmarshal album covers record: %+v", item)
	}
	availability, err := unmarshalAvailability(record.AvailabilityType, record.UserId, item)
	if err != nil {
		return catalog.AlbumId{}, catalogviews.Availability{}, nil, err
	}
	albumId := catalog.AlbumId{
		Owner:      ownermodel.Owner(record.AlbumOwner),
		FolderName: catalog.NewFolderName(record.AlbumFolderName),
	}
	return albumId, availability, unmarshalCovers(record.Covers), nil
}

func unmarshalAvailability(availabilityType, userId string, item map[string]types.AttributeValue) (catalogviews.Availability, error) {
	switch availabilityType {
	case AvailabilityTypeOwner:
		return catalogviews.OwnerAvailability(usermodel.UserId(userId)), nil
	case AvailabilityTypeVisitor:
		return catalogviews.VisitorAvailability(usermodel.UserId(userId)), nil
	default:
		return catalogviews.Availability{}, errors.Errorf("unknown AvailabilityType %q on album record: %+v", availabilityType, item)
	}
}

func marshalTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

func unmarshalTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, value)
}
