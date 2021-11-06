package dynamo

import (
	"github.com/thomasduchatelle/dphoto/delegate/catalog"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	versionTagName       = "Version"
	tableVersion         = "1.1" // to be bumped manually when schema is updated
	dynamoWriteBatchSize = 25
	dynamoReadBatchSize  = 100
	defaultPage          = 50
)

func (r *Rep) InsertMedias(medias []catalog.CreateMediaRequest) error {
	requests := make([]*dynamodb.WriteRequest, len(medias)*2)
	for index, media := range medias {
		mediaEntry, locationEntry, err := marshalMedia(r.RootOwner, &media)
		if err != nil {
			return errors.Wrapf(err, "Failed mapping media %s", fmt.Sprint(media))
		}

		requests[index*2] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: locationEntry,
			},
		}
		requests[index*2+1] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: mediaEntry,
			},
		}
	}

	return r.bufferedWriteItems(requests)
}

func (r *Rep) FindMedias(folderName string, filter catalog.FindMediaFilter) (*catalog.MediaPage, error) {
	builder := newMediaQueryBuilder(r.table)

	builder.WithAlbum(r.RootOwner, folderName)
	if !filter.TimeRange.Start.IsZero() && !filter.TimeRange.End.IsZero() {
		builder.Within(filter.TimeRange.Start, filter.TimeRange.End)
	} else {
		builder.ExcludeAlbumMeta()
	}
	err := builder.AddPagination(filter.PageRequest.Size, filter.PageRequest.NextPage)
	if err != nil {
		return nil, err
	}

	query, err := builder.Build()
	if err != nil {
		return nil, err
	}

	results, err := r.db.Query(query)
	if err != nil {
		return nil, errors.Wrapf(err, "query failed [%s]", query)
	}

	page := new(catalog.MediaPage)
	for _, media := range results.Items {
		data, err := unmarshalMediaMetaData(media)
		if err != nil {
			return nil, err
		}
		page.Content = append(page.Content, data)
	}

	page.NextPage, err = r.marshalPageToken(results.LastEvaluatedKey)
	return page, err
}

func (r *Rep) FindExistingSignatures(signatures []*catalog.MediaSignature) ([]*catalog.MediaSignature, error) {
	keys := make([]map[string]*dynamodb.AttributeValue, len(signatures))
	for index, signature := range signatures {
		key, err := dynamodbattribute.MarshalMap(mediaPrimaryKey(r.RootOwner, signature))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal media keys from signature %+v", signature)
		}

		keys[index] = key
	}

	stream := r.bufferedBatchGetItem(keys, aws.String("SignatureSize, SignatureHash"))

	found := make([]*catalog.MediaSignature, 0, len(signatures))
	for stream.HasNext() {
		attributes := stream.Next()

		var signature MediaData
		err := dynamodbattribute.UnmarshalMap(attributes, &signature)
		if err != nil {
			return nil, err
		}

		found = append(found, &catalog.MediaSignature{
			SignatureSha256: signature.SignatureHash,
			SignatureSize:   signature.SignatureSize,
		})
	}

	return found, stream.Error()
}

func (r *Rep) FindMediaLocationsSignatures(signatures []*catalog.MediaSignature) ([]*catalog.MediaSignatureAndLocation, error) {
	keys := make([]map[string]*dynamodb.AttributeValue, len(signatures))
	for index, signature := range signatures {
		key, err := dynamodbattribute.MarshalMap(mediaLocationPrimaryKey(r.RootOwner, signature))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to marshal media keys from signature %+v", signature)
		}

		keys[index] = key
	}

	stream := r.bufferedBatchGetItem(keys, aws.String("FolderName, Filename, SignatureHash, SignatureSize"))

	locations := make([]*catalog.MediaSignatureAndLocation, 0, len(signatures))
	for stream.HasNext() {
		attributes := stream.Next()

		var location MediaLocationData
		err := dynamodbattribute.UnmarshalMap(attributes, &location)
		if err != nil {
			return nil, err
		}

		locations = append(locations, &catalog.MediaSignatureAndLocation{
			Location: catalog.MediaLocation{
				FolderName: location.FolderName,
				Filename:   location.Filename,
			},
			Signature: catalog.MediaSignature{
				SignatureSha256: location.SignatureHash,
				SignatureSize:   location.SignatureSize,
			},
		})
	}

	return locations, stream.Error()
}

func (r *Rep) bufferedBatchGetItem(keys []map[string]*dynamodb.AttributeValue, projectionExpression *string) Stream {
	return NewGetStream(r, keys, projectionExpression, dynamoReadBatchSize)
}

