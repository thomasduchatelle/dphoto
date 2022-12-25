package catalogdynamo

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogadapters/dynamoutils"
	"strings"
)

func (r *Repository) InsertMedias(owner string, medias []catalog.CreateMediaRequest) error {
	requests := make([]*dynamodb.WriteRequest, len(medias))
	for index, media := range medias {
		mediaEntry, err := marshalMedia(owner, &media)
		if err != nil {
			return errors.Wrapf(err, "Failed mapping media %s", fmt.Sprint(media))
		}

		requests[index] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: mediaEntry,
			},
		}
	}

	return dynamoutils.BufferedWriteItems(r.db, requests, r.table, DynamoWriteBatchSize)
}

func (r *Repository) FindMedias(request *catalog.FindMediaRequest) ([]*catalog.MediaMeta, error) {
	queries, err := newMediaQueryBuilders(r.table, request, "")
	if err != nil {
		return nil, err
	}

	var medias []*catalog.MediaMeta

	crawler := dynamoutils.NewQueryStream(r.db, queries)
	for crawler.HasNext() {
		media, err := unmarshalMediaMetaData(crawler.Next())
		if err != nil {
			return nil, err
		}

		medias = append(medias, media)
	}

	return medias, err
}

func (r *Repository) FindMediaCurrentAlbum(owner, mediaId string) (string, error) {
	key, err := dynamodbattribute.MarshalMap(MediaPrimaryKey(owner, mediaId))
	if err != nil {
		return "", errors.Wrapf(err, "failed to marshal media key %s/%s", owner, mediaId)
	}

	item, err := r.db.GetItem(&dynamodb.GetItemInput{
		Key:                  key,
		ProjectionExpression: aws.String("AlbumIndexPK"),
		TableName:            &r.table,
	})
	if err != nil {
		return "", errors.Wrapf(err, "couldn't get media metadata for media %+v", key)
	}

	if len(item.Item) == 0 {
		return "", catalog.NotFoundError
	}

	if albumIndexPk, ok := item.Item["AlbumIndexPK"]; ok && strings.HasPrefix(*albumIndexPk.S, owner) {
		return (*albumIndexPk.S)[len(owner)+1:], nil
	}

	return "", errors.Errorf("invalid AlbumIndexPK format expected to start with %s ; value: %+v", owner, item.Item)
}

func (r *Repository) FindMediaIds(request *catalog.FindMediaRequest) ([]string, error) {
	queries, err := newMediaQueryBuilders(r.table, request, "Id")
	if err != nil {
		return nil, err
	}

	var mediaIds []string

	crawler := dynamoutils.NewQueryStream(r.db, queries)
	for crawler.HasNext() {
		record := crawler.Next()
		mediaIds = append(mediaIds, *record["Id"].S)
	}

	return mediaIds, err
}

func (r *Repository) TransferMedias(owner string, mediaIds []string, newFolderName string) error {
	for _, id := range mediaIds {
		mediaKey, err := dynamodbattribute.MarshalMap(MediaPrimaryKey(owner, id))
		if err != nil {
			return err
		}

		updateValues, err := dynamodbattribute.MarshalMap(map[string]string{
			":albumPK": AlbumIndexedKeyPK(owner, newFolderName),
		})
		if err != nil {
			return err
		}

		_, err = r.db.UpdateItem(&dynamodb.UpdateItemInput{
			ExpressionAttributeValues: updateValues,
			Key:                       mediaKey,
			TableName:                 &r.table,
			UpdateExpression:          aws.String("SET AlbumIndexPK=:albumPK"),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) FindExistingSignatures(owner string, signatures []*catalog.MediaSignature) ([]*catalog.MediaSignature, error) {
	// note: this implementation expects media id to be an encoded version of its signature

	keys := make([]map[string]*dynamodb.AttributeValue, len(signatures))
	for index, signature := range signatures {
		id, err := catalog.GenerateMediaId(*signature)
		if err != nil {
			return nil, err
		}

		key, err := dynamodbattribute.MarshalMap(MediaPrimaryKey(owner, id))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal media keys from signature %+v", signature)
		}

		keys[index] = key
	}

	stream := dynamoutils.NewGetStream(dynamoutils.NewGetBatchItem(r.db, r.table, *aws.String("Id")), keys, DynamoReadBatchSize)

	found := make([]*catalog.MediaSignature, 0, len(signatures))
	for stream.HasNext() {
		attributes := stream.Next()
		if awsAttr, ok := attributes["Id"]; ok && awsAttr.S != nil {
			signature, err := catalog.DecodeMediaId(*awsAttr.S)
			if err != nil {
				return nil, err
			}

			found = append(found, signature)
		} else {
			return nil, errors.Errorf("Records doesn't have an 'Id' field: %+v", attributes)
		}
	}

	return found, stream.Error()
}

// unmarshalPageToken take a base64 page id and decode it into a dynamodb key. Return nil, nil if nextPageToken is blank
func unmarshalPageToken(nextPageToken string) (map[string]*dynamodb.AttributeValue, error) {
	var startKey map[string]*dynamodb.AttributeValue
	if !isBlank(nextPageToken) {
		startKeyJson, err := base64.StdEncoding.DecodeString(nextPageToken)
		if err != nil {
			return startKey, errors.Wrapf(err, "Invalid nextPageToken %s", nextPageToken)
		}

		err = json.Unmarshal(startKeyJson, &startKey)
		if err != nil {
			return startKey, errors.Wrapf(err, "Invalid nextPageToken %s", startKeyJson)
		}
	}

	return startKey, nil
}

func (r *Repository) unmarshalPageToken(nextPageToken string) (map[string]*dynamodb.AttributeValue, error) {
	return unmarshalPageToken(nextPageToken)
}

// marshalPageToken take a dynamodb key and encode it (JSON+BASE64) to be used by service. Return empty string when key is empty
func (r *Repository) marshalPageToken(key map[string]*dynamodb.AttributeValue) (string, error) {
	if len(key) == 0 {
		return "", nil
	}

	keyJson, err := json.Marshal(key)
	return base64.StdEncoding.EncodeToString(keyJson), err
}

// panic if value can't be converted into an attribute
func mustAttribute(value interface{}) *dynamodb.AttributeValue {
	attribute, err := dynamodbattribute.Marshal(value)
	if err != nil {
		panic(err)
	}
	return attribute
}
