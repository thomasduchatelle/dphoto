package dynamoutilsv2

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/pkg/errors"
)

type queryStream struct {
	ctx            context.Context
	executor       dynamodb.QueryAPIClient
	queries        []*dynamodb.QueryInput   // queries first item is the one to use as long as queryHasNextPage is true
	paginator      *dynamodb.QueryPaginator // paginator is used if not nil and still have next
	internalStream arrayStream              // arrayStream is wrapping the result of the current query page

}

// NewQueryStream creates a stream that will execute the list of queries.
// If a query is paginated, all the pages will be requested before moving on the next query.
func NewQueryStream(ctx context.Context, executor DynamoQuery, queries []*dynamodb.QueryInput) Stream {
	if len(queries) == 0 {
		return NewArrayStream(nil)
	}

	stream := &queryStream{
		ctx:            ctx,
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

func (s *queryStream) Next() map[string]types.AttributeValue {
	next := s.internalStream.Next()

	if !s.internalStream.HasNext() {
		s.populateNextChunk()
	}

	return next
}

func (s *queryStream) populateNextChunk() {
	for !s.internalStream.HasNext() {
		if s.paginator != nil && s.paginator.HasMorePages() {
			page, err := s.paginator.NextPage(s.ctx)
			if err != nil {
				s.internalStream.WithError(errors.Wrapf(err, "couldn't query %+v", s.paginator))
				return
			}

			s.internalStream.appendNextChunk(page.Items)

		} else if len(s.queries) > 0 {
			query := s.queries[0]
			s.queries = s.queries[1:]

			s.paginator = dynamodb.NewQueryPaginator(s.executor, query)
			// ... will fetch data on the next on the next loop

		} else {
			// end of the line
			return
		}
	}
}

func (s *queryStream) Count() int64 {
	return s.internalStream.count
}

func (s *queryStream) Error() error {
	return s.internalStream.Error()
}
