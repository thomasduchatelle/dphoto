// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	aclcore "github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"
)

// ReverseScopesReader is an autogenerated mock type for the ReverseScopesReader type
type ReverseScopesReader struct {
	mock.Mock
}

type ReverseScopesReader_Expecter struct {
	mock *mock.Mock
}

func (_m *ReverseScopesReader) EXPECT() *ReverseScopesReader_Expecter {
	return &ReverseScopesReader_Expecter{mock: &_m.Mock}
}

// ListOwnerScopes provides a mock function with given fields: owner, types
func (_m *ReverseScopesReader) ListOwnerScopes(owner string, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	_va := make([]interface{}, len(types))
	for _i := range types {
		_va[_i] = types[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, owner)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListOwnerScopes")
	}

	var r0 []*aclcore.Scope
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...aclcore.ScopeType) ([]*aclcore.Scope, error)); ok {
		return rf(owner, types...)
	}
	if rf, ok := ret.Get(0).(func(string, ...aclcore.ScopeType) []*aclcore.Scope); ok {
		r0 = rf(owner, types...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*aclcore.Scope)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...aclcore.ScopeType) error); ok {
		r1 = rf(owner, types...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReverseScopesReader_ListOwnerScopes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListOwnerScopes'
type ReverseScopesReader_ListOwnerScopes_Call struct {
	*mock.Call
}

// ListOwnerScopes is a helper method to define mock.On call
//   - owner string
//   - types ...aclcore.ScopeType
func (_e *ReverseScopesReader_Expecter) ListOwnerScopes(owner interface{}, types ...interface{}) *ReverseScopesReader_ListOwnerScopes_Call {
	return &ReverseScopesReader_ListOwnerScopes_Call{Call: _e.mock.On("ListOwnerScopes",
		append([]interface{}{owner}, types...)...)}
}

func (_c *ReverseScopesReader_ListOwnerScopes_Call) Run(run func(owner string, types ...aclcore.ScopeType)) *ReverseScopesReader_ListOwnerScopes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]aclcore.ScopeType, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(aclcore.ScopeType)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *ReverseScopesReader_ListOwnerScopes_Call) Return(_a0 []*aclcore.Scope, _a1 error) *ReverseScopesReader_ListOwnerScopes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ReverseScopesReader_ListOwnerScopes_Call) RunAndReturn(run func(string, ...aclcore.ScopeType) ([]*aclcore.Scope, error)) *ReverseScopesReader_ListOwnerScopes_Call {
	_c.Call.Return(run)
	return _c
}

// NewReverseScopesReader creates a new instance of ReverseScopesReader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReverseScopesReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *ReverseScopesReader {
	mock := &ReverseScopesReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
