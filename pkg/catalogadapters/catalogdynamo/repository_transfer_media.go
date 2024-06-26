package catalogdynamo

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutils"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

func (r *Repository) TransferMediasFromRecords(ctx context.Context, records catalog.MediaTransferRecords) (catalog.TransferredMedias, error) {
	medias := catalog.NewTransferredMedias()

	for albumId, selectors := range records {
		mediaIds, err := r.findMediaIdsFromSelectors(ctx, albumId, selectors)
		if err != nil {
			return medias, err
		}

		if len(mediaIds) > 0 {
			err = r.transferMedias(ctx, albumId, mediaIds)
			if err != nil {
				return medias, err
			}

			medias.Transfers[albumId] = mediaIds
		}
	}

	return medias, nil
}

func (r *Repository) findMediaIdsFromSelectors(ctx context.Context, targetAlbumId catalog.AlbumId, selectors []catalog.MediaSelector) ([]catalog.MediaId, error) {
	request := r.convertSelectorsIntoMediaRequest(targetAlbumId.Owner, selectors)

	queries, err := newMediaQueryBuilders(r.table, request, "Id", types.SelectAllAttributes)
	if err != nil {
		return nil, err
	}

	var mediaIds []catalog.MediaId

	crawler := dynamoutils.NewQueryStream(ctx, r.client, queries)
	for crawler.HasNext() {
		record := crawler.Next()
		if mediaId := readMediaId(record); mediaId != "" {
			mediaIds = append(mediaIds, mediaId)
		}
	}

	return mediaIds, nil
}

func (r *Repository) convertSelectorsIntoMediaRequest(owner ownermodel.Owner, selectors []catalog.MediaSelector) *catalog.FindMediaRequest {
	request := &catalog.FindMediaRequest{
		Owner:            owner,
		AlbumFolderNames: make(map[catalog.FolderName]interface{}),
		Ranges:           nil,
	}

	for _, selector := range selectors {
		for _, album := range selector.FromAlbums {
			request.AlbumFolderNames[album.FolderName] = nil
		}
		request.Ranges = append(request.Ranges, catalog.TimeRange{
			Start: selector.Start,
			End:   selector.End,
		})
	}
	return request
}

func (r *Repository) transferMedias(ctx context.Context, albumId catalog.AlbumId, mediaIds []catalog.MediaId) error {
	for _, id := range mediaIds {
		mediaKey, err := attributevalue.MarshalMap(MediaPrimaryKey(albumId.Owner, id))
		if err != nil {
			return err
		}

		update, err := expression.NewBuilder().WithUpdate(expression.Set(expression.Name("AlbumIndexPK"), expression.Value(AlbumIndexedKeyPK(albumId.Owner, albumId.FolderName)))).Build()

		_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
			ExpressionAttributeValues: update.Values(),
			ExpressionAttributeNames:  update.Names(),
			Key:                       mediaKey,
			TableName:                 &r.table,
			UpdateExpression:          update.Update(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
