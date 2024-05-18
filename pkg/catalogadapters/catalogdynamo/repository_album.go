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
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"time"
)

func (r *Repository) FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
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

		albums = append(albums, album)
	}

	return albums, nil
}

func (r *Repository) UpdateAlbumName(ctx context.Context, albumId catalog.AlbumId, newName string) error {
	update, err := expression.NewBuilder().
		WithUpdate(expression.Set(expression.Name("AlbumName"), expression.Value(newName))).
		WithCondition(expression.Name("PK").Equal(expression.Value(AlbumPrimaryKey(albumId.Owner, albumId.FolderName).PK))).
		Build()
	if err != nil {
		return errors.Wrapf(err, "failed to build update name expression for album %s", albumId)
	}

	albumKey, err := attributevalue.MarshalMap(AlbumPrimaryKey(albumId.Owner, albumId.FolderName))
	if err != nil {
		return errors.Wrapf(err, "failed to build update name expression for album %s", albumId)
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       albumKey,
		TableName:                 &r.table,
		ConditionExpression:       update.Condition(),
		ExpressionAttributeNames:  update.Names(),
		ExpressionAttributeValues: update.Values(),
		UpdateExpression:          update.Update(),
	})
	var conditionalCheckFailedException *types.ConditionalCheckFailedException
	if errors.As(err, &conditionalCheckFailedException) {
		return catalog.MediaNotFoundError
	}
	return errors.Wrapf(err, "failed to build update name expression for album %s", albumId)
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

// countMedias provides an accurate number of medias and can be used to update the count stored in the album record
func (r *Repository) countMedias(ctx context.Context, albumId catalog.AlbumId) (int, error) {
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

func (r *Repository) CountMediasBySelectors(ctx context.Context, owner ownermodel.Owner, selectors []catalog.MediaSelector) (int, error) {
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

		albums = append(albums, album)
	}

	return albums, stream.Error()
}

func (r *Repository) AmendDates(ctx context.Context, albumId catalog.AlbumId, start, end time.Time) error {
	update, err := expression.NewBuilder().
		WithUpdate(expression.
			Set(expression.Name("AlbumStart"), expression.Value(start)).
			Set(expression.Name("AlbumEnd"), expression.Value(end))).
		WithCondition(expression.Name("PK").Equal(expression.Value(AlbumPrimaryKey(albumId.Owner, albumId.FolderName).PK))).
		Build()
	if err != nil {
		return errors.Wrapf(err, "failed to build date update [%s -> %s] expression for album %s", start, end, albumId)
	}

	albumKey, err := attributevalue.MarshalMap(AlbumPrimaryKey(albumId.Owner, albumId.FolderName))
	if err != nil {
		return errors.Wrapf(err, "failed to build date update [%s -> %s] expression for album %s", start, end, albumId)
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		Key:                       albumKey,
		TableName:                 &r.table,
		ConditionExpression:       update.Condition(),
		ExpressionAttributeNames:  update.Names(),
		ExpressionAttributeValues: update.Values(),
		UpdateExpression:          update.Update(),
	})
	var conditionalCheckFailedException *types.ConditionalCheckFailedException
	if errors.As(err, &conditionalCheckFailedException) {
		return catalog.AlbumNotFoundError
	}

	return errors.Wrapf(err, "failed to exec update [%s -> %s] expression for album %s", start, end, albumId)
}

func (r *Repository) CountMedia(ctx context.Context, album ...catalog.AlbumId) (map[catalog.AlbumId]int, error) {
	albumCount := make(map[catalog.AlbumId]int)

	for _, albumId := range album {
		count, err := r.countMedias(ctx, albumId)
		if err != nil {
			return nil, err
		}
		albumCount[albumId] = count
	}

	return albumCount, nil
}
