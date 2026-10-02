package catalogdynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

func (r *Repository) FindCoversByAlbum(ctx context.Context, albumId catalog.AlbumId) ([]catalog.Cover, error) {
	key, err := attributevalue.MarshalMap(CoverPrimaryKey(albumId.Owner, albumId.FolderName))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal cover key for album %s", albumId)
	}

	output, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       key,
		TableName: &r.table,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load covers for album %s", albumId)
	}
	if len(output.Item) == 0 {
		return nil, nil
	}

	return unmarshalCovers(output.Item)
}

func (r *Repository) SaveCovers(ctx context.Context, albumId catalog.AlbumId, covers []catalog.Cover) error {
	if len(covers) > catalog.MaxCoversPerAlbum {
		return errors.Wrapf(catalog.TooManyCoversErr, "cannot save %d covers on album %s", len(covers), albumId)
	}

	if len(covers) == 0 {
		key, err := attributevalue.MarshalMap(CoverPrimaryKey(albumId.Owner, albumId.FolderName))
		if err != nil {
			return errors.Wrapf(err, "failed to marshal cover key for album %s", albumId)
		}
		_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			Key:       key,
			TableName: &r.table,
		})
		return errors.Wrapf(err, "failed to delete covers for album %s", albumId)
	}

	item, err := marshalCovers(albumId, covers)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: &r.table,
	})
	return errors.Wrapf(err, "failed to save %d covers on album %s", len(covers), albumId)
}
