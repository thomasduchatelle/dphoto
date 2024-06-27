// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// AlbumDatesAmendedObserverFunc is an autogenerated mock type for the AlbumDatesAmendedObserverFunc type
type AlbumDatesAmendedObserverFunc struct {
	mock.Mock
}

type AlbumDatesAmendedObserverFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *AlbumDatesAmendedObserverFunc) EXPECT() *AlbumDatesAmendedObserverFunc_Expecter {
	return &AlbumDatesAmendedObserverFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, amendedAlbum
func (_m *AlbumDatesAmendedObserverFunc) Execute(ctx context.Context, amendedAlbum catalog.DatesUpdate) error {
	ret := _m.Called(ctx, amendedAlbum)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, catalog.DatesUpdate) error); ok {
		r0 = rf(ctx, amendedAlbum)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AlbumDatesAmendedObserverFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type AlbumDatesAmendedObserverFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - amendedAlbum catalog.DatesUpdate
func (_e *AlbumDatesAmendedObserverFunc_Expecter) Execute(ctx interface{}, amendedAlbum interface{}) *AlbumDatesAmendedObserverFunc_Execute_Call {
	return &AlbumDatesAmendedObserverFunc_Execute_Call{Call: _e.mock.On("Execute", ctx, amendedAlbum)}
}

func (_c *AlbumDatesAmendedObserverFunc_Execute_Call) Run(run func(ctx context.Context, amendedAlbum catalog.DatesUpdate)) *AlbumDatesAmendedObserverFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(catalog.DatesUpdate))
	})
	return _c
}

func (_c *AlbumDatesAmendedObserverFunc_Execute_Call) Return(_a0 error) *AlbumDatesAmendedObserverFunc_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AlbumDatesAmendedObserverFunc_Execute_Call) RunAndReturn(run func(context.Context, catalog.DatesUpdate) error) *AlbumDatesAmendedObserverFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewAlbumDatesAmendedObserverFunc creates a new instance of AlbumDatesAmendedObserverFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAlbumDatesAmendedObserverFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *AlbumDatesAmendedObserverFunc {
	mock := &AlbumDatesAmendedObserverFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}