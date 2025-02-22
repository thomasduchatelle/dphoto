// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"
	catalogviews "github.com/thomasduchatelle/dphoto/pkg/catalogviews"

	context "context"

	mock "github.com/stretchr/testify/mock"
)

// DriftSynchronizerPort is an autogenerated mock type for the DriftSynchronizerPort type
type DriftSynchronizerPort struct {
	mock.Mock
}

type DriftSynchronizerPort_Expecter struct {
	mock *mock.Mock
}

func (_m *DriftSynchronizerPort) EXPECT() *DriftSynchronizerPort_Expecter {
	return &DriftSynchronizerPort_Expecter{mock: &_m.Mock}
}

// DeleteAlbumSize provides a mock function with given fields: ctx, availability, redirectTo
func (_m *DriftSynchronizerPort) DeleteAlbumSize(ctx context.Context, availability catalogviews.Availability, albumId catalog.AlbumId) error {
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

// DriftSynchronizerPort_DeleteAlbumSize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteAlbumSize'
type DriftSynchronizerPort_DeleteAlbumSize_Call struct {
	*mock.Call
}

// DeleteAlbumSize is a helper method to define mock.On call
//   - ctx context.Context
//   - availability catalogviews.Availability
//   - redirectTo catalog.AlbumId
func (_e *DriftSynchronizerPort_Expecter) DeleteAlbumSize(ctx interface{}, availability interface{}, albumId interface{}) *DriftSynchronizerPort_DeleteAlbumSize_Call {
	return &DriftSynchronizerPort_DeleteAlbumSize_Call{Call: _e.mock.On("DeleteAlbumSize", ctx, availability, albumId)}
}

func (_c *DriftSynchronizerPort_DeleteAlbumSize_Call) Run(run func(ctx context.Context, availability catalogviews.Availability, albumId catalog.AlbumId)) *DriftSynchronizerPort_DeleteAlbumSize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(catalogviews.Availability), args[2].(catalog.AlbumId))
	})
	return _c
}

func (_c *DriftSynchronizerPort_DeleteAlbumSize_Call) Return(_a0 error) *DriftSynchronizerPort_DeleteAlbumSize_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DriftSynchronizerPort_DeleteAlbumSize_Call) RunAndReturn(run func(context.Context, catalogviews.Availability, catalog.AlbumId) error) *DriftSynchronizerPort_DeleteAlbumSize_Call {
	_c.Call.Return(run)
	return _c
}

// InsertAlbumSize provides a mock function with given fields: ctx, albumSize
func (_m *DriftSynchronizerPort) InsertAlbumSize(ctx context.Context, albumSize []catalogviews.MultiUserAlbumSize) error {
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

// DriftSynchronizerPort_InsertAlbumSize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertAlbumSize'
type DriftSynchronizerPort_InsertAlbumSize_Call struct {
	*mock.Call
}

// InsertAlbumSize is a helper method to define mock.On call
//   - ctx context.Context
//   - albumSize []catalogviews.MultiUserAlbumSize
func (_e *DriftSynchronizerPort_Expecter) InsertAlbumSize(ctx interface{}, albumSize interface{}) *DriftSynchronizerPort_InsertAlbumSize_Call {
	return &DriftSynchronizerPort_InsertAlbumSize_Call{Call: _e.mock.On("InsertAlbumSize", ctx, albumSize)}
}

func (_c *DriftSynchronizerPort_InsertAlbumSize_Call) Run(run func(ctx context.Context, albumSize []catalogviews.MultiUserAlbumSize)) *DriftSynchronizerPort_InsertAlbumSize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]catalogviews.MultiUserAlbumSize))
	})
	return _c
}

func (_c *DriftSynchronizerPort_InsertAlbumSize_Call) Return(_a0 error) *DriftSynchronizerPort_InsertAlbumSize_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DriftSynchronizerPort_InsertAlbumSize_Call) RunAndReturn(run func(context.Context, []catalogviews.MultiUserAlbumSize) error) *DriftSynchronizerPort_InsertAlbumSize_Call {
	_c.Call.Return(run)
	return _c
}

// NewDriftSynchronizerPort creates a new instance of DriftSynchronizerPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDriftSynchronizerPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *DriftSynchronizerPort {
	mock := &DriftSynchronizerPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
