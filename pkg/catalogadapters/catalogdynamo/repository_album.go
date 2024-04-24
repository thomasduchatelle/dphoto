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
	dynamoutils "github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutilsv2"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"sort"
)

func (r *Repository) FindAlbumsByOwner(owner string) ([]*catalog.Album, error) {
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
	data, err := r.client.Query(context.TODO(), query)
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
		count, err := r.CountMedias(owner, album.FolderName)
		if err != nil {
			return nil, errors.Wrapf(err, "couldn't count medias in album %s%s", owner, album.FolderName)
		}

		album.TotalCount = count
		albums = append(albums, album)
	}

	sort.Slice(albums, func(i, j int) bool {
		return albums[i].Start.Before(albums[j].Start)
	})

	return albums, nil
}

func (r *Repository) InsertAlbum(album catalog.Album) error {
	if album.Owner == "" || album.FolderName == "" {
		return errors.Errorf("Owner and Foldername are mandatory")
	}

	item, err := marshalAlbum(&album)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_not_exists(PK)"),
		Item:                item,
		TableName:           &r.table,
	})

	return errors.WithStack(errors.Wrapf(err, "failed inserting album '%s'", album.FolderName))
}

func (r *Repository) DeleteEmptyAlbum(owner string, folderName string) error {
	count, err := r.CountMedias(owner, folderName)
	if err != nil {
		return errors.Wrapf(err, "failed to count number of medias in album %s", folderName)
	}

	if count > 0 {
		return catalog.NotEmptyError
	}

	primaryKey, err := attributevalue.MarshalMap(AlbumPrimaryKey(owner, folderName))
	if err != nil {
		return err
	}
	_, err = r.client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		Key:       primaryKey,
		TableName: &r.table,
	})
	return err
}

// CountMedias provides an accurate number of medias and can be used to update the count stored in the album record
func (r *Repository) CountMedias(owner string, folderName string) (int, error) {
	expr, err := expression.NewBuilder().WithKeyCondition(expression.KeyAnd(
		withinAlbum(owner, folderName),
		withExcludingMetaRecord(),
	)).Build()
	if err != nil {
		return 0, err
	}

	query, err := r.client.Query(context.TODO(), &dynamodb.QueryInput{
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

func (r *Repository) FindAlbums(ids ...catalog.AlbumId) ([]*catalog.Album, error) {
	var keys []map[string]types.AttributeValue
	for _, id := range ids {
		key, err := attributevalue.MarshalMap(AlbumPrimaryKey(id.Owner, id.FolderName))
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	stream := dynamoutils.NewGetStream(context.TODO(), dynamoutils.NewGetBatchItem(r.client, r.table, ""), keys, dynamoutils.DynamoReadBatchSize)
	var albums []*catalog.Album
	for stream.HasNext() {
		album, err := unmarshalAlbum(stream.Next())
		if err != nil {
			return nil, err
		}

		// TODO Media should be counted when inserted, not at every request
		album.TotalCount, err = r.CountMedias(album.Owner, album.FolderName)
		if err != nil {
			return nil, errors.Wrapf(err, "couldn't count medias in album %s%s", album.Owner, album.FolderName)
		}

		albums = append(albums, album)
	}

	return albums, stream.Error()
}

func (r *Repository) UpdateAlbum(album catalog.Album) error {
	item, err := marshalAlbum(&album)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_exists(PK)"),
		Item:                item,
		TableName:           &r.table,
	})

	return errors.Wrapf(err, "failed updating album '%s'", album.FolderName)
}
