package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
)

// Stream is inspired from Java streams to chain transformations in a functional programming style
type Stream interface {
	HasNext() bool                             // HasNext returns true if it has another element
	Next() map[string]*dynamodb.AttributeValue // Next return current element and move forward the cursor
	Error() error                              // Error returns the error that interrupted the Stream
	Count() int64                              // Count return the number of element found so far
}

// arrayStream is a wrapper around a slice ; this slice can be replaced when IsLast() is true
type arrayStream struct {
	err     error
	results []map[string]*dynamodb.AttributeValue
	count   int64
}

func NewArrayStream(results []map[string]*dynamodb.AttributeValue) Stream {
	return &arrayStream{
		results: results,
	}
}

type queryStream struct {
	*rep
	internalStream arrayStream                         // arrayStream is wrapping the current internalStream
	queries        []*dynamodb.QueryInput              // queries first item is the one to use as long as queryHasNextPage is true
	nextPageToken  map[string]*dynamodb.AttributeValue // nextPageToken is only usable when queryHasNextPage is true
}

func NewQueryStream(rep *rep, queries []*dynamodb.QueryInput) Stream {
	stream := &queryStream{
		rep:            rep,
		internalStream: arrayStream{},
		queries:        queries,
	}
	stream.populateNextChunk()
	return stream
}

type getStream struct {
	db                   *dynamodb.DynamoDB
	table                string
	projectionExpression *string
	internalStream       arrayStream
	keys                 []map[string]*dynamodb.AttributeValue
	buffer               []map[string]*dynamodb.AttributeValue
}

func NewGetStream(rep *rep, keys []map[string]*dynamodb.AttributeValue, projectionExpression *string, bufferSize int64) Stream {
	stream := &getStream{
		buffer:               make([]map[string]*dynamodb.AttributeValue, 0, bufferSize),
		db:                   rep.db,
		internalStream:       arrayStream{},
		keys:                 keys,
		projectionExpression: projectionExpression,
		table:                rep.table,
	}
	stream.populateNextChunk()
	return stream
}

// ** ARRAY STREAM

func (s *arrayStream) HasNext() bool {
	return s.err == nil && len(s.results) > 0
}

func (s *arrayStream) IsLast() bool {
	return s.err == nil && len(s.results) > 1
}

func (s *arrayStream) Next() (current map[string]*dynamodb.AttributeValue) {
	if s.err != nil {
		panic(errors.Wrapf(s.err, "Next() can't be called when an error occured"))
	}

	if !s.HasNext() {
		panic(errors.Wrapf(s.err, "Next() can't be called when an HasNext() returns false"))
	}

	current, s.results = s.results[0], s.results[1:]
	s.count++

	return current
}

func (s *arrayStream) AppendNextChunk(chunk []map[string]*dynamodb.AttributeValue) {
	results := make([]map[string]*dynamodb.AttributeValue, len(s.results)+len(chunk))
	copy(results, s.results)
	copy(results[len(s.results):], chunk)

	s.results = results
}

func (s *arrayStream) Count() int64 {
	return s.count
}

func (s *arrayStream) Error() error {
	return s.err
}

// WithError keeps exiting error if already present. err can be nil.
func (s *arrayStream) WithError(err error) {
	if s.err == nil {
		s.err = err
	}
}

// ** QUERY STREAM

func (s *queryStream) HasNext() bool {
	return s.internalStream.HasNext()
}

func (s *queryStream) Next() map[string]*dynamodb.AttributeValue {
	next := s.internalStream.Next()

	for !s.internalStream.HasNext() && len(s.queries) > 0 {

		if len(s.nextPageToken) == 0 {
			s.queries = s.queries[1:]
		}

		if len(s.queries) > 0 {
			s.populateNextChunk()
		}
	}

	return next
}

func (s *queryStream) populateNextChunk() {
	query := *s.queries[0]
	query.ExclusiveStartKey = s.nextPageToken
	result, err := s.db.Query(&query)
	if err != nil {
		s.internalStream.WithError(err)
		return
	}

	s.nextPageToken = result.LastEvaluatedKey
	s.internalStream.AppendNextChunk(result.Items)
}

func (s *queryStream) Count() int64 {
	return s.internalStream.count
}

func (s *queryStream) Error() error {
	return s.internalStream.Error()
}

// ** GET STREAM

func (s *getStream) HasNext() bool {
	return s.internalStream.HasNext()
}

func (s *getStream) Next() map[string]*dynamodb.AttributeValue {
	next := s.internalStream.Next()

	for !s.internalStream.HasNext() && (len(s.buffer) > 0 || len(s.keys) > 0) {
		s.populateNextChunk()
	}

	return next
}

func (s *getStream) Count() int64 {
	return s.internalStream.Count()
}

func (s *getStream) Error() error {
	return s.internalStream.Error()
}

func (s *getStream) populateNextChunk() {
	end := cap(s.buffer) - len(s.buffer)
	if end > len(s.keys) {
		end = len(s.keys)
	}

	if end > 0 {
		s.buffer = append(s.buffer, s.keys[:end]...)
		s.keys = s.keys[end:]
	}

	if len(s.buffer) == 0 {
		return
	}

	result, err := s.db.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			s.table: {
				Keys:                 s.buffer,
				ProjectionExpression: s.projectionExpression,
			},
		},
	})
	if err != nil {
		s.internalStream.WithError(err)
		return
	}

	s.buffer = s.buffer[:0]
	if table, ok := result.UnprocessedKeys[s.table]; ok && len(table.Keys) > 0 {
		s.buffer = append(s.buffer, table.Keys...)
	}

	s.internalStream.AppendNextChunk(result.Responses[s.table])
}
