// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	catalogviews "github.com/thomasduchatelle/dphoto/pkg/catalogviews"

	mock "github.com/stretchr/testify/mock"

	usermodel "github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// ListAlbumsProvider is an autogenerated mock type for the ListAlbumsProvider type
type ListAlbumsProvider struct {
	mock.Mock
}

type ListAlbumsProvider_Expecter struct {
	mock *mock.Mock
}

func (_m *ListAlbumsProvider) EXPECT() *ListAlbumsProvider_Expecter {
	return &ListAlbumsProvider_Expecter{mock: &_m.Mock}
}

// ListAlbums provides a mock function with given fields: ctx, user, filter
func (_m *ListAlbumsProvider) ListAlbums(ctx context.Context, user usermodel.CurrentUser, filter catalogviews.ListAlbumsFilter) ([]*catalogviews.VisibleAlbum, error) {
	ret := _m.Called(ctx, user, filter)

	if len(ret) == 0 {
		panic("no return value specified for ListAlbums")
	}

	var r0 []*catalogviews.VisibleAlbum
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, usermodel.CurrentUser, catalogviews.ListAlbumsFilter) ([]*catalogviews.VisibleAlbum, error)); ok {
		return rf(ctx, user, filter)
	}
	if rf, ok := ret.Get(0).(func(context.Context, usermodel.CurrentUser, catalogviews.ListAlbumsFilter) []*catalogviews.VisibleAlbum); ok {
		r0 = rf(ctx, user, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*catalogviews.VisibleAlbum)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, usermodel.CurrentUser, catalogviews.ListAlbumsFilter) error); ok {
		r1 = rf(ctx, user, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAlbumsProvider_ListAlbums_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListAlbums'
type ListAlbumsProvider_ListAlbums_Call struct {
	*mock.Call
}

// ListAlbums is a helper method to define mock.On call
//   - ctx context.Context
//   - user usermodel.CurrentUser
//   - filter catalogviews.ListAlbumsFilter
func (_e *ListAlbumsProvider_Expecter) ListAlbums(ctx interface{}, user interface{}, filter interface{}) *ListAlbumsProvider_ListAlbums_Call {
	return &ListAlbumsProvider_ListAlbums_Call{Call: _e.mock.On("ListAlbums", ctx, user, filter)}
}

func (_c *ListAlbumsProvider_ListAlbums_Call) Run(run func(ctx context.Context, user usermodel.CurrentUser, filter catalogviews.ListAlbumsFilter)) *ListAlbumsProvider_ListAlbums_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(usermodel.CurrentUser), args[2].(catalogviews.ListAlbumsFilter))
	})
	return _c
}

func (_c *ListAlbumsProvider_ListAlbums_Call) Return(_a0 []*catalogviews.VisibleAlbum, _a1 error) *ListAlbumsProvider_ListAlbums_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ListAlbumsProvider_ListAlbums_Call) RunAndReturn(run func(context.Context, usermodel.CurrentUser, catalogviews.ListAlbumsFilter) ([]*catalogviews.VisibleAlbum, error)) *ListAlbumsProvider_ListAlbums_Call {
	_c.Call.Return(run)
	return _c
}

// NewListAlbumsProvider creates a new instance of ListAlbumsProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewListAlbumsProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *ListAlbumsProvider {
	mock := &ListAlbumsProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
