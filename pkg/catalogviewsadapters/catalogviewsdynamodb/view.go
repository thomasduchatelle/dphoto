package catalogviewsdynamodb

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/appdynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

const legacyCountSuffix = "#COUNT"

type AlbumViewRepository struct {
	Client    *dynamodb.Client
	TableName string
}

func legacyDeleteRequest(record map[string]types.AttributeValue) types.WriteRequest {
	pk := record["PK"].(*types.AttributeValueMemberS).Value
	sk := record["SK"].(*types.AttributeValueMemberS).Value
	legacyKey := appdynamodb.TablePk{PK: pk, SK: sk + legacyCountSuffix}
	return types.WriteRequest{
		DeleteRequest: &types.DeleteRequest{Key: legacyKey.ToAttributes()},
	}
}

func (a *AlbumViewRepository) IncrementCountForAllViewers(ctx context.Context, updates []catalogviews.AlbumCountDiff) error {
	for _, update := range updates {
		items, err := a.querySummariesByAlbumIndex(ctx, update.AlbumId)
		if err != nil {
			return errors.Wrapf(err, "failed to list rows for album %v", update.AlbumId)
		}

		for _, item := range items {
			expr, err := expression.NewBuilder().
				WithUpdate(expression.Add(expression.Name("Count"), expression.Value(update.MediaCountDiff))).
				Build()
			if err != nil {
				return errors.Wrapf(err, "failed to build expression for AlbumCountDiff %+v", update)
			}

			_, err = a.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
				TableName: &a.TableName,
				Key: map[string]types.AttributeValue{
					"PK": item["PK"],
					"SK": item["SK"],
				},
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				UpdateExpression:          expr.Update(),
			})
			if err != nil {
				return errors.Wrapf(err, "failed to increment count for album %v", update.AlbumId)
			}
		}
	}

	return nil
}

func (a *AlbumViewRepository) SetCountForAllViewers(ctx context.Context, updates []catalogviews.AlbumCount) error {
	for _, update := range updates {
		items, err := a.querySummariesByAlbumIndex(ctx, update.AlbumId)
		if err != nil {
			return errors.Wrapf(err, "failed to list rows for album %v", update.AlbumId)
		}

		for _, item := range items {
			expr, err := expression.NewBuilder().
				WithUpdate(expression.Set(expression.Name("Count"), expression.Value(update.MediaCount))).
				Build()
			if err != nil {
				return errors.Wrapf(err, "failed to build expression for AlbumCount %+v", update)
			}

			_, err = a.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
				TableName: &a.TableName,
				Key: map[string]types.AttributeValue{
					"PK": item["PK"],
					"SK": item["SK"],
				},
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				UpdateExpression:          expr.Update(),
			})
			if err != nil {
				return errors.Wrapf(err, "failed to set count for album %v", update.AlbumId)
			}
		}
	}

	return nil
}

// PutCoversForAllViewers upserts the covers row for every viewer who already has a
// summary row for the album. The covers row lives at a sibling SK (`...#COVERS`) in
// the viewer's partition so cover writes never touch the main summary row.
func (a *AlbumViewRepository) PutCoversForAllViewers(ctx context.Context, albumId catalog.AlbumId, covers []catalog.Cover) error {
	summaries, err := a.querySummariesByAlbumIndex(ctx, albumId)
	if err != nil {
		return errors.Wrapf(err, "failed to list rows for album %v", albumId)
	}
	if len(summaries) == 0 {
		return nil
	}

	var writes []types.WriteRequest
	for _, item := range summaries {
		availability, err := availabilityFromItem(item)
		if err != nil {
			return err
		}
		record, err := marshalAlbumCovers(availability, albumId, covers)
		if err != nil {
			return err
		}
		writes = append(writes, types.WriteRequest{PutRequest: &types.PutRequest{Item: record}})
	}

	return dynamoutils.BufferedWriteItems(ctx, a.Client, writes, a.TableName, dynamoutils.DynamoWriteBatchSize)
}

// DeleteCoversForAllViewers removes the covers row of every viewer for the album.
// Summary rows are untouched; a next ListSummariesForUser returns the album with an
// empty cover set.
func (a *AlbumViewRepository) DeleteCoversForAllViewers(ctx context.Context, albumId catalog.AlbumId) error {
	items, err := a.queryCoversByAlbumIndex(ctx, albumId)
	if err != nil {
		return errors.Wrapf(err, "failed to list cover rows for album %v", albumId)
	}
	if len(items) == 0 {
		return nil
	}

	writes := make([]types.WriteRequest, 0, len(items))
	for _, item := range items {
		writes = append(writes, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{Key: map[string]types.AttributeValue{
				"PK": item["PK"],
				"SK": item["SK"],
			}},
		})
	}
	return dynamoutils.BufferedWriteItems(ctx, a.Client, writes, a.TableName, dynamoutils.DynamoWriteBatchSize)
}

