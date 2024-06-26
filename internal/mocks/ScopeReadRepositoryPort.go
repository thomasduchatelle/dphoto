// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	aclcore "github.com/thomasduchatelle/dphoto/pkg/acl/aclcore"

	context "context"

	mock "github.com/stretchr/testify/mock"

	ownermodel "github.com/thomasduchatelle/dphoto/pkg/ownermodel"

	usermodel "github.com/thomasduchatelle/dphoto/pkg/usermodel"
)

// ScopeReadRepositoryPort is an autogenerated mock type for the ScopeReadRepositoryPort type
type ScopeReadRepositoryPort struct {
	mock.Mock
}

type ScopeReadRepositoryPort_Expecter struct {
	mock *mock.Mock
}

func (_m *ScopeReadRepositoryPort) EXPECT() *ScopeReadRepositoryPort_Expecter {
	return &ScopeReadRepositoryPort_Expecter{mock: &_m.Mock}
}

// ListScopesByOwner provides a mock function with given fields: ctx, owner, types
func (_m *ScopeReadRepositoryPort) ListScopesByOwner(ctx context.Context, owner ownermodel.Owner, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	_va := make([]interface{}, len(types))
	for _i := range types {
		_va[_i] = types[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, owner)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListScopesByOwner")
	}

	var r0 []*aclcore.Scope
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner, ...aclcore.ScopeType) ([]*aclcore.Scope, error)); ok {
		return rf(ctx, owner, types...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ownermodel.Owner, ...aclcore.ScopeType) []*aclcore.Scope); ok {
		r0 = rf(ctx, owner, types...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*aclcore.Scope)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ownermodel.Owner, ...aclcore.ScopeType) error); ok {
		r1 = rf(ctx, owner, types...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ScopeReadRepositoryPort_ListScopesByOwner_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListScopesByOwner'
type ScopeReadRepositoryPort_ListScopesByOwner_Call struct {
	*mock.Call
}

// ListScopesByOwner is a helper method to define mock.On call
//   - ctx context.Context
//   - owner ownermodel.Owner
//   - types ...aclcore.ScopeType
func (_e *ScopeReadRepositoryPort_Expecter) ListScopesByOwner(ctx interface{}, owner interface{}, types ...interface{}) *ScopeReadRepositoryPort_ListScopesByOwner_Call {
	return &ScopeReadRepositoryPort_ListScopesByOwner_Call{Call: _e.mock.On("ListScopesByOwner",
		append([]interface{}{ctx, owner}, types...)...)}
}

func (_c *ScopeReadRepositoryPort_ListScopesByOwner_Call) Run(run func(ctx context.Context, owner ownermodel.Owner, types ...aclcore.ScopeType)) *ScopeReadRepositoryPort_ListScopesByOwner_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]aclcore.ScopeType, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(aclcore.ScopeType)
			}
		}
		run(args[0].(context.Context), args[1].(ownermodel.Owner), variadicArgs...)
	})
	return _c
}

func (_c *ScopeReadRepositoryPort_ListScopesByOwner_Call) Return(_a0 []*aclcore.Scope, _a1 error) *ScopeReadRepositoryPort_ListScopesByOwner_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ScopeReadRepositoryPort_ListScopesByOwner_Call) RunAndReturn(run func(context.Context, ownermodel.Owner, ...aclcore.ScopeType) ([]*aclcore.Scope, error)) *ScopeReadRepositoryPort_ListScopesByOwner_Call {
	_c.Call.Return(run)
	return _c
}

// ListScopesByUser provides a mock function with given fields: ctx, id, types
func (_m *ScopeReadRepositoryPort) ListScopesByUser(ctx context.Context, id usermodel.UserId, types ...aclcore.ScopeType) ([]*aclcore.Scope, error) {
	_va := make([]interface{}, len(types))
	for _i := range types {
		_va[_i] = types[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, id)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ListScopesByUser")
	}

	var r0 []*aclcore.Scope
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, usermodel.UserId, ...aclcore.ScopeType) ([]*aclcore.Scope, error)); ok {
		return rf(ctx, id, types...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, usermodel.UserId, ...aclcore.ScopeType) []*aclcore.Scope); ok {
		r0 = rf(ctx, id, types...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*aclcore.Scope)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, usermodel.UserId, ...aclcore.ScopeType) error); ok {
		r1 = rf(ctx, id, types...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ScopeReadRepositoryPort_ListScopesByUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListScopesByUser'
type ScopeReadRepositoryPort_ListScopesByUser_Call struct {
	*mock.Call
}

// ListScopesByUser is a helper method to define mock.On call
//   - ctx context.Context
//   - id usermodel.UserId
//   - types ...aclcore.ScopeType
func (_e *ScopeReadRepositoryPort_Expecter) ListScopesByUser(ctx interface{}, id interface{}, types ...interface{}) *ScopeReadRepositoryPort_ListScopesByUser_Call {
	return &ScopeReadRepositoryPort_ListScopesByUser_Call{Call: _e.mock.On("ListScopesByUser",
		append([]interface{}{ctx, id}, types...)...)}
}

func (_c *ScopeReadRepositoryPort_ListScopesByUser_Call) Run(run func(ctx context.Context, id usermodel.UserId, types ...aclcore.ScopeType)) *ScopeReadRepositoryPort_ListScopesByUser_Call {
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

func (_c *ScopeReadRepositoryPort_ListScopesByUser_Call) Return(_a0 []*aclcore.Scope, _a1 error) *ScopeReadRepositoryPort_ListScopesByUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ScopeReadRepositoryPort_ListScopesByUser_Call) RunAndReturn(run func(context.Context, usermodel.UserId, ...aclcore.ScopeType) ([]*aclcore.Scope, error)) *ScopeReadRepositoryPort_ListScopesByUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewScopeReadRepositoryPort creates a new instance of ScopeReadRepositoryPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewScopeReadRepositoryPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *ScopeReadRepositoryPort {
	mock := &ScopeReadRepositoryPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
