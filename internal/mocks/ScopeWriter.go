// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	aclcore "github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
)

// ScopeWriter is an autogenerated mock type for the ScopeWriter type
type ScopeWriter struct {
	mock.Mock
}

// DeleteScopes provides a mock function with given fields: id
func (_m *ScopeWriter) DeleteScopes(id ...aclcore.ScopeId) error {
	_va := make([]interface{}, len(id))
	for _i := range id {
		_va[_i] = id[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteScopes")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...aclcore.ScopeId) error); ok {
		r0 = rf(id...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveIfNewScope provides a mock function with given fields: scope
func (_m *ScopeWriter) SaveIfNewScope(scope aclcore.Scope) error {
	ret := _m.Called(scope)

	if len(ret) == 0 {
		panic("no return value specified for SaveIfNewScope")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(aclcore.Scope) error); ok {
		r0 = rf(scope)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewScopeWriter creates a new instance of ScopeWriter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewScopeWriter(t interface {
	mock.TestingT
	Cleanup(func())
}) *ScopeWriter {
	mock := &ScopeWriter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
