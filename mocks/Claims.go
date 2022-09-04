// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Claims is an autogenerated mock type for the Claims type
type Claims struct {
	mock.Mock
}

// HasApiAccess provides a mock function with given fields: api
func (_m *Claims) HasApiAccess(api string) error {
	ret := _m.Called(api)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(api)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsOwnerOf provides a mock function with given fields: owner
func (_m *Claims) IsOwnerOf(owner string) error {
	ret := _m.Called(owner)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(owner)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewClaims interface {
	mock.TestingT
	Cleanup(func())
}

// NewClaims creates a new instance of Claims. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClaims(t mockConstructorTestingTNewClaims) *Claims {
	mock := &Claims{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
