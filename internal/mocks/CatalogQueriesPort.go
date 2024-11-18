// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	context "context"

	mock "github.com/stretchr/testify/mock"

	ownermodel "github.com/thomasduchatelle/dphoto/pkg/ownermodel"
)

// CatalogQueriesPort is an autogenerated mock type for the CatalogQueriesPort type
type CatalogQueriesPort struct {
	mock.Mock
}

type CatalogQueriesPort_Expecter struct {
	mock *mock.Mock
}

func (_m *CatalogQueriesPort) EXPECT() *CatalogQueriesPort_Expecter {
	return &CatalogQueriesPort_Expecter{mock: &_m.Mock}
}

// FindMediaOwnership provides a mock function with given fields: ctx, owner, mediaId
func (_m *CatalogQueriesPort) FindMediaOwnership(ctx context.Context, owner ownermodel.Owner, mediaId catalog.MediaId) (*catalog.AlbumId, error) {
	ret := _m.Called(ctx, owner, mediaId)

	if len(ret) == 0 {
		panic("no return value specified for FindMediaOwnership")
	}

	var r0 *catalog.AlbumId
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner, catalog.MediaId) (*catalog.AlbumId, error)); ok {
		return rf(ctx, owner, mediaId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner, catalog.MediaId) *catalog.AlbumId); ok {
		r0 = rf(ctx, owner, mediaId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*catalog.AlbumId)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ownermodel.Owner, catalog.MediaId) error); ok {
		r1 = rf(ctx, owner, mediaId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CatalogQueriesPort_FindMediaOwnership_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindMediaOwnership'
type CatalogQueriesPort_FindMediaOwnership_Call struct {
	*mock.Call
}

// FindMediaOwnership is a helper method to define mock.On call
//   - ctx context.Context
//   - owner ownermodel.Owner
//   - mediaId catalog.MediaId
func (_e *CatalogQueriesPort_Expecter) FindMediaOwnership(ctx interface{}, owner interface{}, mediaId interface{}) *CatalogQueriesPort_FindMediaOwnership_Call {
	return &CatalogQueriesPort_FindMediaOwnership_Call{Call: _e.mock.On("FindMediaOwnership", ctx, owner, mediaId)}
}

func (_c *CatalogQueriesPort_FindMediaOwnership_Call) Run(run func(ctx context.Context, owner ownermodel.Owner, mediaId catalog.MediaId)) *CatalogQueriesPort_FindMediaOwnership_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ownermodel.Owner), args[2].(catalog.MediaId))
	})
	return _c
}

func (_c *CatalogQueriesPort_FindMediaOwnership_Call) Return(_a0 *catalog.AlbumId, _a1 error) *CatalogQueriesPort_FindMediaOwnership_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CatalogQueriesPort_FindMediaOwnership_Call) RunAndReturn(run func(context.Context, ownermodel.Owner, catalog.MediaId) (*catalog.AlbumId, error)) *CatalogQueriesPort_FindMediaOwnership_Call {
	_c.Call.Return(run)
	return _c
}

// NewCatalogQueriesPort creates a new instance of CatalogQueriesPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCatalogQueriesPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *CatalogQueriesPort {
	mock := &CatalogQueriesPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}