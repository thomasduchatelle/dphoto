// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// AmendAlbumDateRepositoryPort is an autogenerated mock type for the AmendAlbumDateRepositoryPort type
type AmendAlbumDateRepositoryPort struct {
	mock.Mock
}

type AmendAlbumDateRepositoryPort_Expecter struct {
	mock *mock.Mock
}

func (_m *AmendAlbumDateRepositoryPort) EXPECT() *AmendAlbumDateRepositoryPort_Expecter {
	return &AmendAlbumDateRepositoryPort_Expecter{mock: &_m.Mock}
}

// AmendDates provides a mock function with given fields: ctx, album, start, end
func (_m *AmendAlbumDateRepositoryPort) AmendDates(ctx context.Context, album catalog.AlbumId, start time.Time, end time.Time) error {
	ret := _m.Called(ctx, album, start, end)

	if len(ret) == 0 {
		panic("no return value specified for AmendDates")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, catalog.AlbumId, time.Time, time.Time) error); ok {
		r0 = rf(ctx, album, start, end)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AmendAlbumDateRepositoryPort_AmendDates_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AmendDates'
type AmendAlbumDateRepositoryPort_AmendDates_Call struct {
	*mock.Call
}

// AmendDates is a helper method to define mock.On call
//   - ctx context.Context
//   - album catalog.AlbumId
//   - start time.Time
//   - end time.Time
func (_e *AmendAlbumDateRepositoryPort_Expecter) AmendDates(ctx interface{}, album interface{}, start interface{}, end interface{}) *AmendAlbumDateRepositoryPort_AmendDates_Call {
	return &AmendAlbumDateRepositoryPort_AmendDates_Call{Call: _e.mock.On("AmendDates", ctx, album, start, end)}
}

func (_c *AmendAlbumDateRepositoryPort_AmendDates_Call) Run(run func(ctx context.Context, album catalog.AlbumId, start time.Time, end time.Time)) *AmendAlbumDateRepositoryPort_AmendDates_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(catalog.AlbumId), args[2].(time.Time), args[3].(time.Time))
	})
	return _c
}

func (_c *AmendAlbumDateRepositoryPort_AmendDates_Call) Return(_a0 error) *AmendAlbumDateRepositoryPort_AmendDates_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AmendAlbumDateRepositoryPort_AmendDates_Call) RunAndReturn(run func(context.Context, catalog.AlbumId, time.Time, time.Time) error) *AmendAlbumDateRepositoryPort_AmendDates_Call {
	_c.Call.Return(run)
	return _c
}

// NewAmendAlbumDateRepositoryPort creates a new instance of AmendAlbumDateRepositoryPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAmendAlbumDateRepositoryPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *AmendAlbumDateRepositoryPort {
	mock := &AmendAlbumDateRepositoryPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
