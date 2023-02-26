// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	aclcore "github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
)

// IdentityDetailsStore is an autogenerated mock type for the IdentityDetailsStore type
type IdentityDetailsStore struct {
	mock.Mock
}

// FetchIdentity provides a mock function with given fields: email
func (_m *IdentityDetailsStore) FindIdentity(email string) (*aclcore.Identity, error) {
	ret := _m.Called(email)

	var r0 *aclcore.Identity
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*aclcore.Identity, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) *aclcore.Identity); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*aclcore.Identity)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreIdentity provides a mock function with given fields: identity
func (_m *IdentityDetailsStore) StoreIdentity(identity aclcore.Identity) error {
	ret := _m.Called(identity)

	var r0 error
	if rf, ok := ret.Get(0).(func(aclcore.Identity) error); ok {
		r0 = rf(identity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewIdentityDetailsStore interface {
	mock.TestingT
	Cleanup(func())
}

// NewIdentityDetailsStore creates a new instance of IdentityDetailsStore. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIdentityDetailsStore(t mockConstructorTestingTNewIdentityDetailsStore) *IdentityDetailsStore {
	mock := &IdentityDetailsStore{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
