package catalogdynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	dynamoutils2 "github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"sort"
)

func (r *Repository) FindAlbumsByOwner(owner string) ([]*catalog.Album, error) {
	query := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":owner":     mustAttribute(fmt.Sprintf("%s#ALBUM", owner)),
			":albumOnly": mustAttribute("ALBUM#"),
		},
		KeyConditionExpression: aws.String("PK = :owner AND begins_with(SK, :albumOnly)"),
		TableName:              &r.table,
	}
	data, err := r.db.Query(query)
	if err != nil {
		return nil, errors.Wrapf(err, "DynamoDb Query failed: %s", query)
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

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
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

	primaryKey, err := dynamodbattribute.MarshalMap(AlbumPrimaryKey(owner, folderName))
	if err != nil {
		return err
	}
	_, err = r.db.DeleteItem(&dynamodb.DeleteItemInput{
		Key:       primaryKey,
		TableName: &r.table,
	})
	return err
}

// CountMedias provides an accurate number of medias and can be used to update the count stored in the album record
func (r *Repository) CountMedias(owner string, folderName string) (int, error) {
	albumIndexKey, err := dynamodbattribute.MarshalMap(AlbumIndexedKey(owner, folderName))
	if err != nil {
		return 0, err
	}

	query, err := r.db.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":albumKey":    albumIndexKey["AlbumIndexPK"],
			":excludeMeta": mustAttribute("$"),
		},
		IndexName:              aws.String("AlbumIndex"),
		KeyConditionExpression: aws.String("AlbumIndexPK = :albumKey AND AlbumIndexSK >= :excludeMeta"),
		Select:                 aws.String(dynamodb.SelectCount),
		TableName:              &r.table,
	})
	if err != nil {
		return 0, err
	}
	return int(*query.Count), nil
}

func (r *Repository) FindAlbums(ids ...catalog.AlbumId) ([]*catalog.Album, error) {
	var keys []map[string]*dynamodb.AttributeValue
	for _, id := range ids {
		key, err := dynamodbattribute.MarshalMap(AlbumPrimaryKey(id.Owner, id.FolderName))
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	stream := dynamoutils2.NewGetStream(dynamoutils2.NewGetBatchItem(r.db, r.table, ""), keys, dynamoutils2.DynamoReadBatchSize)
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

	_, err = r.db.PutItem(&dynamodb.PutItemInput{
		ConditionExpression: aws.String("attribute_exists(PK)"),
		Item:                item,
		TableName:           &r.table,
	})

	return errors.WithStack(errors.Wrapf(err, "failed updating album '%s'", album.FolderName))
}
