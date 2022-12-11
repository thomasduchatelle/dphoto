package dynamoutils

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
)

type scanStream struct {
	executor       ScanStreamExecutor
	internalStream arrayStream
	nextPageToken  map[string]*dynamodb.AttributeValue
	tableName      string
}

type ScanStreamExecutor interface {
	Scan(input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error)
}

func NewScanStream(executor ScanStreamExecutor, tableName string) Stream {
	stream := &scanStream{
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

func (s *scanStream) Next() map[string]*dynamodb.AttributeValue {
	next := s.internalStream.Next()

	if !s.internalStream.HasNext() {
		if len(s.nextPageToken) != 0 {
			s.populateNextChunk()
		}
	}

	return next
}

func (s *scanStream) populateNextChunk() {
	result, err := s.executor.Scan(&dynamodb.ScanInput{
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
