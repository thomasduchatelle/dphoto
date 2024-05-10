package catalogdynamo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"sort"
)

func (r *Repository) FindAlbumsByOwner(ctx context.Context, owner catalog.Owner) ([]*catalog.Album, error) {
	expr, err := expression.NewBuilder().WithKeyCondition(expression.KeyAnd(
		expression.Key("PK").Equal(expression.Value(fmt.Sprintf("%s#ALBUM", owner))),
		expression.Key("SK").BeginsWith("ALBUM#"),
	)).Build()
	if err != nil {
		return nil, err
	}

	query := &dynamodb.QueryInput{
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
		KeyConditionExpression:    expr.KeyCondition(),
		TableName:                 &r.table,
	}
	data, err := r.client.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(err, "DynamoDb Query failed: %+v", expr)
	}

	albums := make([]*catalog.Album, 0, len(data.Items))
	for _, attributes := range data.Items {
		album, err := unmarshalAlbum(attributes)
		if err != nil {
			return nil, err
		}

		// TODO Media should be counted when inserted, not at every request
		count, err := r.CountMedias(ctx, catalog.AlbumId{Owner: owner, FolderName: album.FolderName})
		if err != nil {
			return nil, errors.Wrapf(err, "couldn't count medias in album %s%s", owner, album.FolderName)
		}

		album.TotalCount = count
		albums = append(albums, album)
	}

	// TODO Is it really necessary to sort here?
	sort.Slice(albums, func(i, j int) bool {
		return albums[i].Start.Before(albums[j].Start)
	})

	return albums, nil
}

func (r *Repository) UpdateAlbumName(ctx context.Context, albumId catalog.AlbumId, newName string) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) InsertAlbum(ctx context.Context, album catalog.Album) error {
	if album.Owner == "" || album.FolderName == "" {
		return errors.Errorf("Owner and Foldername are mandatory")
	}

	item, err := marshalAlbum(&album)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
		Item:                item,
		TableName:           &r.table,
	})

	return errors.WithStack(errors.Wrapf(err, "failed inserting album '%s'", album.FolderName))
}

func (r *Repository) DeleteEmptyAlbum(ctx context.Context, albumId catalog.AlbumId) error {
	// TODO DeleteEmptyAlbum method should not exists, DeleteAlbum should be used instead (no count check)
	count, err := r.CountMedias(ctx, albumId)
	if err != nil {
		return errors.Wrapf(err, "failed to count number of medias in album %s", albumId.FolderName)
	}

	if count > 0 {
		return catalog.AlbumIsNotEmptyError
	}

	return r.DeleteAlbum(ctx, albumId)
}

func (r *Repository) DeleteAlbum(ctx context.Context, albumId catalog.AlbumId) error {
	primaryKey, err := attributevalue.MarshalMap(AlbumPrimaryKey(albumId.Owner, albumId.FolderName))
	if err != nil {
		return err
	}
	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       primaryKey,
		TableName: &r.table,
	})
	return err
}

// CountMedias provides an accurate number of medias and can be used to update the count stored in the album record
func (r *Repository) CountMedias(ctx context.Context, albumId catalog.AlbumId) (int, error) {
	expr, err := expression.NewBuilder().WithKeyCondition(expression.KeyAnd(
		withinAlbum(albumId.Owner, albumId.FolderName),
		withExcludingMetaRecord(),
	)).Build()
	if err != nil {
		return 0, err
	}

	query, err := r.client.Query(ctx, &dynamodb.QueryInput{
		ExpressionAttributeValues: expr.Values(),
		ExpressionAttributeNames:  expr.Names(),
		IndexName:                 aws.String(albumIndex),
		KeyConditionExpression:    expr.KeyCondition(),
		Select:                    types.SelectCount,
		TableName:                 &r.table,
	})
	if err != nil {
		return 0, err
	}
	return int(query.Count), nil
}

func (r *Repository) CountMediasBySelectors(ctx context.Context, owner catalog.Owner, selectors []catalog.MediaSelector) (int, error) {
	if len(selectors) == 0 {
		return 0, nil
	}

	request := r.convertSelectorsIntoMediaRequest(owner, selectors)
	queries, err := newMediaQueryBuilders(r.table, request, "", types.SelectCount)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, query := range queries {
		output, err := r.client.Query(ctx, query)
		if err != nil {
			return 0, err
		}
		count += int(output.Count)
	}

	return count, nil
}

func (r *Repository) FindAlbumById(ctx context.Context, id catalog.AlbumId) (*catalog.Album, error) {
	albums, err := r.FindAlbumByIds(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(albums) == 0 {
		return nil, catalog.AlbumNotFoundError
	}
	return albums[0], nil
}

func (r *Repository) FindAlbumByIds(ctx context.Context, ids ...catalog.AlbumId) ([]*catalog.Album, error) {
	var keys []map[string]types.AttributeValue
	for _, id := range ids {
		key, err := attributevalue.MarshalMap(AlbumPrimaryKey(id.Owner, id.FolderName))
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	stream := dynamoutils.NewGetStream(ctx, dynamoutils.NewGetBatchItem(r.client, r.table, ""), keys, dynamoutils.DynamoReadBatchSize)
	var albums []*catalog.Album
	for stream.HasNext() {
		album, err := unmarshalAlbum(stream.Next())
		if err != nil {
			return nil, err
		}

		// TODO Media should be counted when inserted, not at every request
		album.TotalCount, err = r.CountMedias(ctx, album.AlbumId)
		if err != nil {
			return nil, errors.Wrapf(err, "couldn't count medias in album %s%s", album.Owner, album.FolderName)
		}

		albums = append(albums, album)
	}

	return albums, stream.Error()
}

func (r *Repository) UpdateAlbum(ctx context.Context, album catalog.Album) error {
	item, err := marshalAlbum(&album)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_exists(PK)"),
		Item:                item,
		TableName:           &r.table,
	})

	return errors.Wrapf(err, "failed updating album '%s'", album.FolderName)
}
