// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	aclcore "github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
	catalog "github.com/thomasduchatelle/dphoto/pkg/catalog"

	mock "github.com/stretchr/testify/mock"

	ownermodel "github.com/thomasduchatelle/dphoto/pkg/ownermodel"

	usermodel "github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// CatalogRules is an autogenerated mock type for the CatalogRules type
type CatalogRules struct {
	mock.Mock
}

type CatalogRules_Expecter struct {
	mock *mock.Mock
}

func (_m *CatalogRules) EXPECT() *CatalogRules_Expecter {
	return &CatalogRules_Expecter{mock: &_m.Mock}
}

// CanListMediasFromAlbum provides a mock function with given fields: id
func (_m *CatalogRules) CanListMediasFromAlbum(id catalog.AlbumId) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for CanListMediasFromAlbum")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(catalog.AlbumId) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CatalogRules_CanListMediasFromAlbum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CanListMediasFromAlbum'
type CatalogRules_CanListMediasFromAlbum_Call struct {
	*mock.Call
}

// CanListMediasFromAlbum is a helper method to define mock.On call
//   - id catalog.AlbumId
func (_e *CatalogRules_Expecter) CanListMediasFromAlbum(id interface{}) *CatalogRules_CanListMediasFromAlbum_Call {
	return &CatalogRules_CanListMediasFromAlbum_Call{Call: _e.mock.On("CanListMediasFromAlbum", id)}
}

func (_c *CatalogRules_CanListMediasFromAlbum_Call) Run(run func(id catalog.AlbumId)) *CatalogRules_CanListMediasFromAlbum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(catalog.AlbumId))
	})
	return _c
}

func (_c *CatalogRules_CanListMediasFromAlbum_Call) Return(_a0 error) *CatalogRules_CanListMediasFromAlbum_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CatalogRules_CanListMediasFromAlbum_Call) RunAndReturn(run func(catalog.AlbumId) error) *CatalogRules_CanListMediasFromAlbum_Call {
	_c.Call.Return(run)
	return _c
}

// CanManageAlbum provides a mock function with given fields: id
func (_m *CatalogRules) CanManageAlbum(id catalog.AlbumId) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for CanShareAlbum")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(catalog.AlbumId) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CatalogRules_CanManageAlbum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CanShareAlbum'
type CatalogRules_CanManageAlbum_Call struct {
	*mock.Call
}

// CanManageAlbum is a helper method to define mock.On call
//   - id catalog.AlbumId
func (_e *CatalogRules_Expecter) CanManageAlbum(id interface{}) *CatalogRules_CanManageAlbum_Call {
	return &CatalogRules_CanManageAlbum_Call{Call: _e.mock.On("CanShareAlbum", id)}
}

func (_c *CatalogRules_CanManageAlbum_Call) Run(run func(id catalog.AlbumId)) *CatalogRules_CanManageAlbum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(catalog.AlbumId))
	})
	return _c
}

func (_c *CatalogRules_CanManageAlbum_Call) Return(_a0 error) *CatalogRules_CanManageAlbum_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CatalogRules_CanManageAlbum_Call) RunAndReturn(run func(catalog.AlbumId) error) *CatalogRules_CanManageAlbum_Call {
	_c.Call.Return(run)
	return _c
}

// CanReadMedia provides a mock function with given fields: owner, id
func (_m *CatalogRules) CanReadMedia(owner ownermodel.Owner, id catalog.MediaId) error {
	ret := _m.Called(owner, id)

	if len(ret) == 0 {
		panic("no return value specified for CanReadMedia")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(ownermodel.Owner, catalog.MediaId) error); ok {
		r0 = rf(owner, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CatalogRules_CanReadMedia_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CanReadMedia'
type CatalogRules_CanReadMedia_Call struct {
	*mock.Call
}

// CanReadMedia is a helper method to define mock.On call
//   - owner ownermodel.Owner
//   - id catalog.MediaId
func (_e *CatalogRules_Expecter) CanReadMedia(owner interface{}, id interface{}) *CatalogRules_CanReadMedia_Call {
	return &CatalogRules_CanReadMedia_Call{Call: _e.mock.On("CanReadMedia", owner, id)}
}

func (_c *CatalogRules_CanReadMedia_Call) Run(run func(owner ownermodel.Owner, id catalog.MediaId)) *CatalogRules_CanReadMedia_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(ownermodel.Owner), args[1].(catalog.MediaId))
	})
	return _c
}

func (_c *CatalogRules_CanReadMedia_Call) Return(_a0 error) *CatalogRules_CanReadMedia_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CatalogRules_CanReadMedia_Call) RunAndReturn(run func(ownermodel.Owner, catalog.MediaId) error) *CatalogRules_CanReadMedia_Call {
	_c.Call.Return(run)
	return _c
}

