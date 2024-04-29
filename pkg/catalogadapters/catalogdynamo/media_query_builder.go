package catalogdynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
)

const (
	albumIndex = "AlbumIndex"
)

func newMediaQueryBuilders(table string, request *catalog.FindMediaRequest, projectionName string) ([]*dynamodb.QueryInput, error) {
	var queries []*dynamodb.QueryInput

	builderWithProjection := func() expression.Builder {
		builder := expression.NewBuilder()
		if projectionName != "" {
			return builder.WithProjection(expression.NamesList(expression.Name(projectionName)))
		}

		return builder
	}

	for folderName, _ := range request.AlbumFolderNames {
		if len(request.Ranges) == 0 {
			expr, err := builderWithProjection().WithKeyCondition(expression.KeyAnd(
				withinAlbum(request.Owner, folderName),
				withExcludingMetaRecord(),
			)).Build()
			if err != nil {
				return nil, err
			}

			queries = append(queries, &dynamodb.QueryInput{
				TableName:                 &table,
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				IndexName:                 aws.String(albumIndex),
				KeyConditionExpression:    expr.KeyCondition(),
				ProjectionExpression:      expr.Projection(),
			})
		}

		for _, timeRange := range request.Ranges {
			expr, err := builderWithProjection().WithKeyCondition(expression.KeyAnd(
				withinAlbum(request.Owner, folderName),
				withinRange(timeRange),
			)).Build()
			if err != nil {
				return nil, err
			}

			queries = append(queries, &dynamodb.QueryInput{
				TableName:                 &table,
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				IndexName:                 aws.String(albumIndex),
				KeyConditionExpression:    expr.KeyCondition(),
				ProjectionExpression:      expr.Projection(),
			})
		}
	}
	return queries, nil
}

func withinAlbum(owner catalog.Owner, folderName catalog.FolderName) expression.KeyConditionBuilder {
	return expression.Key("AlbumIndexPK").Equal(expression.Value(AlbumIndexedKey(owner, folderName).AlbumIndexPK))
}

func withExcludingMetaRecord() expression.KeyConditionBuilder {
	return expression.Key("AlbumIndexSK").GreaterThanEqual(expression.Value("$"))
}

func withinRange(timeRange catalog.TimeRange) expression.KeyConditionBuilder {
	return expression.Key("AlbumIndexSK").Between(
		expression.Value(fmt.Sprintf("MEDIA#%s#", timeRange.Start.Format(IsoTime))),
		expression.Value(fmt.Sprintf("MEDIA#%s#", timeRange.End.Format(IsoTime))), // exclusive
	)
}
