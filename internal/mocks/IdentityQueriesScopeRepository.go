// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	aclcore "github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
)

// IdentityQueriesScopeRepository is an autogenerated mock type for the IdentityQueriesScopeRepository type
type IdentityQueriesScopeRepository struct {
	mock.Mock
}

// ListScopesByOwners provides a mock function with given fields: owners, types
func (_m *IdentityQueriesScopeRepository) ListScopesByOwners(owners []string, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	_va := make([]interface{}, len(types))
	for _i := range types {
		_va[_i] = types[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, owners)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListScopesByOwners")
	}

	var r0 []*aclcore.Scope
	var r1 error
	if rf, ok := ret.Get(0).(func([]string, ...aclcore.ScopeType) ([]*aclcore.Scope, error)); ok {
		return rf(owners, types...)
	}
	if rf, ok := ret.Get(0).(func([]string, ...aclcore.ScopeType) []*aclcore.Scope); ok {
		r0 = rf(owners, types...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*aclcore.Scope)
		}
	}

	if rf, ok := ret.Get(1).(func([]string, ...aclcore.ScopeType) error); ok {
		r1 = rf(owners, types...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIdentityQueriesScopeRepository creates a new instance of IdentityQueriesScopeRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIdentityQueriesScopeRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IdentityQueriesScopeRepository {
	mock := &IdentityQueriesScopeRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
