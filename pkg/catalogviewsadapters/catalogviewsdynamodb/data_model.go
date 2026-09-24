package catalogviewsdynamodb

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"time"
)

const (
	AvailabilityTypeOwner   = "OWNED"
	AvailabilityTypeVisitor = "VISITOR"
)

type AlbumSummaryRecord struct {
	appdynamodb.TablePk
	AlbumOwner       string
	AlbumFolderName  string
	AvailabilityType string
	UserId           string
	Count            int
	AlbumName        string `dynamodbav:",omitempty"`
	AlbumStart       string `dynamodbav:",omitempty"`
	AlbumEnd         string `dynamodbav:",omitempty"`
	AlbumViewIndexPK string
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

func unmarshalAlbumSummary(item map[string]types.AttributeValue) (*catalogviews.UserAlbumSummary, error) {
	record := &AlbumSummaryRecord{}
	err := attributevalue.UnmarshalMap(item, record)

	availability := catalogviews.OwnerAvailability(usermodel.UserId(record.UserId))
	if record.AvailabilityType == AvailabilityTypeVisitor {
		availability = catalogviews.VisitorAvailability(usermodel.UserId(record.UserId))
	}

	start, _ := unmarshalTime(record.AlbumStart)
	end, _ := unmarshalTime(record.AlbumEnd)

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
	}, errors.Wrapf(err, "failed to unmarshal album summary record: %+v", item)
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
