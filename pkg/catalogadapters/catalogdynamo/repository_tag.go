package catalogdynamo

import (
	"context"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

const mediaByTagIndex = "MediaByTagIndex"

// AddTag attaches the tag to the owner's media. The media's #METADATA Tags String Set is the
// source of truth: the tag is added with an atomic, idempotent ADD (concurrent writers merge
// instead of clobbering), which also bumps the media's TagVersion. The returned version stamps
// the derived TAG# sibling (feeding the medias-by-tag GSI), so a stale concurrent RemoveTag of
// the same tag cannot delete a sibling that a later AddTag owns. It is idempotent and
// owner-only (MediaNotFoundError on a media the owner does not have).
//
// The version-guarded, non-transactional reconciliation and its accepted crash residual are
// recorded in docs/adr/0006-media-tag-writes-version-guarded-not-transactional.md.
func (r *Repository) AddTag(ctx context.Context, owner ownermodel.Owner, mediaId catalog.MediaId, tag catalog.Tag) error {
	version, err := r.updateMediaTags(ctx, owner, mediaId, expression.Add(expression.Name("Tags"), expression.Value(stringSet(tag))))
	if err != nil {
		return err
	}
	return r.putTagSibling(ctx, owner, mediaId, tag, version)
}

// RemoveTag detaches the tag from the owner's media: it is removed from the #METADATA Tags
// String Set (the source of truth) with an atomic, idempotent DELETE that bumps the media's
// TagVersion, then the derived TAG# sibling is deleted — but only if no later AddTag has
// claimed it with a higher version. It is idempotent (removing an absent tag, or a tag of a
// media that does not exist, succeeds).
func (r *Repository) RemoveTag(ctx context.Context, owner ownermodel.Owner, mediaId catalog.MediaId, tag catalog.Tag) error {
	version, err := r.updateMediaTags(ctx, owner, mediaId, expression.Delete(expression.Name("Tags"), expression.Value(stringSet(tag))))
	if err != nil {
		if errors.Is(err, catalog.MediaNotFoundError) {
			return nil
		}
		return err
	}
	return r.deleteTagSibling(ctx, owner, mediaId, tag, version)
}

// updateMediaTags applies a String Set update (ADD or DELETE) to the media's #METADATA row and
// atomically bumps its TagVersion, returning the new version. Because DynamoDB serialises the
// #METADATA updates, the version is monotonic: a set mutation that logically happens later
// always carries a higher version. The update is guarded by attribute_exists(PK) so an absent
// (not owned) media is rejected as MediaNotFoundError rather than created.
func (r *Repository) updateMediaTags(ctx context.Context, owner ownermodel.Owner, mediaId catalog.MediaId, tagsUpdate expression.UpdateBuilder) (int64, error) {
	key, err := attributevalue.MarshalMap(MediaPrimaryKey(owner, mediaId))
	if err != nil {
		return 0, errors.Wrapf(err, "failed to marshal media key %s/%s", owner, mediaId)
	}

	update := tagsUpdate.Add(expression.Name("TagVersion"), expression.Value(1))
	expr, err := expression.NewBuilder().
		WithUpdate(update).
		WithCondition(expression.Name("PK").AttributeExists()).
		Build()
	if err != nil {
		return 0, err
	}

	output, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 &r.table,
		Key:                       key,
		UpdateExpression:          expr.Update(),
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	})
	var conditionalCheckFailedException *types.ConditionalCheckFailedException
	if errors.As(err, &conditionalCheckFailedException) {
		return 0, errors.Wrapf(catalog.MediaNotFoundError, "cannot update tags of media %s/%s", owner, mediaId)
	}
	if err != nil {
		return 0, errors.Wrapf(err, "failed to persist tags of media %s/%s", owner, mediaId)
	}

	version, ok := output.Attributes["TagVersion"].(*types.AttributeValueMemberN)
	if !ok {
		return 0, errors.Errorf("UpdateItem did not return the new TagVersion of media %s/%s", owner, mediaId)
	}
	return strconv.ParseInt(version.Value, 10, 64)
}

// putTagSibling writes the derived TAG# index item stamped with the given version. It only
// overwrites a sibling stamped with an older version, so a stale AddTag cannot resurrect a
// sibling that a later RemoveTag deleted. Losing the version race is a no-op success.
func (r *Repository) putTagSibling(ctx context.Context, owner ownermodel.Owner, mediaId catalog.MediaId, tag catalog.Tag, version int64) error {
	item, err := marshalTag(owner, mediaId, tag, version)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal tag %s for media %s/%s", tag, owner, mediaId)
	}

	cond, err := expression.NewBuilder().
		WithCondition(expression.Or(
			expression.Name("TagVersion").AttributeNotExists(),
			expression.Name("TagVersion").LessThan(expression.Value(version)),
		)).
		Build()
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:                      item,
		TableName:                 &r.table,
		ConditionExpression:       cond.Condition(),
		ExpressionAttributeNames:  cond.Names(),
		ExpressionAttributeValues: cond.Values(),
	})
	var conditionalCheckFailedException *types.ConditionalCheckFailedException
	if errors.As(err, &conditionalCheckFailedException) {
		return nil
	}
	return errors.Wrapf(err, "failed to add tag %s to media %s/%s", tag, owner, mediaId)
}

