package dynamoutils

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoGetStreamExecutor creates the query and execute in on the dynamoDB
type DynamoGetStreamExecutor struct {
	db                   DynamoBatchGetItem
	table                string
	projectionExpression *string
}

type getStream struct {
	ctx            context.Context
	executor       GetStreamAdapter
	internalStream arrayStream
	keys           []map[string]types.AttributeValue
	buffer         []map[string]types.AttributeValue
}

// NewGetStream batches requests to get each item by their natural key
func NewGetStream(ctx context.Context, executor GetStreamAdapter, keys []map[string]types.AttributeValue, bufferSize int64) Stream {
	stream := &getStream{
		ctx:            ctx,
		executor:       executor,
		buffer:         make([]map[string]types.AttributeValue, 0, bufferSize),
		internalStream: arrayStream{},
		keys:           keys,
	}
	stream.populateNextChunk()
	return stream
}

func (s *getStream) HasNext() bool {
	return s.internalStream.HasNext()
}

func (s *getStream) Next() map[string]types.AttributeValue {
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

	result, err := s.executor.BatchGet(s.ctx, s.buffer)
	if err != nil {
		s.internalStream.WithError(err)
		return
	}

	s.buffer = s.buffer[:0]
	for _, unprocessedKeys := range result.UnprocessedKeys {
		s.buffer = append(s.buffer, unprocessedKeys.Keys...)
	}

	for _, records := range result.Responses {
		s.internalStream.appendNextChunk(records)
	}
}
