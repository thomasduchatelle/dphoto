package catalogviewsdynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
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

func (a *AlbumViewRepository) DeleteAlbumSize(ctx context.Context, availability catalogviews.Availability, albumId catalog.AlbumId) error {
	key := albumSizeKey(availability, albumId)

	_, err := a.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &a.TableName,
		Key:       key.ToAttributes(),
	})

	return errors.Wrapf(err, "failed to delete album size for album %v and user %v", albumId, availability)
}

// TODO GetAvailabilitiesByUser shouldn't return a slice of AlbumSize: the AlbumSize.User field is not used

func (a *AlbumViewRepository) GetAvailabilitiesByUser(ctx context.Context, user usermodel.UserId) ([]catalogviews.AlbumSize, error) {
	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.Key("PK").Equal(expression.Value(albumsViewPK(user)))).
		Build()

	if err != nil {
		return nil, errors.Wrapf(err, "failed to build expression for user %v", user)
	}

	paginator := dynamodb.NewQueryPaginator(a.Client, &dynamodb.QueryInput{
		TableName:                 &a.TableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})

	var albumSizes []catalogviews.AlbumSize
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.Items {
			albumSize, err := unmarshalAlbumSize(item)
			if err != nil {
				return nil, err
			}

			albumSizes = append(albumSizes, *albumSize)
		}
	}

	return albumSizes, nil
}
