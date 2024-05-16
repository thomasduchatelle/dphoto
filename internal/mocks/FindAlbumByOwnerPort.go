// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	context "context"

	mock "github.com/stretchr/testify/mock"

	ownermodel "github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// FindAlbumByOwnerPort is an autogenerated mock type for the FindAlbumByOwnerPort type
type FindAlbumByOwnerPort struct {
	mock.Mock
}

type FindAlbumByOwnerPort_Expecter struct {
	mock *mock.Mock
}

func (_m *FindAlbumByOwnerPort) EXPECT() *FindAlbumByOwnerPort_Expecter {
	return &FindAlbumByOwnerPort_Expecter{mock: &_m.Mock}
}

// FindAlbumsByOwner provides a mock function with given fields: ctx, owner
func (_m *FindAlbumByOwnerPort) FindAlbumsByOwner(ctx context.Context, owner ownermodel.Owner) ([]*catalog.Album, error) {
	ret := _m.Called(ctx, owner)

	if len(ret) == 0 {
		panic("no return value specified for FindAlbumsByOwner")
	}

	var r0 []*catalog.Album
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner) ([]*catalog.Album, error)); ok {
		return rf(ctx, owner)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner) []*catalog.Album); ok {
		r0 = rf(ctx, owner)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*catalog.Album)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ownermodel.Owner) error); ok {
		r1 = rf(ctx, owner)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindAlbumByOwnerPort_FindAlbumsByOwner_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindAlbumsByOwner'
type FindAlbumByOwnerPort_FindAlbumsByOwner_Call struct {
	*mock.Call
}

// FindAlbumsByOwner is a helper method to define mock.On call
//   - ctx context.Context
//   - owner ownermodel.Owner
func (_e *FindAlbumByOwnerPort_Expecter) FindAlbumsByOwner(ctx interface{}, owner interface{}) *FindAlbumByOwnerPort_FindAlbumsByOwner_Call {
	return &FindAlbumByOwnerPort_FindAlbumsByOwner_Call{Call: _e.mock.On("FindAlbumsByOwner", ctx, owner)}
}

func (_c *FindAlbumByOwnerPort_FindAlbumsByOwner_Call) Run(run func(ctx context.Context, owner ownermodel.Owner)) *FindAlbumByOwnerPort_FindAlbumsByOwner_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ownermodel.Owner))
	})
	return _c
}

func (_c *FindAlbumByOwnerPort_FindAlbumsByOwner_Call) Return(_a0 []*catalog.Album, _a1 error) *FindAlbumByOwnerPort_FindAlbumsByOwner_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *FindAlbumByOwnerPort_FindAlbumsByOwner_Call) RunAndReturn(run func(context.Context, ownermodel.Owner) ([]*catalog.Album, error)) *FindAlbumByOwnerPort_FindAlbumsByOwner_Call {
	_c.Call.Return(run)
	return _c
}

// NewFindAlbumByOwnerPort creates a new instance of FindAlbumByOwnerPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFindAlbumByOwnerPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *FindAlbumByOwnerPort {
	mock := &FindAlbumByOwnerPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
