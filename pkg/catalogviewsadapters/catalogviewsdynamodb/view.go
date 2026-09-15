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

// IncrementCounts applies `ADD Count :d` on each viewer row. It does NOT touch display fields:
// the update expression only mentions Count and the identity attributes needed to bootstrap a
// missing row (owner/folder/user/availability) so an existing row's AlbumName/AlbumStart/AlbumEnd
// survive untouched.
func (a *AlbumViewRepository) IncrementCounts(ctx context.Context, updates []catalogviews.AlbumMediaCountDiff) error {
	var inputs []*dynamodb.UpdateItemInput

	for _, update := range updates {
		for _, user := range update.Users {
			expr, err := expression.NewBuilder().
				WithUpdate(expression.
					Add(expression.Name("Count"), expression.Value(update.MediaCountDiff)).
					Set(expression.Name("UserId"), expression.Value(user.UserId.Value())).
					Set(expression.Name("AvailabilityType"), expression.Value(marshalAvailabilityType(user))).
					Set(expression.Name("AlbumOwner"), expression.Value(update.AlbumId.Owner)).
					Set(expression.Name("AlbumFolderName"), expression.Value(update.AlbumId.FolderName.String())),
				).
				Build()
			if err != nil {
				return errors.Wrapf(err, "failed to build expression for AlbumMediaCountDiff %+v", update)
			}

			inputs = append(inputs, &dynamodb.UpdateItemInput{
				TableName:                 &a.TableName,
				Key:                       albumSummaryKey(user, update.AlbumId).ToAttributes(),
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				UpdateExpression:          expr.Update(),
			})
		}
	}

	for _, input := range inputs {
		_, err := a.Client.UpdateItem(ctx, input)
		if err != nil {
			return errors.Wrapf(err, "failed to increment count for album %v", input.Key)
		}
	}

	return nil
}

// SetCounts applies `SET Count = :c` on each viewer row. It does NOT touch display fields.
func (a *AlbumViewRepository) SetCounts(ctx context.Context, updates []catalogviews.AlbumMediaCountForUsers) error {
	var inputs []*dynamodb.UpdateItemInput

	for _, update := range updates {
		for _, user := range update.Users {
			expr, err := expression.NewBuilder().
				WithUpdate(expression.
					Set(expression.Name("Count"), expression.Value(update.MediaCount)).
					Set(expression.Name("UserId"), expression.Value(user.UserId.Value())).
					Set(expression.Name("AvailabilityType"), expression.Value(marshalAvailabilityType(user))).
					Set(expression.Name("AlbumOwner"), expression.Value(update.AlbumId.Owner)).
					Set(expression.Name("AlbumFolderName"), expression.Value(update.AlbumId.FolderName.String())),
				).
				Build()
			if err != nil {
				return errors.Wrapf(err, "failed to build expression for AlbumMediaCountForUsers %+v", update)
			}

			inputs = append(inputs, &dynamodb.UpdateItemInput{
				TableName:                 &a.TableName,
				Key:                       albumSummaryKey(user, update.AlbumId).ToAttributes(),
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				UpdateExpression:          expr.Update(),
			})
		}
	}

	for _, input := range inputs {
		_, err := a.Client.UpdateItem(ctx, input)
		if err != nil {
			return errors.Wrapf(err, "failed to set count for album %v", input.Key)
		}
	}

	return nil
}

// SetDisplayFields updates AlbumName/AlbumStart/AlbumEnd on the viewer rows. It does NOT touch
// Count.
func (a *AlbumViewRepository) SetDisplayFields(ctx context.Context, albumId catalog.AlbumId, users []catalogviews.Availability, name string, start, end time.Time) error {
	var inputs []*dynamodb.UpdateItemInput

	for _, user := range users {
		expr, err := expression.NewBuilder().
			WithUpdate(expression.
				Set(expression.Name("AlbumName"), expression.Value(name)).
				Set(expression.Name("AlbumStart"), expression.Value(marshalTime(start))).
				Set(expression.Name("AlbumEnd"), expression.Value(marshalTime(end))).
				Set(expression.Name("UserId"), expression.Value(user.UserId.Value())).
				Set(expression.Name("AvailabilityType"), expression.Value(marshalAvailabilityType(user))).
				Set(expression.Name("AlbumOwner"), expression.Value(albumId.Owner)).
				Set(expression.Name("AlbumFolderName"), expression.Value(albumId.FolderName.String())),
			).
			Build()
		if err != nil {
			return errors.Wrapf(err, "failed to build expression for SetDisplayFields %+v/%+v", albumId, user)
		}

		inputs = append(inputs, &dynamodb.UpdateItemInput{
			TableName:                 &a.TableName,
			Key:                       albumSummaryKey(user, albumId).ToAttributes(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
		})
	}

	for _, input := range inputs {
		_, err := a.Client.UpdateItem(ctx, input)
		if err != nil {
			return errors.Wrapf(err, "failed to update display fields for album %v", input.Key)
		}
	}

	return nil
}

// PutSummaries is a full-row upsert (Put). It writes all attributes including Count.
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

// DeleteAllRowsForAlbum scans the projection to find every viewer row for the album and deletes
// each one. Cost is bounded by the number of viewers of the album.
func (a *AlbumViewRepository) DeleteAllRowsForAlbum(ctx context.Context, albumId catalog.AlbumId) error {
	expr, err := expression.NewBuilder().
		WithFilter(expression.
			Name("AlbumOwner").Equal(expression.Value(albumId.Owner.Value())).
			And(expression.Name("AlbumFolderName").Equal(expression.Value(albumId.FolderName.String())))).
		Build()
	if err != nil {
		return errors.Wrapf(err, "failed to build expression for DeleteAllRowsForAlbum %v", albumId)
	}

	paginator := dynamodb.NewScanPaginator(a.Client, &dynamodb.ScanInput{
		TableName:                 &a.TableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
	})

	var deletes []types.WriteRequest
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return errors.Wrapf(err, "failed to scan rows for album %v", albumId)
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
		// legacy entries cleanup: 'summary.Availability.UserId' wasn't set before 2024-06-30 and will be ignored during a drift reconciliation
		if slices.Contains(owner, summary.AlbumSummary.AlbumId.Owner) && summary.Availability.UserId != "" {
			filtered = append(filtered, summary)
		}
	}

	return filtered, err
}
