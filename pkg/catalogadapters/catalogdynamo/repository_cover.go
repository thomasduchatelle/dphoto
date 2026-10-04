package catalogdynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
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

// FindCoversByAlbums batches a BatchGetItem over the canonical cover record of each
// album. Missing covers are omitted from the result map (an album with no cover record
// is simply not present in the output).
func (r *Repository) FindCoversByAlbums(ctx context.Context, albumIds ...catalog.AlbumId) (map[catalog.AlbumId][]catalog.Cover, error) {
	if len(albumIds) == 0 {
		return nil, nil
	}

	keys := make([]map[string]types.AttributeValue, 0, len(albumIds))
	for _, id := range albumIds {
		key, err := attributevalue.MarshalMap(CoverPrimaryKey(id.Owner, id.FolderName))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal cover key for album %s", id)
		}
		keys = append(keys, key)
	}

	result := make(map[catalog.AlbumId][]catalog.Cover)
	stream := dynamoutils.NewGetStream(ctx, dynamoutils.NewGetBatchItem(r.client, r.table, ""), keys, dynamoutils.DynamoReadBatchSize)
	for stream.HasNext() {
		item := stream.Next()
		covers, err := unmarshalCovers(item)
		if err != nil {
			return nil, err
		}
		albumId, err := unmarshalCoverAlbumId(item)
		if err != nil {
			return nil, err
		}
		result[albumId] = covers
	}
	return result, stream.Error()
}

func unmarshalCoverAlbumId(item map[string]types.AttributeValue) (catalog.AlbumId, error) {
	var record CoverRecord
	if err := attributevalue.UnmarshalMap(item, &record); err != nil {
		return catalog.AlbumId{}, errors.Wrapf(err, "failed to unmarshal cover record for AlbumId: %+v", item)
	}
	return catalog.AlbumId{
		Owner:      ownermodel.Owner(record.AlbumOwner),
		FolderName: catalog.NewFolderName(record.AlbumFolderName),
	}, nil
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
