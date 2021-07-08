package dynamo

import (
	"duchatelle.io/dphoto/dphoto/catalog"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"sort"
)

func (r *Rep) FindAllAlbums() ([]*catalog.Album, error) {
	query := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":owner":     mustAttribute(r.RootOwner),
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
		unmarshalAlbum, err := unmarshalAlbum(attributes)
		if err != nil {
			return nil, err
		}

		albums = append(albums, unmarshalAlbum)
	}

	sort.Slice(albums, func(i, j int) bool {
		return albums[i].Start.Before(albums[j].Start)
	})

	return albums, nil
}

func (r *Rep) InsertAlbum(album catalog.Album) error {
	item, err := marshalAlbum(r.RootOwner, &album)
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

func (r *Rep) DeleteEmptyAlbum(folderName string) error {
	count, err := r.CountMedias(folderName)
	if err != nil {
		return errors.Wrapf(err, "failed to count number of medias in album %s", folderName)
	}

	if count > 0 {
		return catalog.NotEmptyError
	}

	primaryKey, err := dynamodbattribute.MarshalMap(albumPrimaryKey(r.RootOwner, folderName))
	if err != nil {
		return err
	}
	_, err = r.db.DeleteItem(&dynamodb.DeleteItemInput{
		Key:       primaryKey,
		TableName: &r.table,
	})
	return err
}

func (r *Rep) CountMedias(folderName string) (int, error) {
	albumIndexKey, err := dynamodbattribute.MarshalMap(albumIndexedKey(r.RootOwner, folderName))
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

func (r *Rep) FindAlbum(folderName string) (*catalog.Album, error) {
	key, err := dynamodbattribute.MarshalMap(albumPrimaryKey(r.RootOwner, folderName))
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

func (r *Rep) UpdateAlbum(album catalog.Album) error {
	item, err := marshalAlbum(r.RootOwner, &album)
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
