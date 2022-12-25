package catalogdynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/thomasduchatelle/dphoto/pkg/catalog"
	"strings"
	"time"
)

type mediaQueryBuilder struct {
	query        *dynamodb.QueryInput
	keyCondition []string
	values       map[string]interface{}
}

func newMediaQueryBuilder(table string) *mediaQueryBuilder {
	return &mediaQueryBuilder{
		query: &dynamodb.QueryInput{
			IndexName: aws.String("AlbumIndex"),
			TableName: &table,
		},
		values: make(map[string]interface{}),
	}
}

func newMediaQueryBuilders(table string, request *catalog.FindMediaRequest, projectionExpression string) ([]*dynamodb.QueryInput, error) {
	var queries []*dynamodb.QueryInput

	for folderName, _ := range request.AlbumFolderNames {
		if len(request.Ranges) == 0 {
			builder := newMediaQueryBuilder(table)
			builder.WithAlbum(request.Owner, folderName)
			builder.WithoutAlbumMeta()

			builder.WithProjection(projectionExpression)

			query, err := builder.Build()
			if err != nil {
				return nil, err
			}

			queries = append(queries, query)
		}

		for _, timeRange := range request.Ranges {
			builder := newMediaQueryBuilder(table)
			builder.WithAlbum(request.Owner, folderName)
			builder.WithinTimeRange(timeRange.Start, timeRange.End)

			builder.WithProjection(projectionExpression)

			query, err := builder.Build()
			if err != nil {
				return nil, err
			}

			queries = append(queries, query)
		}
	}
	return queries, nil
}

func (b *mediaQueryBuilder) WithAlbum(owner, folderName string) {
	b.values[":albumKey"] = AlbumIndexedKey(owner, folderName).AlbumIndexPK
	b.keyCondition = append(b.keyCondition, "AlbumIndexPK = :albumKey")
}

func (b *mediaQueryBuilder) WithinTimeRange(start, end time.Time) {
	b.values[":from"] = fmt.Sprintf("MEDIA#%s#", start.Format(IsoTime))
	b.values[":to"] = fmt.Sprintf("MEDIA#%s#", end.Format(IsoTime)) // exclusive
	b.keyCondition = append(b.keyCondition, "AlbumIndexSK BETWEEN :from AND :to")
}

func (b *mediaQueryBuilder) WithoutAlbumMeta() {
	b.values[":excludeMeta"] = "$"
	b.keyCondition = append(b.keyCondition, "AlbumIndexSK >= :excludeMeta")
}

func (b *mediaQueryBuilder) WithPagination(limit int64, nextPageToken string) (err error) {
	if limit <= 0 && nextPageToken == "" {
		return
	}

	b.query.Limit = &limit
	b.query.ExclusiveStartKey, err = unmarshalPageToken(nextPageToken)

	return
}

func (b *mediaQueryBuilder) WithProjection(projectionExpression string) {
	if projectionExpression != "" {
		b.query.ProjectionExpression = &projectionExpression
		b.query.Select = aws.String(dynamodb.SelectSpecificAttributes)
	}
}

func (b *mediaQueryBuilder) Build() (*dynamodb.QueryInput, error) {
	b.query.KeyConditionExpression = aws.String(strings.Join(b.keyCondition, " AND "))

	values, err := dynamodbattribute.MarshalMap(b.values)
	b.query.ExpressionAttributeValues = values

	return b.query, err
}
