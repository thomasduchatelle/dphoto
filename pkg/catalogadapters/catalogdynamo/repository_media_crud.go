package catalogdynamo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"strings"
)

func (r *Repository) InsertMedias(ctx context.Context, owner ownermodel.Owner, medias []catalog.CreateMediaRequest) error {
	requests := make([]types.WriteRequest, len(medias))
	for index, media := range medias {
		mediaEntry, err := marshalMedia(owner, &media)
		if err != nil {
			return errors.Wrapf(err, "Failed mapping media %s", fmt.Sprint(media))
		}

		requests[index] = types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: mediaEntry,
			},
		}
	}

	return dynamoutils.BufferedWriteItems(ctx, r.client, requests, r.table, dynamoutils.DynamoWriteBatchSize)
}

func (r *Repository) FindMedias(ctx context.Context, request *catalog.FindMediaRequest) ([]*catalog.MediaMeta, error) {
	queries, err := newMediaQueryBuilders(r.table, request, "", types.SelectAllAttributes)
	if err != nil {
		return nil, err
	}

	var medias []*catalog.MediaMeta

	crawler := dynamoutils.NewQueryStream(ctx, r.client, queries)
	for crawler.HasNext() {
		media, err := unmarshalMediaMetaData(crawler.Next())
		if err != nil {
			return nil, err
		}

		medias = append(medias, media)
	}

	return medias, err
}

func (r *Repository) FindMediaCurrentAlbum(ctx context.Context, owner ownermodel.Owner, mediaId catalog.MediaId) (*catalog.AlbumId, error) {
	key, err := attributevalue.MarshalMap(MediaPrimaryKey(owner, mediaId))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal media key %s/%s", owner, mediaId)
	}

	item, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key:                  key,
		ProjectionExpression: aws.String("AlbumIndexPK"),
		TableName:            &r.table,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't get media metadata for media %+v", key)
	}

	if len(item.Item) == 0 {
		return nil, catalog.AlbumNotFoundError
	}

	if albumIndexPk, ok := item.Item["AlbumIndexPK"]; ok {
		if value, ok := albumIndexPk.(*types.AttributeValueMemberS); ok && strings.HasPrefix(value.Value, string(owner)) {
			return &catalog.AlbumId{
				Owner:      owner,
				FolderName: catalog.NewFolderName((value.Value)[len(owner)+1:]),
			}, nil
		}
	}

	return nil, errors.Errorf("invalid AlbumIndexPK format expected to start with %s ; value: %+v", owner, item.Item)
}

func (r *Repository) FindMediaIds(ctx context.Context, request *catalog.FindMediaRequest) ([]catalog.MediaId, error) {
	queries, err := newMediaQueryBuilders(r.table, request, "Id", types.SelectAllAttributes)
	if err != nil {
		return nil, err
	}

	var mediaIds []catalog.MediaId

	crawler := dynamoutils.NewQueryStream(ctx, r.client, queries)
	for crawler.HasNext() {
		record := crawler.Next()
		if id, ok := record["Id"].(*types.AttributeValueMemberS); ok {
			mediaIds = append(mediaIds, catalog.MediaId(id.Value))
		}
	}

	return mediaIds, err
}

func (r *Repository) FindExistingSignatures(ctx context.Context, owner ownermodel.Owner, signatures []*catalog.MediaSignature) ([]*catalog.MediaSignature, error) {
	// note: this implementation expects media id to be an encoded version of its signature

	var keys []map[string]types.AttributeValue
	uniqueSignatures := make(map[catalog.MediaSignature]interface{})
	for _, signature := range signatures {
		if _, found := uniqueSignatures[*signature]; !found {
			id, err := catalog.GenerateMediaId(*signature)
			if err != nil {
				return nil, err
			}

			key, err := attributevalue.MarshalMap(MediaPrimaryKey(owner, id))
			if err != nil {
				return nil, errors.Wrapf(err, "failed to marshal media keys from signature %+v", signature)
			}

			keys = append(keys, key)
		}
		uniqueSignatures[*signature] = nil
	}

	stream := dynamoutils.NewGetStream(ctx, dynamoutils.NewGetBatchItem(r.client, r.table, *aws.String("Id")), keys, dynamoutils.DynamoReadBatchSize)

	found := make([]*catalog.MediaSignature, 0, len(signatures))
	for stream.HasNext() {
		attributes := stream.Next()
		if awsAttr, ok := attributes["Id"]; ok {
			if value, ok := awsAttr.(*types.AttributeValueMemberS); ok && value.Value != "" {
				signature, err := catalog.DecodeMediaId(catalog.MediaId(value.Value))
				if err != nil {
					return nil, err
				}

				found = append(found, signature)
			} else {
				return nil, errors.Errorf("Records Id field is empty or not a String: %+v", attributes)
			}
		} else {
			return nil, errors.Errorf("Records doesn't have an 'Id' field: %+v", attributes)
		}
	}

	return found, stream.Error()
}

func (r *Repository) FindSignatures(ctx context.Context, owner ownermodel.Owner, signatures []catalog.MediaSignature) (map[catalog.MediaSignature]catalog.MediaId, error) {
	// note: this implementation expects media id to be an encoded version of its signature

	var keys []map[string]types.AttributeValue
	uniqueSignatures := make(map[catalog.MediaSignature]interface{})
	for _, signature := range signatures {
		if _, found := uniqueSignatures[signature]; !found {
			id, err := catalog.GenerateMediaId(signature)
			if err != nil {
				return nil, err
			}

			key, err := attributevalue.MarshalMap(MediaPrimaryKey(owner, id))
			if err != nil {
				return nil, errors.Wrapf(err, "failed to marshal media keys from signature %+v", signature)
			}

			keys = append(keys, key)
		}
		uniqueSignatures[signature] = nil
	}

	stream := dynamoutils.NewGetStream(ctx, dynamoutils.NewGetBatchItem(r.client, r.table, *aws.String("Id")), keys, dynamoutils.DynamoReadBatchSize)

	found := make(map[catalog.MediaSignature]catalog.MediaId)
	for stream.HasNext() {
		attributes := stream.Next()
		if awsAttr, ok := attributes["Id"]; ok {
			if value, ok := awsAttr.(*types.AttributeValueMemberS); ok && value.Value != "" {
				mediaId := catalog.MediaId(value.Value)
				signature, err := catalog.DecodeMediaId(mediaId)
				if err != nil {
					return nil, err
				}

				found[*signature] = mediaId
			} else {
				return nil, errors.Errorf("Records Id field is empty or not a String: %+v", attributes)
			}
		} else {
			return nil, errors.Errorf("Records doesn't have an 'Id' field: %+v", attributes)
		}
	}

	return found, stream.Error()
}
