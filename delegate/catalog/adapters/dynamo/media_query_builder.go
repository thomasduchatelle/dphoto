package dynamo

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
func (b *mediaQueryBuilder) WithAlbum(owner, folderName string) {
	b.values[":albumKey"] = albumIndexedKey(owner, folderName).AlbumIndexPK
	b.keyCondition = append(b.keyCondition, "AlbumIndexPK = :albumKey")
}

func (b *mediaQueryBuilder) Within(start, end time.Time) {
	b.values[":from"] = fmt.Sprintf("MEDIA#%s#", start.Format(IsoTime))
	b.values[":to"] = fmt.Sprintf("MEDIA#%s#", end.Format(IsoTime)) // exclusive
	b.keyCondition = append(b.keyCondition, "AlbumIndexSK BETWEEN :from AND :to")
}

func (b *mediaQueryBuilder) ExcludeAlbumMeta() {
	b.values[":excludeMeta"] = "$"
	b.keyCondition = append(b.keyCondition, "AlbumIndexSK >= :excludeMeta")
}

func (b *mediaQueryBuilder) AddPagination(limit int64, nextPageToken string) error {
	if limit <= 0 {
		limit = defaultPage
	}
	b.query.Limit = &limit

	startKey, err := unmarshalPageToken(nextPageToken)
	b.query.ExclusiveStartKey = startKey

	return err
}

func (b *mediaQueryBuilder) WithProjection(projectionExpression string) {
	b.query.ProjectionExpression = &projectionExpression
	b.query.Select = aws.String(dynamodb.SelectSpecificAttributes)
}

func (b *mediaQueryBuilder) Build() (*dynamodb.QueryInput, error) {
	b.query.KeyConditionExpression = aws.String(strings.Join(b.keyCondition, " AND "))

	values, err := dynamodbattribute.MarshalMap(b.values)
	b.query.ExpressionAttributeValues = values

	return b.query, err
}
