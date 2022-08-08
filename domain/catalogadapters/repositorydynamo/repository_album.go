package repositorydynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/domain/catalog"
	"sort"
)

func (r *rep) FindAllAlbums(owner string) ([]*catalog.Album, error) {
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

func (r *rep) InsertAlbum(album catalog.Album) error {
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

func (r *rep) DeleteEmptyAlbum(owner string, folderName string) error {
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
func (r *rep) CountMedias(owner string, folderName string) (int, error) {
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

func (r *rep) FindAlbum(owner string, folderName string) (*catalog.Album, error) {
	key, err := dynamodbattribute.MarshalMap(AlbumPrimaryKey(owner, folderName))
	if err != nil {
		return nil, err
	}

	attributes, err := r.db.GetItem(&dynamodb.GetItemInput{
		Key:       key,
		TableName: &r.table,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find album '%s'", folderName)
	}
	if len(attributes.Item) == 0 {
		return nil, catalog.NotFoundError
	}

	return unmarshalAlbum(attributes.Item)
}

func (r *rep) UpdateAlbum(album catalog.Album) error {
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
