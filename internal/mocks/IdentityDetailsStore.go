// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	aclcore "github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"

	usermodel "github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// IdentityDetailsStore is an autogenerated mock type for the IdentityDetailsStore type
type IdentityDetailsStore struct {
	mock.Mock
}

type IdentityDetailsStore_Expecter struct {
	mock *mock.Mock
}

func (_m *IdentityDetailsStore) EXPECT() *IdentityDetailsStore_Expecter {
	return &IdentityDetailsStore_Expecter{mock: &_m.Mock}
}

// FindIdentity provides a mock function with given fields: email
func (_m *IdentityDetailsStore) FindIdentity(email usermodel.UserId) (*aclcore.Identity, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for FindIdentity")
	}

	var r0 *aclcore.Identity
	var r1 error
	if rf, ok := ret.Get(0).(func(usermodel.UserId) (*aclcore.Identity, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(usermodel.UserId) *aclcore.Identity); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*aclcore.Identity)
		}
	}

	if rf, ok := ret.Get(1).(func(usermodel.UserId) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IdentityDetailsStore_FindIdentity_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindIdentity'
type IdentityDetailsStore_FindIdentity_Call struct {
	*mock.Call
}

// FindIdentity is a helper method to define mock.On call
//   - email usermodel.UserId
func (_e *IdentityDetailsStore_Expecter) FindIdentity(email interface{}) *IdentityDetailsStore_FindIdentity_Call {
	return &IdentityDetailsStore_FindIdentity_Call{Call: _e.mock.On("FindIdentity", email)}
}

func (_c *IdentityDetailsStore_FindIdentity_Call) Run(run func(email usermodel.UserId)) *IdentityDetailsStore_FindIdentity_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(usermodel.UserId))
	})
	return _c
}

func (_c *IdentityDetailsStore_FindIdentity_Call) Return(_a0 *aclcore.Identity, _a1 error) *IdentityDetailsStore_FindIdentity_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *IdentityDetailsStore_FindIdentity_Call) RunAndReturn(run func(usermodel.UserId) (*aclcore.Identity, error)) *IdentityDetailsStore_FindIdentity_Call {
	_c.Call.Return(run)
	return _c
}

// StoreIdentity provides a mock function with given fields: identity
func (_m *IdentityDetailsStore) StoreIdentity(identity aclcore.Identity) error {
	ret := _m.Called(identity)

	if len(ret) == 0 {
		panic("no return value specified for StoreIdentity")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(aclcore.Identity) error); ok {
		r0 = rf(identity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IdentityDetailsStore_StoreIdentity_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StoreIdentity'
type IdentityDetailsStore_StoreIdentity_Call struct {
	*mock.Call
}

// StoreIdentity is a helper method to define mock.On call
//   - identity aclcore.Identity
func (_e *IdentityDetailsStore_Expecter) StoreIdentity(identity interface{}) *IdentityDetailsStore_StoreIdentity_Call {
	return &IdentityDetailsStore_StoreIdentity_Call{Call: _e.mock.On("StoreIdentity", identity)}
}

func (_c *IdentityDetailsStore_StoreIdentity_Call) Run(run func(identity aclcore.Identity)) *IdentityDetailsStore_StoreIdentity_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(aclcore.Identity))
	})
	return _c
}

func (_c *IdentityDetailsStore_StoreIdentity_Call) Return(_a0 error) *IdentityDetailsStore_StoreIdentity_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *IdentityDetailsStore_StoreIdentity_Call) RunAndReturn(run func(aclcore.Identity) error) *IdentityDetailsStore_StoreIdentity_Call {
	_c.Call.Return(run)
	return _c
}

// NewIdentityDetailsStore creates a new instance of IdentityDetailsStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIdentityDetailsStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *IdentityDetailsStore {
	mock := &IdentityDetailsStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
