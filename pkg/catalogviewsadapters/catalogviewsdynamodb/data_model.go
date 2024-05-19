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
)

type AlbumSizeRecord struct {
	appdynamodb.TablePk
	AlbumOwner      string
	AlbumFolderName string
	Count           int
}

func albumsViewPK(user usermodel.UserId) string {
	return fmt.Sprintf("USER#%s#ALBUMS_VIEW", user.Value())
}

func albumSizeKey(user catalogviews.Availability, albumId catalog.AlbumId) appdynamodb.TablePk {
	belongType := "OWNED"
	if !user.AsOwner {
		belongType = "VISITOR"
	}

	recordKey := appdynamodb.TablePk{
		PK: albumsViewPK(user.UserId),
		SK: fmt.Sprintf("%s#%s#%s#COUNT", belongType, albumId.Owner, albumId.FolderName.String()),
	}
	return recordKey
}

func marshalAlbumSize(albumSize catalogviews.AlbumSize) ([]map[string]types.AttributeValue, error) {
	var items []map[string]types.AttributeValue
	for _, user := range albumSize.Users {
		recordKey := albumSizeKey(user, albumSize.AlbumId)

		item, err := attributevalue.MarshalMap(AlbumSizeRecord{
			TablePk:         recordKey,
			AlbumOwner:      albumSize.AlbumId.Owner.Value(),
			AlbumFolderName: albumSize.AlbumId.FolderName.String(),
			Count:           albumSize.MediaCount,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal album size record: %+v", albumSize)
		}

		items = append(items, item)
	}

	return items, nil
}

func unmarshalAlbumSize(item map[string]types.AttributeValue) (*catalogviews.AlbumSize, error) {
	record := &AlbumSizeRecord{}
	err := attributevalue.UnmarshalMap(item, record)

	return &catalogviews.AlbumSize{
		AlbumId: catalog.AlbumId{
			Owner:      ownermodel.Owner(record.AlbumOwner),
			FolderName: catalog.NewFolderName(record.AlbumFolderName),
		},
		MediaCount: record.Count,
	}, errors.Wrapf(err, "failed to unmarshal album size record: %+v", item)
}
