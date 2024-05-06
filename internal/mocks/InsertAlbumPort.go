// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"
)

// InsertAlbumPort is an autogenerated mock type for the InsertAlbumPort type
type InsertAlbumPort struct {
	mock.Mock
}

type InsertAlbumPort_Expecter struct {
	mock *mock.Mock
}

func (_m *InsertAlbumPort) EXPECT() *InsertAlbumPort_Expecter {
	return &InsertAlbumPort_Expecter{mock: &_m.Mock}
}

// InsertAlbum provides a mock function with given fields: ctx, album
func (_m *InsertAlbumPort) InsertAlbum(ctx context.Context, album catalog.Album) error {
	ret := _m.Called(ctx, album)

	if len(ret) == 0 {
		panic("no return value specified for InsertAlbum")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, catalog.Album) error); ok {
		r0 = rf(ctx, album)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertAlbumPort_InsertAlbum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertAlbum'
type InsertAlbumPort_InsertAlbum_Call struct {
	*mock.Call
}

// InsertAlbum is a helper method to define mock.On call
//   - ctx context.Context
//   - album catalog.Album
func (_e *InsertAlbumPort_Expecter) InsertAlbum(ctx interface{}, album interface{}) *InsertAlbumPort_InsertAlbum_Call {
	return &InsertAlbumPort_InsertAlbum_Call{Call: _e.mock.On("InsertAlbum", ctx, album)}
}

func (_c *InsertAlbumPort_InsertAlbum_Call) Run(run func(ctx context.Context, album catalog.Album)) *InsertAlbumPort_InsertAlbum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(catalog.Album))
	})
	return _c
}

func (_c *InsertAlbumPort_InsertAlbum_Call) Return(_a0 error) *InsertAlbumPort_InsertAlbum_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *InsertAlbumPort_InsertAlbum_Call) RunAndReturn(run func(context.Context, catalog.Album) error) *InsertAlbumPort_InsertAlbum_Call {
	_c.Call.Return(run)
	return _c
}

// NewInsertAlbumPort creates a new instance of InsertAlbumPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInsertAlbumPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *InsertAlbumPort {
	mock := &InsertAlbumPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