func (a *AlbumViewRepository) SetDisplayFieldsForAllViewers(ctx context.Context, albumId catalog.AlbumId, name string, start, end time.Time) error {
	items, err := a.querySummariesByAlbumIndex(ctx, albumId)
	if err != nil {
		return errors.Wrapf(err, "failed to list rows for album %v", albumId)
	}

	for _, item := range items {
		expr, err := expression.NewBuilder().
			WithUpdate(expression.
				Set(expression.Name("AlbumName"), expression.Value(name)).
				Set(expression.Name("AlbumStart"), expression.Value(marshalTime(start))).
				Set(expression.Name("AlbumEnd"), expression.Value(marshalTime(end))),
			).
			Build()
		if err != nil {
			return errors.Wrapf(err, "failed to build expression for SetDisplayFieldsForAllViewers %+v", albumId)
		}

		_, err = a.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			TableName: &a.TableName,
			Key: map[string]types.AttributeValue{
				"PK": item["PK"],
				"SK": item["SK"],
			},
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
		})
		if err != nil {
			return errors.Wrapf(err, "failed to update display fields for album %v", albumId)
		}
	}

	return nil
}

func (a *AlbumViewRepository) RenameAlbum(ctx context.Context, existingId, renamedId catalog.AlbumId, newName string) error {
	allItems, err := a.queryAllByAlbumIndex(ctx, existingId)
	if err != nil {
		return errors.Wrapf(err, "failed to list rows for album %v", existingId)
	}
	if len(allItems) == 0 {
		return nil
	}

	var ownerRow *catalogviews.UserAlbumSummary
	viewers := make([]catalogviews.Availability, 0)
	coversPerViewer := make(map[string][]catalog.Cover)
	for _, item := range allItems {
		recordType := recordTypeOf(item)
		switch recordType {
		case RecordTypeSummary:
			summary, err := unmarshalAlbumSummary(item)
			if err != nil {
				return err
			}
			viewers = append(viewers, summary.Availability)
			if summary.Availability.AsOwner {
				ownerRow = summary
			}
		case RecordTypeCovers:
			_, availability, covers, err := unmarshalAlbumCovers(item)
			if err != nil {
				return err
			}
			coversPerViewer[availability.String()] = covers
		}
	}
	if ownerRow == nil {
		return errors.Errorf("no owner row for album %v: cannot RenameAlbum without a source of truth", existingId)
	}

	renamedSummary := catalogviews.AlbumSummary{
		AlbumId:    renamedId,
		Name:       newName,
		Start:      ownerRow.AlbumSummary.Start,
		End:        ownerRow.AlbumSummary.End,
		MediaCount: ownerRow.AlbumSummary.MediaCount,
	}

	var writes []types.WriteRequest
	for _, item := range allItems {
		writes = append(writes, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					"PK": item["PK"],
					"SK": item["SK"],
				},
			},
		})
	}

	newItems, err := marshalAlbumSummary(catalogviews.AlbumSummaryForUsers{
		AlbumSummary: renamedSummary,
		Users:        viewers,
	})
	if err != nil {
		return err
	}
	for _, item := range newItems {
		writes = append(writes, types.WriteRequest{
			PutRequest: &types.PutRequest{Item: item},
		})
	}

	for _, viewer := range viewers {
		covers, ok := coversPerViewer[viewer.String()]
		if !ok || len(covers) == 0 {
			continue
		}
		coverItem, err := marshalAlbumCovers(viewer, renamedId, covers)
		if err != nil {
			return err
		}
		writes = append(writes, types.WriteRequest{
			PutRequest: &types.PutRequest{Item: coverItem},
		})
	}

	return dynamoutils.BufferedWriteItems(ctx, a.Client, writes, a.TableName, dynamoutils.DynamoWriteBatchSize)
}

