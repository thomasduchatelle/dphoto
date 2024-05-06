// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// DeleteAlbumRepositoryPort is an autogenerated mock type for the DeleteAlbumRepositoryPort type
type DeleteAlbumRepositoryPort struct {
	mock.Mock
}

type DeleteAlbumRepositoryPort_Expecter struct {
	mock *mock.Mock
}

func (_m *DeleteAlbumRepositoryPort) EXPECT() *DeleteAlbumRepositoryPort_Expecter {
	return &DeleteAlbumRepositoryPort_Expecter{mock: &_m.Mock}
}

// DeleteAlbum provides a mock function with given fields: ctx, albumId
func (_m *DeleteAlbumRepositoryPort) DeleteAlbum(ctx context.Context, albumId catalog.AlbumId) error {
	ret := _m.Called(ctx, albumId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAlbum")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, catalog.AlbumId) error); ok {
		r0 = rf(ctx, albumId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteAlbumRepositoryPort_DeleteAlbum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteAlbum'
type DeleteAlbumRepositoryPort_DeleteAlbum_Call struct {
	*mock.Call
}

// DeleteAlbum is a helper method to define mock.On call
//   - ctx context.Context
//   - albumId catalog.AlbumId
func (_e *DeleteAlbumRepositoryPort_Expecter) DeleteAlbum(ctx interface{}, albumId interface{}) *DeleteAlbumRepositoryPort_DeleteAlbum_Call {
	return &DeleteAlbumRepositoryPort_DeleteAlbum_Call{Call: _e.mock.On("DeleteAlbum", ctx, albumId)}
}

func (_c *DeleteAlbumRepositoryPort_DeleteAlbum_Call) Run(run func(ctx context.Context, albumId catalog.AlbumId)) *DeleteAlbumRepositoryPort_DeleteAlbum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(catalog.AlbumId))
	})
	return _c
}

func (_c *DeleteAlbumRepositoryPort_DeleteAlbum_Call) Return(_a0 error) *DeleteAlbumRepositoryPort_DeleteAlbum_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DeleteAlbumRepositoryPort_DeleteAlbum_Call) RunAndReturn(run func(context.Context, catalog.AlbumId) error) *DeleteAlbumRepositoryPort_DeleteAlbum_Call {
	_c.Call.Return(run)
	return _c
}

// NewDeleteAlbumRepositoryPort creates a new instance of DeleteAlbumRepositoryPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDeleteAlbumRepositoryPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *DeleteAlbumRepositoryPort {
	mock := &DeleteAlbumRepositoryPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}