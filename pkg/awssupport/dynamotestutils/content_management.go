package dynamotestutils

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	dynamoutils "github.com/thomasduchatelle/dphoto/pkg/awssupport/dynamoutilsv2"
	"sort"
)

func (d *DynamodbTestContext) WithDbContent(ctx context.Context, entries []map[string]types.AttributeValue) error {
	err := d.clearContent(ctx)
	if err != nil {
		return err
	}

	if len(entries) > 0 {
		err = d.insertItems(ctx, entries, err)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DynamodbTestContext) insertItems(ctx context.Context, entries []map[string]types.AttributeValue, err error) error {
	requests := make([]types.WriteRequest, len(entries), len(entries))
	for i, item := range entries {
		requests[i] = types.WriteRequest{
			PutRequest: &types.PutRequest{Item: item},
		}
	}

	err = dynamoutils.BufferedWriteItems(ctx, d.Client, requests, d.Table, dynamoutils.DynamoWriteBatchSize)
	return err
}

func (d *DynamodbTestContext) clearContent(ctx context.Context) error {
	var deletionRequests []types.WriteRequest

	stream := dynamoutils.NewScanStream(ctx, d.Client, d.Table)
	for stream.HasNext() {
		entry := stream.Next()
		key := make(map[string]types.AttributeValue)
		key["PK"], _ = entry["PK"]
		key["SK"], _ = entry["SK"]
		deletionRequests = append(deletionRequests, types.WriteRequest{DeleteRequest: &types.DeleteRequest{Key: key}})
	}

	if stream.Error() != nil {
		return stream.Error()
	}

	return dynamoutils.BufferedWriteItems(ctx, d.Client, deletionRequests, d.Table, dynamoutils.DynamoWriteBatchSize)
}

func (d *DynamodbTestContext) Got(ctx context.Context) ([]map[string]types.AttributeValue, error) {
	items, err := dynamoutils.AsSlice(dynamoutils.NewScanStream(ctx, d.Client, d.Table))
	if err != nil {
		return nil, err
	}

	sort.Slice(items, func(i, j int) bool {
		var iPK, jPK, iSK, jSK string

		erriPK := attributevalue.Unmarshal(items[i]["PK"], &iPK)
		errjPK := attributevalue.Unmarshal(items[j]["PK"], &jPK)
		erriSK := attributevalue.Unmarshal(items[i]["SK"], &iSK)
		errjSK := attributevalue.Unmarshal(items[j]["SK"], &jSK)
		if erriPK != nil || errjPK != nil || erriSK != nil || errjSK != nil {
			assert.FailNowf(d.T, "sorting-content", "Unmarshal silenty failed: erriPK = %s, errjPK = %s, erriSK = %s, errjSK = %s", erriPK, errjPK, erriSK, errjSK)
			return false
		}

		if iPK == jPK {
			return iSK < jSK
		}
		return iPK < jPK
	})

	return items, err
}

func (d *DynamodbTestContext) EqualContent(ctx context.Context, wantItems []map[string]types.AttributeValue) (bool, error) {
	gotItems, err := d.Got(ctx)
	if err != nil {
		return false, err
	}

	got := make([]map[string]interface{}, 0)
	want := make([]map[string]interface{}, 0)

	err = attributevalue.UnmarshalListOfMaps(gotItems, &got)
	if err != nil {
		return false, err
	}

	err = attributevalue.UnmarshalListOfMaps(wantItems, &want)
	if err != nil {
		return false, err
	}

	return assert.Equal(d.T, want, got), nil

}
