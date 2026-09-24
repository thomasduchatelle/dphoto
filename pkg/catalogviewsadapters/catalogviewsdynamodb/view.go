package catalogviewsdynamodb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/catalogviews"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
	"github.com/thomasduchatelle/dphoto/pkg/usermodel"
	"slices"
	"time"
)

type AlbumViewRepository struct {
	Client    *dynamodb.Client
	TableName string
}

func (a *AlbumViewRepository) IncrementCountForAllViewers(ctx context.Context, updates []catalogviews.AlbumCountDiff) error {
	for _, update := range updates {
		items, err := a.queryByAlbumIndex(ctx, update.AlbumId)
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
		items, err := a.queryByAlbumIndex(ctx, update.AlbumId)
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

func (a *AlbumViewRepository) SetDisplayFieldsForAllViewers(ctx context.Context, albumId catalog.AlbumId, name string, start, end time.Time) error {
	items, err := a.queryByAlbumIndex(ctx, albumId)
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
	items, err := a.queryByAlbumIndex(ctx, existingId)
	if err != nil {
		return errors.Wrapf(err, "failed to list rows for album %v", existingId)
	}
	if len(items) == 0 {
		return nil
	}

	var ownerRow *catalogviews.UserAlbumSummary
	viewers := make([]catalogviews.Availability, 0, len(items))
	for _, item := range items {
		summary, err := unmarshalAlbumSummary(item)
		if err != nil {
			return err
		}
		viewers = append(viewers, summary.Availability)
		if summary.Availability.AsOwner {
			ownerRow = summary
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
	for _, item := range items {
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

	_, err := a.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &a.TableName,
		Key:       key.ToAttributes(),
	})

	return errors.Wrapf(err, "failed to delete row for album %v and user %v", albumId, availability)
}

func (a *AlbumViewRepository) DeleteAllRowsForAlbum(ctx context.Context, albumId catalog.AlbumId) error {
	pages, err := a.queryByAlbumIndex(ctx, albumId)
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

func (a *AlbumViewRepository) queryByAlbumIndex(ctx context.Context, albumId catalog.AlbumId) ([]map[string]types.AttributeValue, error) {
	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.Key("AlbumViewIndexPK").Equal(expression.Value(albumViewByAlbumIndexPK(albumId)))).
		Build()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build expression for AlbumViewByAlbumIndex query %v", albumId)
	}

	indexName := "AlbumViewByAlbumIndex"
	paginator := dynamodb.NewQueryPaginator(a.Client, &dynamodb.QueryInput{
		TableName:                 &a.TableName,
		IndexName:                 &indexName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})

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
		Build()

	if err != nil {
		return nil, errors.Wrapf(err, "failed to build expression for user %v", userId)
	}

	paginator := dynamodb.NewQueryPaginator(a.Client, &dynamodb.QueryInput{
		TableName:                 &a.TableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})

	var summaries []catalogviews.UserAlbumSummary
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.Items {
			summary, err := unmarshalAlbumSummary(item)
			if err != nil {
				return nil, err
			}

			summaries = append(summaries, *summary)
		}
	}

	return summaries, nil
}

func (a *AlbumViewRepository) ListSummariesForUserAndOwners(ctx context.Context, userId usermodel.UserId, owner ...ownermodel.Owner) ([]catalogviews.UserAlbumSummary, error) {
	summaries, err := a.ListSummariesForUser(ctx, userId)

	var filtered []catalogviews.UserAlbumSummary
	for _, summary := range summaries {
		if slices.Contains(owner, summary.AlbumSummary.AlbumId.Owner) && summary.Availability.UserId != "" {
			filtered = append(filtered, summary)
		}
	}

	return filtered, err
}
