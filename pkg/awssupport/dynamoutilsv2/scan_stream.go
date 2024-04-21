package dynamoutilsv2

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
)

type scanStream struct {
	ctx            context.Context
	executor       ScanStreamExecutor
	internalStream arrayStream
	nextPageToken  map[string]types.AttributeValue
	tableName      string
}

type ScanStreamExecutor interface {
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

func NewScanStream(ctx context.Context, executor ScanStreamExecutor, tableName string) Stream {
	stream := &scanStream{
		ctx:            ctx,
		executor:       executor,
		internalStream: arrayStream{},
		tableName:      tableName,
	}
	stream.populateNextChunk()
	return stream
}

func (s *scanStream) HasNext() bool {
	return s.internalStream.HasNext()
}

func (s *scanStream) Next() map[string]types.AttributeValue {
	next := s.internalStream.Next()

	if !s.internalStream.HasNext() {
		if len(s.nextPageToken) != 0 {
			s.populateNextChunk()
		}
	}

	return next
}

func (s *scanStream) populateNextChunk() {
	result, err := s.executor.Scan(s.ctx, &dynamodb.ScanInput{
		ExclusiveStartKey: s.nextPageToken,
		TableName:         &s.tableName,
	})
	if err != nil {
		s.internalStream.WithError(errors.Wrapf(err, "couldn't scan %s", s.tableName))
		return
	}

	s.nextPageToken = result.LastEvaluatedKey
	s.internalStream.appendNextChunk(result.Items)
}

func (s *scanStream) Count() int64 {
	return s.internalStream.count
}

func (s *scanStream) Error() error {
	return s.internalStream.Error()
}