func (r *Rep) bufferedQueriesCrawler(queries []*dynamodb.QueryInput) Stream {
	return NewQueryStream(r, queries)
}

func (r *Rep) bufferedWriteItems(requests []*dynamodb.WriteRequest) error {
	retry := backoff.NewExponentialBackOff()
	buffer := make([]*dynamodb.WriteRequest, 0, dynamoWriteBatchSize)

	for len(buffer) > 0 || len(requests) > 0 {
		end := cap(buffer) - len(buffer)
		if end > len(requests) {
			end = len(requests)
		}

		if end > 0 {
			buffer = append(buffer, requests[:end]...)
			requests = requests[end:]
		}

		err := backoff.RetryNotify(func() error {
			result, err := r.db.BatchWriteItem(&dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]*dynamodb.WriteRequest{
					r.table: buffer,
				},
			})

			if err != nil {
				return err
			}

			buffer = buffer[:0]
			if unprocessed, ok := result.UnprocessedItems[r.table]; ok && len(unprocessed) > 0 {
				buffer = append(buffer, unprocessed...)
			}

			return nil

		}, retry, func(err error, duration time.Duration) {
			log.WithFields(log.Fields{
				"Duration": duration,
			}).WithError(err).Warnf("Retrying inserting media (buffer len %d)", len(buffer))
		})
		if err != nil {
			return errors.Wrapf(err, "failed to insert batch %+v", buffer)
		}
	}

	return nil
}

func (r *Rep) bufferedUpdateItems(updates []*dynamodb.UpdateItemInput) error {
	for _, update := range updates {
		_, err := r.db.UpdateItem(update)
		if err != nil {
			return err
		}
	}

	return nil
}

// unmarshalPageToken take a base64 page id and decode it into a dynamodb key. Return nil, nil if nextPageToken is blank
func unmarshalPageToken(nextPageToken string) (map[string]*dynamodb.AttributeValue, error) {
	var startKey map[string]*dynamodb.AttributeValue
	if !isBlank(nextPageToken) {
		startKeyJson, err := base64.StdEncoding.DecodeString(nextPageToken)
		if err != nil {
			return startKey, errors.Wrapf(err, "Invalid nextPageToken %s", nextPageToken)
		}

		err = json.Unmarshal(startKeyJson, &startKey)
		if err != nil {
			return startKey, errors.Wrapf(err, "Invalid nextPageToken %s", startKeyJson)
		}
	}

	return startKey, nil
}

func (r *Rep) unmarshalPageToken(nextPageToken string) (map[string]*dynamodb.AttributeValue, error) {
	return unmarshalPageToken(nextPageToken)
}

// marshalPageToken take a dynamodb key and encode it (JSON+BASE64) to be used by service. Return empty string when key is empty
func (r *Rep) marshalPageToken(key map[string]*dynamodb.AttributeValue) (string, error) {
	if len(key) == 0 {
		return "", nil
	}

	keyJson, err := json.Marshal(key)
	return base64.StdEncoding.EncodeToString(keyJson), err
}

// panic if value can't be converted into an attribute
func mustAttribute(value interface{}) *dynamodb.AttributeValue {
	attribute, err := dynamodbattribute.Marshal(value)
	if err != nil {
		panic(err)
	}
	return attribute
}

func (r *Rep) FindMediaLocations(signature catalog.MediaSignature) ([]*catalog.MediaLocation, error) {
	mediaPK := mediaPrimaryKey(r.RootOwner, &signature)
	queryValues, err := dynamodbattribute.MarshalMap(map[string]string{
		":mediaPK":    mediaPK.PK,
		":noMetadata": "$",
	})
	if err != nil {
		return nil, err
	}

	result, err := r.db.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: queryValues,
		KeyConditionExpression:    aws.String("PK = :mediaPK AND SK > :noMetadata"),
		TableName:                 &r.table,
	})
	if err != nil {
		return nil, err
	}

	location, orders, err := unmarshalMediaItems(result.Items)
	if location == nil {
		return nil, errors.Errorf("Location not found for media %+v", signature)
	}

	locations := make([]*catalog.MediaLocation, len(orders)+1)
	locations[0] = &catalog.MediaLocation{
		FolderName: location.FolderName,
		Filename:   location.Filename,
	}
	for index, order := range orders {
		locations[index+1] = &catalog.MediaLocation{
			FolderName: order.DestinationFolder,
			Filename:   location.Filename,
		}
	}

	return locations, nil
}