func (a *AlbumViewRepository) PutSummaries(ctx context.Context, summaries []catalogviews.AlbumSummaryForUsers) error {
	var items []types.WriteRequest

	for _, summary := range summaries {
		records, err := marshalAlbumSummary(summary)
		if err != nil {
			return err
		}

		for _, record := range records {
			items = append(items, types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: record,
				},
			})
			items = append(items, legacyDeleteRequest(record))
		}
	}

	if len(items) > 0 {
		err := dynamoutils.BufferedWriteItems(ctx, a.Client, items, a.TableName, dynamoutils.DynamoWriteBatchSize)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AlbumViewRepository) DeleteRow(ctx context.Context, availability catalogviews.Availability, albumId catalog.AlbumId) error {
	key := albumSummaryKey(availability, albumId)
	coversKey := albumCoversKey(availability, albumId)
	legacyKey := appdynamodb.TablePk{PK: key.PK, SK: key.SK + legacyCountSuffix}

	writes := []types.WriteRequest{
		{DeleteRequest: &types.DeleteRequest{Key: key.ToAttributes()}},
		{DeleteRequest: &types.DeleteRequest{Key: coversKey.ToAttributes()}},
		{DeleteRequest: &types.DeleteRequest{Key: legacyKey.ToAttributes()}},
	}

	err := dynamoutils.BufferedWriteItems(ctx, a.Client, writes, a.TableName, dynamoutils.DynamoWriteBatchSize)
	return errors.Wrapf(err, "failed to delete row for album %v and user %v", albumId, availability)
}

func (a *AlbumViewRepository) DeleteAllRowsForAlbum(ctx context.Context, albumId catalog.AlbumId) error {
	pages, err := a.queryAllByAlbumIndex(ctx, albumId)
	if err != nil {
		return errors.Wrapf(err, "failed to list rows for album %v", albumId)
	}

	var deletes []types.WriteRequest
	for _, item := range pages {
		deletes = append(deletes, types.WriteRequest{
			DeleteRequest: &types.DeleteRequest{
				Key: map[string]types.AttributeValue{
					"PK": item["PK"],
					"SK": item["SK"],
				},
			},
		})
	}

	if len(deletes) == 0 {
		return nil
	}
	return dynamoutils.BufferedWriteItems(ctx, a.Client, deletes, a.TableName, dynamoutils.DynamoWriteBatchSize)
}

func (a *AlbumViewRepository) querySummariesByAlbumIndex(ctx context.Context, albumId catalog.AlbumId) ([]map[string]types.AttributeValue, error) {
	return a.queryByAlbumIndex(ctx, albumId, &RecordTypeSummaryFilter)
}

func (a *AlbumViewRepository) queryCoversByAlbumIndex(ctx context.Context, albumId catalog.AlbumId) ([]map[string]types.AttributeValue, error) {
	return a.queryByAlbumIndex(ctx, albumId, &RecordTypeCoversFilter)
}

func (a *AlbumViewRepository) queryAllByAlbumIndex(ctx context.Context, albumId catalog.AlbumId) ([]map[string]types.AttributeValue, error) {
	return a.queryByAlbumIndex(ctx, albumId, nil)
}

var (
	RecordTypeSummaryFilter = RecordTypeSummary
	RecordTypeCoversFilter  = RecordTypeCovers
)

func (a *AlbumViewRepository) queryByAlbumIndex(ctx context.Context, albumId catalog.AlbumId, filterRecordType *string) ([]map[string]types.AttributeValue, error) {
	builder := expression.NewBuilder().
		WithKeyCondition(expression.Key("AlbumViewIndexPK").Equal(expression.Value(albumViewByAlbumIndexPK(albumId))))

	if filterRecordType != nil {
		if *filterRecordType == RecordTypeSummary {
			builder = builder.WithFilter(
				expression.Name("RecordType").Equal(expression.Value(RecordTypeSummary)).
					Or(expression.Name("RecordType").AttributeNotExists()),
			)
		} else {
			builder = builder.WithFilter(expression.Name("RecordType").Equal(expression.Value(*filterRecordType)))
		}
	}

	expr, err := builder.Build()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build expression for AlbumViewByAlbumIndex query %v", albumId)
	}

	indexName := "AlbumViewByAlbumIndex"
	queryInput := &dynamodb.QueryInput{
		TableName:                 &a.TableName,
		IndexName:                 &indexName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}
	if expr.Filter() != nil {
		queryInput.FilterExpression = expr.Filter()
	}

	paginator := dynamodb.NewQueryPaginator(a.Client, queryInput)

	var items []map[string]types.AttributeValue
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		items = append(items, page.Items...)
	}
	return items, nil
}

