// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	aclcore "github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
)

// ScopesReader is an autogenerated mock type for the ScopesReader type
type ScopesReader struct {
	mock.Mock
}

// FindScopesById provides a mock function with given fields: ids
func (_m *ScopesReader) FindScopesById(ids ...aclcore.ScopeId) ([]*aclcore.Scope, error) {
	_va := make([]interface{}, len(ids))
	for _i := range ids {
		_va[_i] = ids[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindScopesById")
	}

	var r0 []*aclcore.Scope
	var r1 error
	if rf, ok := ret.Get(0).(func(...aclcore.ScopeId) ([]*aclcore.Scope, error)); ok {
		return rf(ids...)
	}
	if rf, ok := ret.Get(0).(func(...aclcore.ScopeId) []*aclcore.Scope); ok {
		r0 = rf(ids...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*aclcore.Scope)
		}
	}

	if rf, ok := ret.Get(1).(func(...aclcore.ScopeId) error); ok {
		r1 = rf(ids...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListUserScopes provides a mock function with given fields: email, types
func (_m *ScopesReader) ListUserScopes(email string, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	_va := make([]interface{}, len(types))
	for _i := range types {
		_va[_i] = types[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, email)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListUserScopes")
	}

	var r0 []*aclcore.Scope
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...aclcore.ScopeType) ([]*aclcore.Scope, error)); ok {
		return rf(email, types...)
	}
	if rf, ok := ret.Get(0).(func(string, ...aclcore.ScopeType) []*aclcore.Scope); ok {
		r0 = rf(email, types...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*aclcore.Scope)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...aclcore.ScopeType) error); ok {
		r1 = rf(email, types...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewScopesReader creates a new instance of ScopesReader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewScopesReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *ScopesReader {
	mock := &ScopesReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
