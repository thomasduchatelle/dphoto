package catalogdynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

const mediaByTagIndex = "MediaByTagIndex"

func (r *Repository) AddTag(ctx context.Context, albumId catalog.AlbumId, mediaId catalog.MediaId, tag catalog.Tag) error {
	item, err := marshalTag(albumId, mediaId, tag)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal tag %s for media %s/%s", tag, albumId.Owner, mediaId)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: &r.table,
	})
	return errors.Wrapf(err, "failed to add tag %s to media %s/%s", tag, albumId.Owner, mediaId)
}

func (r *Repository) RemoveTag(ctx context.Context, owner ownermodel.Owner, mediaId catalog.MediaId, tag catalog.Tag) error {
	key, err := attributevalue.MarshalMap(TagPrimaryKey(owner, mediaId, tag))
	if err != nil {
		return errors.Wrapf(err, "failed to marshal tag key %s for media %s/%s", tag, owner, mediaId)
	}

	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:       key,
		TableName: &r.table,
	})
	return errors.Wrapf(err, "failed to remove tag %s from media %s/%s", tag, owner, mediaId)
}

// FindMediasByTag queries the MediaByTagIndex to list the medias carrying the tag in
// O(matches), then reads each matched media's partition to attach its metadata and its full
// tag set, so each returned media is ready for the Tags page.
func (r *Repository) FindMediasByTag(ctx context.Context, owner ownermodel.Owner, tag catalog.Tag) ([]*catalog.MediaWithAlbum, error) {
	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.Key("MediaByTagPK").Equal(expression.Value(MediaByTagIndexedKey(owner, "", tag).MediaByTagPK))).
		Build()
	if err != nil {
		return nil, err
	}

	stream := dynamoutils.NewQueryStream(ctx, r.client, []*dynamodb.QueryInput{{
		TableName:                 &r.table,
		IndexName:                 aws.String(mediaByTagIndex),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}})

	var mediaIds []catalog.MediaId
	albumByMedia := make(map[catalog.MediaId]catalog.AlbumId)
	for stream.HasNext() {
		mediaId, _, albumId, err := unmarshalTag(stream.Next())
		if err != nil {
			return nil, err
		}
		mediaIds = append(mediaIds, mediaId)
		albumByMedia[mediaId] = albumId
	}
	if err := stream.Error(); err != nil {
		return nil, err
	}

	medias := make([]*catalog.MediaWithAlbum, 0, len(mediaIds))
	for _, mediaId := range mediaIds {
		media, err := r.findMediaWithTags(ctx, owner, mediaId)
		if err != nil {
			return nil, err
		}
		if media == nil {
			continue
		}
		medias = append(medias, &catalog.MediaWithAlbum{
			MediaMeta: *media,
			AlbumId:   albumByMedia[mediaId],
		})
	}

	return medias, nil
}

// findMediaWithTags reads a single media partition and returns the media metadata with all
// its tags attached. It returns nil when the media partition holds no #METADATA row.
func (r *Repository) findMediaWithTags(ctx context.Context, owner ownermodel.Owner, mediaId catalog.MediaId) (*catalog.MediaMeta, error) {
	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.Key("PK").Equal(expression.Value(MediaPrimaryKey(owner, mediaId).PK))).
		Build()
	if err != nil {
		return nil, err
	}

	stream := dynamoutils.NewQueryStream(ctx, r.client, []*dynamodb.QueryInput{{
		TableName:                 &r.table,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}})

	var media *catalog.MediaMeta
	var tags []catalog.Tag
	for stream.HasNext() {
		item := stream.Next()
		if isTagRecord(item) {
			_, tag, _, err := unmarshalTag(item)
			if err != nil {
				return nil, err
			}
			tags = append(tags, tag)
			continue
		}
		if isMediaMetadata(item) {
			media, err = unmarshalMediaMetaData(item)
			if err != nil {
				return nil, err
			}
		}
	}
	if err := stream.Error(); err != nil {
		return nil, err
	}

	if media == nil {
		return nil, nil
	}
	media.Tags = tags
	return media, nil
}

func isMediaMetadata(attributes map[string]types.AttributeValue) bool {
	if sk, ok := attributes["SK"].(*types.AttributeValueMemberS); ok {
		return sk.Value == "#METADATA"
	}
	return false
}
