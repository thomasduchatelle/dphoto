package dynamoutils

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
)

type queryStream struct {
	executor       DynamoQuery
	internalStream arrayStream                         // arrayStream is wrapping the current internalStream
	queries        []*dynamodb.QueryInput              // queries first item is the one to use as long as queryHasNextPage is true
	nextPageToken  map[string]*dynamodb.AttributeValue // nextPageToken is only usable when queryHasNextPage is true
}

// NewQueryStream creates a stream that will execute the list of queries.
// If a query is paginated, all the pages will be requested before moving on the next query.
func NewQueryStream(executor DynamoQuery, queries []*dynamodb.QueryInput) Stream {
	if len(queries) == 0 {
		return NewArrayStream(nil)
	}

	stream := &queryStream{
		executor:       executor,
		internalStream: arrayStream{},
		queries:        queries,
	}
	stream.populateNextChunk()
	return stream
}

func (s *queryStream) HasNext() bool {
	return s.internalStream.HasNext()
}

func (s *queryStream) Next() map[string]*dynamodb.AttributeValue {
	next := s.internalStream.Next()

	if !s.internalStream.HasNext() {
		s.populateNextChunk()
	}

	return next
}

func (s *queryStream) populateNextChunk() {
	for !s.internalStream.HasNext() && len(s.queries) > 0 {
		query := *s.queries[0]
		query.ExclusiveStartKey = s.nextPageToken
		result, err := s.executor.Query(&query)
		if err != nil {
			s.internalStream.WithError(errors.Wrapf(err, "couldn't query %+v", query))
			return
		}

		s.nextPageToken = result.LastEvaluatedKey
		s.internalStream.appendNextChunk(result.Items)

		if len(s.nextPageToken) == 0 {
			s.queries = s.queries[1:]
		}
	}
}

func (s *queryStream) Count() int64 {
	return s.internalStream.count
}

func (s *queryStream) Error() error {
	return s.internalStream.Error()
}
