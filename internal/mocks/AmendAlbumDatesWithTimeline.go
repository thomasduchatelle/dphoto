// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// AmendAlbumDatesWithTimeline is an autogenerated mock type for the AmendAlbumDatesWithTimeline type
type AmendAlbumDatesWithTimeline struct {
	mock.Mock
}

type AmendAlbumDatesWithTimeline_Expecter struct {
	mock *mock.Mock
}

func (_m *AmendAlbumDatesWithTimeline) EXPECT() *AmendAlbumDatesWithTimeline_Expecter {
	return &AmendAlbumDatesWithTimeline_Expecter{mock: &_m.Mock}
}

// AmendAlbumDates provides a mock function with given fields: ctx, timeline, albumId, start, end
func (_m *AmendAlbumDatesWithTimeline) AmendAlbumDates(ctx context.Context, timeline *catalog.TimelineAggregate, albumId catalog.AlbumId, start time.Time, end time.Time) error {
	ret := _m.Called(ctx, timeline, albumId, start, end)

	if len(ret) == 0 {
		panic("no return value specified for AmendAlbumDates")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *catalog.TimelineAggregate, catalog.AlbumId, time.Time, time.Time) error); ok {
		r0 = rf(ctx, timeline, albumId, start, end)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AmendAlbumDatesWithTimeline_AmendAlbumDates_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AmendAlbumDates'
type AmendAlbumDatesWithTimeline_AmendAlbumDates_Call struct {
	*mock.Call
}

// AmendAlbumDates is a helper method to define mock.On call
//   - ctx context.Context
//   - timeline *catalog.TimelineAggregate
//   - albumId catalog.AlbumId
//   - start time.Time
//   - end time.Time
func (_e *AmendAlbumDatesWithTimeline_Expecter) AmendAlbumDates(ctx interface{}, timeline interface{}, albumId interface{}, start interface{}, end interface{}) *AmendAlbumDatesWithTimeline_AmendAlbumDates_Call {
	return &AmendAlbumDatesWithTimeline_AmendAlbumDates_Call{Call: _e.mock.On("AmendAlbumDates", ctx, timeline, albumId, start, end)}
}

func (_c *AmendAlbumDatesWithTimeline_AmendAlbumDates_Call) Run(run func(ctx context.Context, timeline *catalog.TimelineAggregate, albumId catalog.AlbumId, start time.Time, end time.Time)) *AmendAlbumDatesWithTimeline_AmendAlbumDates_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*catalog.TimelineAggregate), args[2].(catalog.AlbumId), args[3].(time.Time), args[4].(time.Time))
	})
	return _c
}

func (_c *AmendAlbumDatesWithTimeline_AmendAlbumDates_Call) Return(_a0 error) *AmendAlbumDatesWithTimeline_AmendAlbumDates_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AmendAlbumDatesWithTimeline_AmendAlbumDates_Call) RunAndReturn(run func(context.Context, *catalog.TimelineAggregate, catalog.AlbumId, time.Time, time.Time) error) *AmendAlbumDatesWithTimeline_AmendAlbumDates_Call {
	_c.Call.Return(run)
	return _c
}

// NewAmendAlbumDatesWithTimeline creates a new instance of AmendAlbumDatesWithTimeline. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAmendAlbumDatesWithTimeline(t interface {
	mock.TestingT
	Cleanup(func())
}) *AmendAlbumDatesWithTimeline {
	mock := &AmendAlbumDatesWithTimeline{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
