package dynamoutils

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/pkg/errors"
)

type arrayStream struct {
	err     error
	results []map[string]*dynamodb.AttributeValue
	count   int64
}

// NewArrayStream creates a stream from a slice ; the slice will be updated.
func NewArrayStream(results []map[string]*dynamodb.AttributeValue) Stream {
	return &arrayStream{
		results: results,
	}
}

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

func (s *arrayStream) appendNextChunk(chunk []map[string]*dynamodb.AttributeValue) {
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
