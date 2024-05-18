// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	aclcore "github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"

	mock "github.com/stretchr/testify/mock"

	usermodel "github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// ScopesReader is an autogenerated mock type for the ScopesReader type
type ScopesReader struct {
	mock.Mock
}

type ScopesReader_Expecter struct {
	mock *mock.Mock
}

func (_m *ScopesReader) EXPECT() *ScopesReader_Expecter {
	return &ScopesReader_Expecter{mock: &_m.Mock}
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

// ScopesReader_FindScopesById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindScopesById'
type ScopesReader_FindScopesById_Call struct {
	*mock.Call
}

// FindScopesById is a helper method to define mock.On call
//   - ids ...aclcore.ScopeId
func (_e *ScopesReader_Expecter) FindScopesById(ids ...interface{}) *ScopesReader_FindScopesById_Call {
	return &ScopesReader_FindScopesById_Call{Call: _e.mock.On("FindScopesById",
		append([]interface{}{}, ids...)...)}
}

func (_c *ScopesReader_FindScopesById_Call) Run(run func(ids ...aclcore.ScopeId)) *ScopesReader_FindScopesById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]aclcore.ScopeId, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(aclcore.ScopeId)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *ScopesReader_FindScopesById_Call) Return(_a0 []*aclcore.Scope, _a1 error) *ScopesReader_FindScopesById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ScopesReader_FindScopesById_Call) RunAndReturn(run func(...aclcore.ScopeId) ([]*aclcore.Scope, error)) *ScopesReader_FindScopesById_Call {
	_c.Call.Return(run)
	return _c
}

// ListScopesByUser provides a mock function with given fields: ctx, email, types
func (_m *ScopesReader) ListScopesByUser(ctx context.Context, email usermodel.UserId, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	_va := make([]interface{}, len(types))
	for _i := range types {
		_va[_i] = types[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, email)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListScopesByUser")
	}

	var r0 []*aclcore.Scope
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, usermodel.UserId, ...aclcore.ScopeType) ([]*aclcore.Scope, error)); ok {
		return rf(ctx, email, types...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, usermodel.UserId, ...aclcore.ScopeType) []*aclcore.Scope); ok {
		r0 = rf(ctx, email, types...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*aclcore.Scope)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, usermodel.UserId, ...aclcore.ScopeType) error); ok {
		r1 = rf(ctx, email, types...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ScopesReader_ListScopesByUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListScopesByUser'
type ScopesReader_ListScopesByUser_Call struct {
	*mock.Call
}

// ListScopesByUser is a helper method to define mock.On call
//   - ctx context.Context
//   - email usermodel.UserId
//   - types ...aclcore.ScopeType
func (_e *ScopesReader_Expecter) ListScopesByUser(ctx interface{}, email interface{}, types ...interface{}) *ScopesReader_ListScopesByUser_Call {
	return &ScopesReader_ListScopesByUser_Call{Call: _e.mock.On("ListScopesByUser",
		append([]interface{}{ctx, email}, types...)...)}
}

func (_c *ScopesReader_ListScopesByUser_Call) Run(run func(ctx context.Context, email usermodel.UserId, types ...aclcore.ScopeType)) *ScopesReader_ListScopesByUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]aclcore.ScopeType, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(aclcore.ScopeType)
			}
		}
		run(args[0].(context.Context), args[1].(usermodel.UserId), variadicArgs...)
	})
	return _c
}

func (_c *ScopesReader_ListScopesByUser_Call) Return(_a0 []*aclcore.Scope, _a1 error) *ScopesReader_ListScopesByUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ScopesReader_ListScopesByUser_Call) RunAndReturn(run func(context.Context, usermodel.UserId, ...aclcore.ScopeType) ([]*aclcore.Scope, error)) *ScopesReader_ListScopesByUser_Call {
	_c.Call.Return(run)
	return _c
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
