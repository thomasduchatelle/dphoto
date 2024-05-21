// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	context "context"

	mock "github.com/stretchr/testify/mock"
)

// FindAlbumsByIdsFunc is an autogenerated mock type for the FindAlbumsByIdsFunc type
type FindAlbumsByIdsFunc struct {
	mock.Mock
}

type FindAlbumsByIdsFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *FindAlbumsByIdsFunc) EXPECT() *FindAlbumsByIdsFunc_Expecter {
	return &FindAlbumsByIdsFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, ids
func (_m *FindAlbumsByIdsFunc) Execute(ctx context.Context, ids []catalog.AlbumId) ([]*catalog.Album, error) {
	ret := _m.Called(ctx, ids)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 []*catalog.Album
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []catalog.AlbumId) ([]*catalog.Album, error)); ok {
		return rf(ctx, ids)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []catalog.AlbumId) []*catalog.Album); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*catalog.Album)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []catalog.AlbumId) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAlbumsByIdsFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type FindAlbumsByIdsFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - ids []catalog.AlbumId
func (_e *FindAlbumsByIdsFunc_Expecter) Execute(ctx interface{}, ids interface{}) *FindAlbumsByIdsFunc_Execute_Call {
	return &FindAlbumsByIdsFunc_Execute_Call{Call: _e.mock.On("Execute", ctx, ids)}
}

func (_c *FindAlbumsByIdsFunc_Execute_Call) Run(run func(ctx context.Context, ids []catalog.AlbumId)) *FindAlbumsByIdsFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]catalog.AlbumId))
	})
	return _c
}

func (_c *FindAlbumsByIdsFunc_Execute_Call) Return(_a0 []*catalog.Album, _a1 error) *FindAlbumsByIdsFunc_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FindAlbumsByIdsFunc_Execute_Call) RunAndReturn(run func(context.Context, []catalog.AlbumId) ([]*catalog.Album, error)) *FindAlbumsByIdsFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewFindAlbumsByIdsFunc creates a new instance of FindAlbumsByIdsFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFindAlbumsByIdsFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *FindAlbumsByIdsFunc {
	mock := &FindAlbumsByIdsFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}