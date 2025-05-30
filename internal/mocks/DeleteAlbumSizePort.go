// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"
	catalogviews "github.com/thomasduchatelle/dphoto/pkg/catalogviews"

	context "context"

	mock "github.com/stretchr/testify/mock"
)

// DeleteAlbumSizePort is an autogenerated mock type for the DeleteAlbumSizePort type
type DeleteAlbumSizePort struct {
	mock.Mock
}

type DeleteAlbumSizePort_Expecter struct {
	mock *mock.Mock
}

func (_m *DeleteAlbumSizePort) EXPECT() *DeleteAlbumSizePort_Expecter {
	return &DeleteAlbumSizePort_Expecter{mock: &_m.Mock}
}

// DeleteAlbumSize provides a mock function with given fields: ctx, availability, redirectTo
func (_m *DeleteAlbumSizePort) DeleteAlbumSize(ctx context.Context, availability catalogviews.Availability, albumId catalog.AlbumId) error {
	ret := _m.Called(ctx, availability, albumId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAlbumSize")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, catalogviews.Availability, catalog.AlbumId) error); ok {
		r0 = rf(ctx, availability, albumId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAlbumSizePort_DeleteAlbumSize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteAlbumSize'
type DeleteAlbumSizePort_DeleteAlbumSize_Call struct {
	*mock.Call
}

// DeleteAlbumSize is a helper method to define mock.On call
//   - ctx context.Context
//   - availability catalogviews.Availability
//   - redirectTo catalog.AlbumId
func (_e *DeleteAlbumSizePort_Expecter) DeleteAlbumSize(ctx interface{}, availability interface{}, albumId interface{}) *DeleteAlbumSizePort_DeleteAlbumSize_Call {
	return &DeleteAlbumSizePort_DeleteAlbumSize_Call{Call: _e.mock.On("DeleteAlbumSize", ctx, availability, albumId)}
}

func (_c *DeleteAlbumSizePort_DeleteAlbumSize_Call) Run(run func(ctx context.Context, availability catalogviews.Availability, albumId catalog.AlbumId)) *DeleteAlbumSizePort_DeleteAlbumSize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(catalogviews.Availability), args[2].(catalog.AlbumId))
	})
	return _c
}

func (_c *DeleteAlbumSizePort_DeleteAlbumSize_Call) Return(_a0 error) *DeleteAlbumSizePort_DeleteAlbumSize_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DeleteAlbumSizePort_DeleteAlbumSize_Call) RunAndReturn(run func(context.Context, catalogviews.Availability, catalog.AlbumId) error) *DeleteAlbumSizePort_DeleteAlbumSize_Call {
	_c.Call.Return(run)
	return _c
}

// NewDeleteAlbumSizePort creates a new instance of DeleteAlbumSizePort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDeleteAlbumSizePort(t interface {
	mock.TestingT
	Cleanup(func())
}) *DeleteAlbumSizePort {
	mock := &DeleteAlbumSizePort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