// deleteTagSibling removes the derived TAG# index item, but only if no later AddTag has claimed
// it with a version higher than the one of the RemoveTag reconciling it. Losing the version
// race (a newer sibling exists) is a no-op success, and deleting an absent sibling succeeds too.
func (r *Repository) deleteTagSibling(ctx context.Context, owner ownermodel.Owner, mediaId catalog.MediaId, tag catalog.Tag, version int64) error {
	key, err := attributevalue.MarshalMap(TagPrimaryKey(owner, mediaId, tag))
	if err != nil {
		return errors.Wrapf(err, "failed to marshal tag key %s for media %s/%s", tag, owner, mediaId)
	}

	cond, err := expression.NewBuilder().
		WithCondition(expression.Or(
			expression.Name("TagVersion").AttributeNotExists(),
			expression.Name("TagVersion").LessThan(expression.Value(version)),
		)).
		Build()
	if err != nil {
		return err
	}

	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		Key:                       key,
		TableName:                 &r.table,
		ConditionExpression:       cond.Condition(),
		ExpressionAttributeNames:  cond.Names(),
		ExpressionAttributeValues: cond.Values(),
	})
	var conditionalCheckFailedException *types.ConditionalCheckFailedException
	if errors.As(err, &conditionalCheckFailedException) {
		return nil
	}
	return errors.Wrapf(err, "failed to remove tag %s from media %s/%s", tag, owner, mediaId)
}

func stringSet(tag catalog.Tag) *types.AttributeValueMemberSS {
	return &types.AttributeValueMemberSS{Value: []string{tag.String()}}
}

// FindMediasByTag queries the MediaByTagIndex to list the medias carrying the tag in
// O(matches), then batch-reads their #METADATA rows to attach each media's metadata, its full
// tag set, and the album it belongs to, so each returned media is ready for the Tags page.
//
// The version-guarded sibling reconciliation (see AddTag/RemoveTag) keeps the index consistent
// with the #METADATA Tags source of truth under concurrency, but a failure between the two
// writes can still leave an orphan sibling. Such a GSI hit whose #METADATA no longer carries the
// tag is dropped here (the truth is what the customer sees) and logged as a WARN to surface it.
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

	var keys []map[string]types.AttributeValue
	for stream.HasNext() {
		mediaId, err := unmarshalTagMediaId(stream.Next())
		if err != nil {
			return nil, err
		}
		key, err := attributevalue.MarshalMap(MediaPrimaryKey(owner, mediaId))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal media key %s/%s", owner, mediaId)
		}
		keys = append(keys, key)
	}
	if err := stream.Error(); err != nil {
		return nil, err
	}

	metadataStream := dynamoutils.NewGetStream(ctx, dynamoutils.NewGetBatchItem(r.client, r.table, ""), keys, dynamoutils.DynamoReadBatchSize)

	var medias []*catalog.MediaWithAlbum
	var drifted []catalog.MediaId
	for metadataStream.HasNext() {
		media, albumId, err := unmarshalMediaWithAlbum(metadataStream.Next())
		if err != nil {
			return nil, err
		}
		if !containsTag(media.Tags, tag) {
			drifted = append(drifted, media.Id)
			continue
		}
		medias = append(medias, &catalog.MediaWithAlbum{MediaMeta: *media, AlbumId: albumId})
	}
	if len(drifted) > 0 {
		log.WithField("Owner", owner).Warnf("medias-by-tag index drift for tag %s: %d media(s) returned by the GSI no longer carry the tag and were dropped: %v", tag, len(drifted), drifted)
	}

	return medias, metadataStream.Error()
}

func containsTag(tags []catalog.Tag, tag catalog.Tag) bool {
	for _, t := range tags {
		if t.String() == tag.String() {
			return true
		}
	}
	return false
}

// unmarshalMediaWithAlbum reads a media's #METADATA row, returning its metadata (tags included)
// and the album it belongs to, resolved from the AlbumIndexPK.
func unmarshalMediaWithAlbum(item map[string]types.AttributeValue) (*catalog.MediaMeta, catalog.AlbumId, error) {
	media, err := unmarshalMediaMetaData(item)
	if err != nil {
		return nil, catalog.AlbumId{}, err
	}

	albumId := catalog.AlbumId{}
	if value, ok := item["AlbumIndexPK"].(*types.AttributeValueMemberS); ok {
		if owner, folderName, found := strings.Cut(value.Value, "#"); found {
			albumId = catalog.AlbumId{Owner: ownermodel.Owner(owner), FolderName: catalog.NewFolderName(folderName)}
		}
	}
	return media, albumId, nil
}
