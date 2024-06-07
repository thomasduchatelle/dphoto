// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"
	catalogviews "github.com/thomasduchatelle/dphoto/pkg/catalogviews"

	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ListUserWhoCanAccessAlbumPort is an autogenerated mock type for the ListUserWhoCanAccessAlbumPort type
type ListUserWhoCanAccessAlbumPort struct {
	mock.Mock
}

type ListUserWhoCanAccessAlbumPort_Expecter struct {
	mock *mock.Mock
}

func (_m *ListUserWhoCanAccessAlbumPort) EXPECT() *ListUserWhoCanAccessAlbumPort_Expecter {
	return &ListUserWhoCanAccessAlbumPort_Expecter{mock: &_m.Mock}
}

// ListUsersWhoCanAccessAlbum provides a mock function with given fields: ctx, albumId
func (_m *ListUserWhoCanAccessAlbumPort) ListUsersWhoCanAccessAlbum(ctx context.Context, albumId ...catalog.AlbumId) (map[catalog.AlbumId][]catalogviews.Availability, error) {
	_va := make([]interface{}, len(albumId))
	for _i := range albumId {
		_va[_i] = albumId[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListUsersWhoCanAccessAlbum")
	}

	var r0 map[catalog.AlbumId][]catalogviews.Availability
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...catalog.AlbumId) (map[catalog.AlbumId][]catalogviews.Availability, error)); ok {
		return rf(ctx, albumId...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...catalog.AlbumId) map[catalog.AlbumId][]catalogviews.Availability); ok {
		r0 = rf(ctx, albumId...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[catalog.AlbumId][]catalogviews.Availability)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...catalog.AlbumId) error); ok {
		r1 = rf(ctx, albumId...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUserWhoCanAccessAlbumPort_ListUsersWhoCanAccessAlbum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListUsersWhoCanAccessAlbum'
type ListUserWhoCanAccessAlbumPort_ListUsersWhoCanAccessAlbum_Call struct {
	*mock.Call
}

// ListUsersWhoCanAccessAlbum is a helper method to define mock.On call
//   - ctx context.Context
//   - albumId ...catalog.AlbumId
func (_e *ListUserWhoCanAccessAlbumPort_Expecter) ListUsersWhoCanAccessAlbum(ctx interface{}, albumId ...interface{}) *ListUserWhoCanAccessAlbumPort_ListUsersWhoCanAccessAlbum_Call {
	return &ListUserWhoCanAccessAlbumPort_ListUsersWhoCanAccessAlbum_Call{Call: _e.mock.On("ListUsersWhoCanAccessAlbum",
		append([]interface{}{ctx}, albumId...)...)}
}

func (_c *ListUserWhoCanAccessAlbumPort_ListUsersWhoCanAccessAlbum_Call) Run(run func(ctx context.Context, albumId ...catalog.AlbumId)) *ListUserWhoCanAccessAlbumPort_ListUsersWhoCanAccessAlbum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]catalog.AlbumId, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(catalog.AlbumId)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *ListUserWhoCanAccessAlbumPort_ListUsersWhoCanAccessAlbum_Call) Return(_a0 map[catalog.AlbumId][]catalogviews.Availability, _a1 error) *ListUserWhoCanAccessAlbumPort_ListUsersWhoCanAccessAlbum_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ListUserWhoCanAccessAlbumPort_ListUsersWhoCanAccessAlbum_Call) RunAndReturn(run func(context.Context, ...catalog.AlbumId) (map[catalog.AlbumId][]catalogviews.Availability, error)) *ListUserWhoCanAccessAlbumPort_ListUsersWhoCanAccessAlbum_Call {
	_c.Call.Return(run)
	return _c
}

// NewListUserWhoCanAccessAlbumPort creates a new instance of ListUserWhoCanAccessAlbumPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewListUserWhoCanAccessAlbumPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *ListUserWhoCanAccessAlbumPort {
	mock := &ListUserWhoCanAccessAlbumPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