// Owner provides a mock function with given fields:
func (_m *CatalogRules) Owner() (*ownermodel.Owner, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Owner")
	}

	var r0 *ownermodel.Owner
	var r1 error
	if rf, ok := ret.Get(0).(func() (*ownermodel.Owner, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *ownermodel.Owner); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ownermodel.Owner)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CatalogRules_Owner_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Owner'
type CatalogRules_Owner_Call struct {
	*mock.Call
}

// Owner is a helper method to define mock.On call
func (_e *CatalogRules_Expecter) Owner() *CatalogRules_Owner_Call {
	return &CatalogRules_Owner_Call{Call: _e.mock.On("Owner")}
}

func (_c *CatalogRules_Owner_Call) Run(run func()) *CatalogRules_Owner_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CatalogRules_Owner_Call) Return(_a0 *ownermodel.Owner, _a1 error) *CatalogRules_Owner_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CatalogRules_Owner_Call) RunAndReturn(run func() (*ownermodel.Owner, error)) *CatalogRules_Owner_Call {
	_c.Call.Return(run)
	return _c
}

// SharedByUserGrid provides a mock function with given fields: owner
func (_m *CatalogRules) SharedByUserGrid(owner ownermodel.Owner) (map[string]map[usermodel.UserId]aclcore.ScopeType, error) {
	ret := _m.Called(owner)

	if len(ret) == 0 {
		panic("no return value specified for SharedByUserGrid")
	}

	var r0 map[string]map[usermodel.UserId]aclcore.ScopeType
	var r1 error
	if rf, ok := ret.Get(0).(func(ownermodel.Owner) (map[string]map[usermodel.UserId]aclcore.ScopeType, error)); ok {
		return rf(owner)
	}
	if rf, ok := ret.Get(0).(func(ownermodel.Owner) map[string]map[usermodel.UserId]aclcore.ScopeType); ok {
		r0 = rf(owner)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]map[usermodel.UserId]aclcore.ScopeType)
		}
	}

	if rf, ok := ret.Get(1).(func(ownermodel.Owner) error); ok {
		r1 = rf(owner)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CatalogRules_SharedByUserGrid_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SharedByUserGrid'
type CatalogRules_SharedByUserGrid_Call struct {
	*mock.Call
}

// SharedByUserGrid is a helper method to define mock.On call
//   - owner ownermodel.Owner
func (_e *CatalogRules_Expecter) SharedByUserGrid(owner interface{}) *CatalogRules_SharedByUserGrid_Call {
	return &CatalogRules_SharedByUserGrid_Call{Call: _e.mock.On("SharedByUserGrid", owner)}
}

func (_c *CatalogRules_SharedByUserGrid_Call) Run(run func(owner ownermodel.Owner)) *CatalogRules_SharedByUserGrid_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(ownermodel.Owner))
	})
	return _c
}

func (_c *CatalogRules_SharedByUserGrid_Call) Return(_a0 map[string]map[usermodel.UserId]aclcore.ScopeType, _a1 error) *CatalogRules_SharedByUserGrid_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CatalogRules_SharedByUserGrid_Call) RunAndReturn(run func(ownermodel.Owner) (map[string]map[usermodel.UserId]aclcore.ScopeType, error)) *CatalogRules_SharedByUserGrid_Call {
	_c.Call.Return(run)
	return _c
}

// SharedWithUserAlbum provides a mock function with given fields:
func (_m *CatalogRules) SharedWithUserAlbum() ([]catalog.AlbumId, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SharedWithUserAlbum")
	}

	var r0 []catalog.AlbumId
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]catalog.AlbumId, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []catalog.AlbumId); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]catalog.AlbumId)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CatalogRules_SharedWithUserAlbum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SharedWithUserAlbum'
type CatalogRules_SharedWithUserAlbum_Call struct {
	*mock.Call
}

// SharedWithUserAlbum is a helper method to define mock.On call
func (_e *CatalogRules_Expecter) SharedWithUserAlbum() *CatalogRules_SharedWithUserAlbum_Call {
	return &CatalogRules_SharedWithUserAlbum_Call{Call: _e.mock.On("SharedWithUserAlbum")}
}

func (_c *CatalogRules_SharedWithUserAlbum_Call) Run(run func()) *CatalogRules_SharedWithUserAlbum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *CatalogRules_SharedWithUserAlbum_Call) Return(_a0 []catalog.AlbumId, _a1 error) *CatalogRules_SharedWithUserAlbum_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CatalogRules_SharedWithUserAlbum_Call) RunAndReturn(run func() ([]catalog.AlbumId, error)) *CatalogRules_SharedWithUserAlbum_Call {
	_c.Call.Return(run)
	return _c
}

// NewCatalogRules creates a new instance of CatalogRules. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCatalogRules(t interface {
	mock.TestingT
	Cleanup(func())
}) *CatalogRules {
	mock := &CatalogRules{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
