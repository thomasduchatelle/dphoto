// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalogviews "github.com/thomasduchatelle/dphoto/pkg/catalogviews"

	mock "github.com/stretchr/testify/mock"
)

// InsertAlbumSizePort is an autogenerated mock type for the InsertAlbumSizePort type
type InsertAlbumSizePort struct {
	mock.Mock
}

type InsertAlbumSizePort_Expecter struct {
	mock *mock.Mock
}

func (_m *InsertAlbumSizePort) EXPECT() *InsertAlbumSizePort_Expecter {
	return &InsertAlbumSizePort_Expecter{mock: &_m.Mock}
}

// InsertAlbumSize provides a mock function with given fields: ctx, albumSize
func (_m *InsertAlbumSizePort) InsertAlbumSize(ctx context.Context, albumSize []catalogviews.MultiUserAlbumSize) error {
	ret := _m.Called(ctx, albumSize)

	if len(ret) == 0 {
		panic("no return value specified for InsertAlbumSize")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []catalogviews.MultiUserAlbumSize) error); ok {
		r0 = rf(ctx, albumSize)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertAlbumSizePort_InsertAlbumSize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertAlbumSize'
type InsertAlbumSizePort_InsertAlbumSize_Call struct {
	*mock.Call
}

// InsertAlbumSize is a helper method to define mock.On call
//   - ctx context.Context
//   - albumSize []catalogviews.MultiUserAlbumSize
func (_e *InsertAlbumSizePort_Expecter) InsertAlbumSize(ctx interface{}, albumSize interface{}) *InsertAlbumSizePort_InsertAlbumSize_Call {
	return &InsertAlbumSizePort_InsertAlbumSize_Call{Call: _e.mock.On("InsertAlbumSize", ctx, albumSize)}
}

func (_c *InsertAlbumSizePort_InsertAlbumSize_Call) Run(run func(ctx context.Context, albumSize []catalogviews.MultiUserAlbumSize)) *InsertAlbumSizePort_InsertAlbumSize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]catalogviews.MultiUserAlbumSize))
	})
	return _c
}

func (_c *InsertAlbumSizePort_InsertAlbumSize_Call) Return(_a0 error) *InsertAlbumSizePort_InsertAlbumSize_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *InsertAlbumSizePort_InsertAlbumSize_Call) RunAndReturn(run func(context.Context, []catalogviews.MultiUserAlbumSize) error) *InsertAlbumSizePort_InsertAlbumSize_Call {
	_c.Call.Return(run)
	return _c
}

// NewInsertAlbumSizePort creates a new instance of InsertAlbumSizePort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInsertAlbumSizePort(t interface {
	mock.TestingT
	Cleanup(func())
}) *InsertAlbumSizePort {
	mock := &InsertAlbumSizePort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}