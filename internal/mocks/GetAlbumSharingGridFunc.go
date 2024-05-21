// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	context "context"

	mock "github.com/stretchr/testify/mock"

	ownermodel "github.com/thomasduchatelle/dphoto/pkg/ownermodel"

	usermodel "github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// GetAlbumSharingGridFunc is an autogenerated mock type for the GetAlbumSharingGridFunc type
type GetAlbumSharingGridFunc struct {
	mock.Mock
}

type GetAlbumSharingGridFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *GetAlbumSharingGridFunc) EXPECT() *GetAlbumSharingGridFunc_Expecter {
	return &GetAlbumSharingGridFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: ctx, owner
func (_m *GetAlbumSharingGridFunc) Execute(ctx context.Context, owner ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error) {
	ret := _m.Called(ctx, owner)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 map[catalog.AlbumId][]usermodel.UserId
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error)); ok {
		return rf(ctx, owner)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner) map[catalog.AlbumId][]usermodel.UserId); ok {
		r0 = rf(ctx, owner)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[catalog.AlbumId][]usermodel.UserId)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ownermodel.Owner) error); ok {
		r1 = rf(ctx, owner)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAlbumSharingGridFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type GetAlbumSharingGridFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - owner ownermodel.Owner
func (_e *GetAlbumSharingGridFunc_Expecter) Execute(ctx interface{}, owner interface{}) *GetAlbumSharingGridFunc_Execute_Call {
	return &GetAlbumSharingGridFunc_Execute_Call{Call: _e.mock.On("Execute", ctx, owner)}
}

func (_c *GetAlbumSharingGridFunc_Execute_Call) Run(run func(ctx context.Context, owner ownermodel.Owner)) *GetAlbumSharingGridFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ownermodel.Owner))
	})
	return _c
}

func (_c *GetAlbumSharingGridFunc_Execute_Call) Return(_a0 map[catalog.AlbumId][]usermodel.UserId, _a1 error) *GetAlbumSharingGridFunc_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *GetAlbumSharingGridFunc_Execute_Call) RunAndReturn(run func(context.Context, ownermodel.Owner) (map[catalog.AlbumId][]usermodel.UserId, error)) *GetAlbumSharingGridFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewGetAlbumSharingGridFunc creates a new instance of GetAlbumSharingGridFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGetAlbumSharingGridFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *GetAlbumSharingGridFunc {
	mock := &GetAlbumSharingGridFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}