func (a *AlbumViewRepository) ListSummariesForUser(ctx context.Context, userId usermodel.UserId) ([]catalogviews.UserAlbumSummary, error) {
	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.Key("PK").Equal(expression.Value(albumsViewPK(userId)))).
		WithFilter(expression.Name("AvailabilityType").AttributeExists()).
		Build()

	if err != nil {
		return nil, errors.Wrapf(err, "failed to build expression for user %v", userId)
	}

	paginator := dynamodb.NewQueryPaginator(a.Client, &dynamodb.QueryInput{
		TableName:                 &a.TableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
	})

	type albumKey struct {
		albumId      catalog.AlbumId
		availability catalogviews.Availability
	}
	summariesByKey := make(map[string]*catalogviews.UserAlbumSummary)
	coversByKey := make(map[string][]catalog.Cover)
	var order []string
	toKey := func(k albumKey) string {
		return k.availability.String() + "|" + k.albumId.String()
	}

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.Items {
			switch recordTypeOf(item) {
			case RecordTypeCovers:
				albumId, availability, covers, err := unmarshalAlbumCovers(item)
				if err != nil {
					return nil, err
				}
				key := toKey(albumKey{albumId: albumId, availability: availability})
				coversByKey[key] = covers
			default:
				summary, err := unmarshalAlbumSummary(item)
				if err != nil {
					return nil, err
				}
				key := toKey(albumKey{albumId: summary.AlbumSummary.AlbumId, availability: summary.Availability})
				if _, exists := summariesByKey[key]; !exists {
					order = append(order, key)
				}
				summariesByKey[key] = summary
			}
		}
	}

	if len(order) == 0 {
		return nil, nil
	}
	summaries := make([]catalogviews.UserAlbumSummary, 0, len(order))
	for _, key := range order {
		summary := summariesByKey[key]
		if covers, ok := coversByKey[key]; ok {
			summary.AlbumSummary.Covers = covers
		}
		summaries = append(summaries, *summary)
	}

	return summaries, nil
}

func (a *AlbumViewRepository) DeleteLegacyRowsForUser(ctx context.Context, userId usermodel.UserId) error {
	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.Key("PK").Equal(expression.Value(albumsViewPK(userId)))).
		WithFilter(expression.Name("AvailabilityType").AttributeNotExists()).
		Build()
	if err != nil {
		return errors.Wrapf(err, "failed to build expression for user %v", userId)
	}

	paginator := dynamodb.NewQueryPaginator(a.Client, &dynamodb.QueryInput{
		TableName:                 &a.TableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      aws.String("PK, SK"),
	})

	var deletes []types.WriteRequest
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to list legacy rows for user %v", userId)
		}
		for _, item := range page.Items {
			deletes = append(deletes, types.WriteRequest{
				DeleteRequest: &types.DeleteRequest{
					Key: map[string]types.AttributeValue{
						"PK": item["PK"],
						"SK": item["SK"],
					},
				},
			})
		}
	}

	if len(deletes) == 0 {
		return nil
	}
	return dynamoutils.BufferedWriteItems(ctx, a.Client, deletes, a.TableName, dynamoutils.DynamoWriteBatchSize)
}

func (a *AlbumViewRepository) ListSummariesForUserAndOwners(ctx context.Context, userId usermodel.UserId, owner ...ownermodel.Owner) ([]catalogviews.UserAlbumSummary, error) {
	summaries, err := a.ListSummariesForUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	var filtered []catalogviews.UserAlbumSummary
	for _, summary := range summaries {
		if slices.Contains(owner, summary.AlbumSummary.AlbumId.Owner) {
			filtered = append(filtered, summary)
		}
	}

	return filtered, nil
}

func recordTypeOf(item map[string]types.AttributeValue) string {
	if attr, ok := item["RecordType"]; ok {
		if s, ok := attr.(*types.AttributeValueMemberS); ok {
			return s.Value
		}
	}
	sk, ok := item["SK"].(*types.AttributeValueMemberS)
	if ok && strings.HasSuffix(sk.Value, coversSKSuffix) {
		return RecordTypeCovers
	}
	return RecordTypeSummary
}

func availabilityFromItem(item map[string]types.AttributeValue) (catalogviews.Availability, error) {
	availability, ok := item["AvailabilityType"].(*types.AttributeValueMemberS)
	if !ok {
		return catalogviews.Availability{}, errors.Errorf("missing AvailabilityType on row: %+v", item)
	}
	userId, ok := item["UserId"].(*types.AttributeValueMemberS)
	if !ok {
		return catalogviews.Availability{}, errors.Errorf("missing UserId on row: %+v", item)
	}
	return unmarshalAvailability(availability.Value, userId.Value, item)
}
