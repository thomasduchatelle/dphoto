package catalogviewsdynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

type AlbumViewRepository struct {
	Client    *dynamodb.Client
	TableName string
}

func (a *AlbumViewRepository) InsertAlbumSize(ctx context.Context, albumSizes []catalogviews.AlbumSize) error {
	var items []types.WriteRequest

	for _, albumSize := range albumSizes {
		sizes, err := marshalAlbumSize(albumSize)
		if err != nil {
			return err
		}

		for _, size := range sizes {
			items = append(items, types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: size,
				},
			})
		}
	}

	if len(items) > 0 {
		err := dynamoutils.BufferedWriteItems(ctx, a.Client, items, a.TableName, dynamoutils.DynamoWriteBatchSize)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO GetAvailabilitiesByUser shouldn't return a slice of AlbumSize: the AlbumSize.User field is not used

func (a *AlbumViewRepository) GetAvailabilitiesByUser(ctx context.Context, user usermodel.UserId) ([]catalogviews.AlbumSize, error) {
	//TODO implement me
	panic("implement me")
}
